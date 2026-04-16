package parsing

import (
	"errors"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

// ErrInvalidLineRange is returned when LineRange bounds are invalid.
var ErrInvalidLineRange = errors.New("line_range must be 1-based inclusive with end >= start")

const ellipsis = "..."

// ExtractOptions configures skeleton extraction.
type ExtractOptions struct {
	IncludeImports bool
	Symbol         string
	LineRange      *[2]int // 1-based inclusive lo, hi; nil means full file
}

// ExtractSkeleton returns signatures and docstrings without bodies (coderay parity).
func ExtractSkeleton(path, content string, opt ExtractOptions) (string, error) {
	if opt.LineRange != nil {
		lo, hi := opt.LineRange[0], opt.LineRange[1]
		if lo < 1 || hi < 1 || hi < lo {
			return "", ErrInvalidLineRange
		}
	}

	ctx := NewParseContext(path, content)
	if ctx == nil {
		return content, nil
	}

	absDisplay, err := filepath.Abs(path)
	if err != nil {
		absDisplay = path
	}

	p := newSkeletonParser(ctx, absDisplay, opt)
	var lines []string
	var panicked bool
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
				log.Printf("skeleton: extraction failed: %v", r)
			}
		}()
		lines = p.collectLines()
	}()
	if panicked {
		return content, nil
	}
	if len(lines) == 0 {
		if opt.LineRange != nil {
			a, b := opt.LineRange[0], opt.LineRange[1]
			return "# No skeleton entries in file line range " + strconv.Itoa(a) + "-" + strconv.Itoa(b) + ".", nil
		}
		if opt.Symbol != "" {
			return symbolNotFoundHint(ctx, opt.Symbol), nil
		}
	}
	return strings.Join(lines, "\n"), nil
}

func symbolNotFoundHint(ctx *ParseContext, symbol string) string {
	helper := NewBaseTreeSitterParser(ctx)
	tree := helper.GetTree()
	if tree == nil {
		return "# Symbol '" + symbol + "' not found. Available symbols: (none)"
	}
	root := tree.RootNode()
	if root == nil {
		return "# Symbol '" + symbol + "' not found. Available symbols: (none)"
	}
	langCfg := ctx.LangCfg
	names := map[string]struct{}{}
	for i := 0; i < int(root.ChildCount()); i++ {
		node := root.Child(i)
		if node == nil || node.IsNull() {
			continue
		}
		n := unwrapDecorated(node, langCfg)
		nameNode := n.ChildByFieldName("name")
		if nameNode != nil && !nameNode.IsNull() {
			names[strings.TrimSpace(helper.NodeText(nameNode))] = struct{}{}
		}
	}
	if len(names) == 0 {
		return "# Symbol '" + symbol + "' not found. Available symbols: (none)"
	}
	list := make([]string, 0, len(names))
	for s := range names {
		list = append(list, s)
	}
	sortStrings(list)
	return "# Symbol '" + symbol + "' not found. Available symbols: " + strings.Join(list, ", ")
}

func unwrapDecorated(node *sitter.Node, langCfg *LanguageConfig) *sitter.Node {
	n := node
	if inList(langCfg.Cst.DecoratorScopeTypes, n.Type()) {
		for j := 0; j < int(n.NamedChildCount()); j++ {
			inner := n.NamedChild(j)
			if inner == nil {
				continue
			}
			if inList(langCfg.Cst.FunctionScopeTypes, inner.Type()) ||
				inList(langCfg.Cst.ClassScopeTypes, inner.Type()) {
				return inner
			}
		}
	}
	return n
}

func sortStrings(a []string) {
	// tiny sort without importing sort for small sets
	for i := 0; i < len(a); i++ {
		for j := i + 1; j < len(a); j++ {
			if a[j] < a[i] {
				a[i], a[j] = a[j], a[i]
			}
		}
	}
}

type skeletonParser struct {
	*BaseTreeSitterParser
	absDisplay                 string
	includeImports             bool
	symbol                     string
	symbolParts                []string
	fileLineRange0             *[2]int // 0-based inclusive
	rootID                     uintptr
	rootNode                   *sitter.Node
	omitNextSkeletonPathLine   bool
	seen                       map[uintptr]struct{}
}

func newSkeletonParser(ctx *ParseContext, absDisplay string, opt ExtractOptions) *skeletonParser {
	p := &skeletonParser{
		BaseTreeSitterParser: NewBaseTreeSitterParser(ctx),
		absDisplay:           absDisplay,
		includeImports:       opt.IncludeImports,
		symbol:               opt.Symbol,
		seen:                 map[uintptr]struct{}{},
	}
	if opt.Symbol != "" {
		p.symbolParts = strings.Split(opt.Symbol, ".")
	}
	if opt.LineRange != nil {
		a, b := opt.LineRange[0], opt.LineRange[1]
		p.fileLineRange0 = &[2]int{a - 1, b - 1}
	}
	return p
}

func inclusiveRows1Based(n *sitter.Node) (int, int) {
	sp := n.StartPoint()
	ep := n.EndPoint()
	return int(sp.Row) + 1, int(ep.Row) + 1
}

func (p *skeletonParser) collectLines() []string {
	tree := p.GetTree()
	if tree == nil {
		return nil
	}
	root := tree.RootNode()
	if root == nil {
		return nil
	}
	p.rootNode = root
	p.rootID = root.ID()
	lines := []string{}
	p.dfs(root, &lines, 0)
	return lines
}

func (p *skeletonParser) emitPathHeader(lines *[]string, decl *sitter.Node) {
	lo, hi := inclusiveRows1Based(decl)
	*lines = append(*lines, p.absDisplay+":"+strconv.Itoa(lo)+"-"+strconv.Itoa(hi))
}

func (p *skeletonParser) emitFirst(lines *[]string, decl *sitter.Node, indent, text string) {
	if p.omitNextSkeletonPathLine {
		p.omitNextSkeletonPathLine = false
	} else {
		p.emitPathHeader(lines, decl)
	}
	*lines = append(*lines, indent+text)
}

func (p *skeletonParser) emitCont(lines *[]string, indent, text string) {
	*lines = append(*lines, indent+text)
}

func (p *skeletonParser) extractText(node *sitter.Node) string {
	if node == nil || node.IsNull() {
		return ""
	}
	exprType := p.ctx.LangCfg.Skeleton.DocstringExprType
	if node.Type() != exprType {
		return ""
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		sub := node.Child(i)
		if sub == nil {
			continue
		}
		if sub.Type() == "string" {
			return strings.TrimSpace(p.NodeText(sub))
		}
	}
	return ""
}

func (p *skeletonParser) getDocstring(node *sitter.Node) string {
	if node == nil || node.IsNull() {
		return ""
	}
	langCfg := p.ctx.LangCfg
	scopeTypes := append(append([]string{}, langCfg.Cst.FunctionScopeTypes...), langCfg.Cst.ClassScopeTypes...)
	if inList(scopeTypes, node.Type()) {
		var body *sitter.Node
		for i := 0; i < int(node.ChildCount()); i++ {
			ch := node.Child(i)
			if ch == nil {
				continue
			}
			if inList(langCfg.Skeleton.BodyBlockTypes, ch.Type()) {
				body = ch
				break
			}
		}
		if body != nil {
			for i := 0; i < int(body.ChildCount()); i++ {
				ch := body.Child(i)
				if t := p.extractText(ch); t != "" {
					return t
				}
			}
		}
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		ch := node.Child(i)
		if t := p.extractText(ch); t != "" {
			return t
		}
	}
	return ""
}

func (p *skeletonParser) getSignatureLine(node *sitter.Node) string {
	text := p.NodeText(node)
	for _, delim := range []string{":\n", "{\n", ":\r\n", "{\r\n"} {
		if idx := strings.Index(text, delim); idx >= 0 {
			return text[:idx+1]
		}
	}
	if idx := strings.IndexByte(text, '\n'); idx >= 0 {
		return text[:idx]
	}
	return text
}

func (p *skeletonParser) emitSignatureBlock(lines *[]string, node *sitter.Node, indent string, withEllipsis bool) {
	p.emitFirst(lines, node, indent, p.getSignatureLine(node))
	if doc := p.getDocstring(node); doc != "" {
		p.emitCont(lines, indent, "    "+doc)
	}
	if withEllipsis {
		p.emitCont(lines, indent, "    "+ellipsis)
	}
}

func (p *skeletonParser) nodeName(node *sitter.Node) string {
	if node == nil {
		return ""
	}
	nameNode := node.ChildByFieldName("name")
	if nameNode != nil && !nameNode.IsNull() {
		return strings.TrimSpace(p.NodeText(nameNode))
	}
	return ""
}

func (p *skeletonParser) matchesSymbol(node *sitter.Node, depth int) bool {
	if p.symbol == "" {
		return true
	}
	if depth >= len(p.symbolParts) {
		return true
	}
	return p.nodeName(node) == p.symbolParts[depth]
}

func (p *skeletonParser) isSymbolTarget(node *sitter.Node, depth int) bool {
	if p.symbol == "" {
		return true
	}
	if depth >= len(p.symbolParts)-1 {
		return true
	}
	return inList(p.ctx.LangCfg.Cst.ClassScopeTypes, node.Type())
}

func (p *skeletonParser) rowsVsLineFilter(ns, ne int, contained bool) bool {
	if p.fileLineRange0 == nil {
		return true
	}
	lo, hi := p.fileLineRange0[0], p.fileLineRange0[1]
	if contained {
		return ns >= lo && ne <= hi
	}
	return !(ne < lo || ns > hi)
}

func (p *skeletonParser) nodeOverlapsLineFilter(n *sitter.Node) bool {
	ns := int(n.StartPoint().Row)
	ne := int(n.EndPoint().Row)
	return p.rowsVsLineFilter(ns, ne, false)
}

func (p *skeletonParser) nodeFullyInLineFilter(n *sitter.Node) bool {
	ns := int(n.StartPoint().Row)
	ne := int(n.EndPoint().Row)
	return p.rowsVsLineFilter(ns, ne, true)
}

func (p *skeletonParser) wouldEmitDecl(node *sitter.Node, depth int) bool {
	return p.isSymbolTarget(node, depth) && p.nodeFullyInLineFilter(node)
}

func (p *skeletonParser) isCallArgument(node *sitter.Node) bool {
	parent := node.Parent()
	if parent == nil || parent.IsNull() {
		return false
	}
	return parent.Type() == "arguments"
}

func (p *skeletonParser) decoratedInner(node *sitter.Node) *sitter.Node {
	langCfg := p.ctx.LangCfg
	for i := 0; i < int(node.NamedChildCount()); i++ {
		ch := node.NamedChild(i)
		if ch == nil {
			continue
		}
		if inList(langCfg.Cst.FunctionScopeTypes, ch.Type()) ||
			inList(langCfg.Cst.ClassScopeTypes, ch.Type()) {
			return ch
		}
	}
	return nil
}

func (p *skeletonParser) dfs(node *sitter.Node, lines *[]string, depth int) {
	if node == nil || node.IsNull() {
		return
	}
	if _, ok := p.seen[node.ID()]; ok {
		return
	}

	if p.fileLineRange0 != nil && node.ID() != p.rootID && !p.nodeOverlapsLineFilter(node) {
		return
	}

	indent := strings.Repeat("    ", depth)
	ntype := node.Type()
	langCfg := p.ctx.LangCfg
	skel := langCfg.Skeleton
	cst := langCfg.Cst

	if inList(cst.ImportTypes, ntype) {
		if p.includeImports && p.symbol == "" && p.nodeFullyInLineFilter(node) {
			p.emitFirst(lines, node, indent, strings.TrimSpace(p.NodeText(node)))
		}
		return
	}

	if inList(skel.SymbolTypes, ntype) {
		if inList(cst.DecoratorScopeTypes, ntype) {
			p.dfsDecorated(node, lines, depth, indent)
			return
		}
		if inList(cst.FunctionScopeTypes, ntype) {
			p.dfsFunction(node, lines, depth, indent, ntype)
			return
		}
		if inList(cst.ClassScopeTypes, ntype) {
			p.dfsClass(node, lines, depth, indent)
			return
		}
		p.seen[node.ID()] = struct{}{}
		for i := 0; i < int(node.ChildCount()); i++ {
			p.dfs(node.Child(i), lines, depth)
		}
		return
	}

	if depth == 0 && inList(skel.TopLevelExprTypes, ntype) && p.symbol == "" {
		text := strings.TrimSpace(p.NodeText(node))
		if text != "" && p.topLevelExprHeuristic(text) && p.nodeFullyInLineFilter(node) {
			p.emitFirst(lines, node, "", text)
		}
		for i := 0; i < int(node.ChildCount()); i++ {
			p.dfs(node.Child(i), lines, depth)
		}
		return
	}

	p.seen[node.ID()] = struct{}{}
	for i := 0; i < int(node.ChildCount()); i++ {
		p.dfs(node.Child(i), lines, depth)
	}
}

func (p *skeletonParser) topLevelExprHeuristic(text string) bool {
	if strings.HasPrefix(text, `"""`) || strings.HasPrefix(text, "'''") {
		return true
	}
	if strings.HasPrefix(text, `"`) || strings.HasPrefix(text, "'") {
		return true
	}
	return strings.Contains(text, "=")
}

func (p *skeletonParser) dfsDecorated(node *sitter.Node, lines *[]string, depth int, indent string) {
	langCfg := p.ctx.LangCfg
	inner := p.decoratedInner(node)
	if inner != nil && !p.matchesSymbol(inner, depth) {
		p.seen[node.ID()] = struct{}{}
		return
	}
	if inner != nil && p.wouldEmitDecl(inner, depth) {
		p.emitPathHeader(lines, inner)
	}
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)
		if child == nil {
			continue
		}
		if child.Type() == "decorator" {
			if inner != nil && p.isSymbolTarget(inner, depth) && p.nodeFullyInLineFilter(child) {
				p.emitCont(lines, indent, strings.TrimSpace(p.NodeText(child)))
			}
			p.seen[child.ID()] = struct{}{}
			continue
		}
		if inList(langCfg.Cst.FunctionScopeTypes, child.Type()) ||
			inList(langCfg.Cst.ClassScopeTypes, child.Type()) {
			if inner != nil && p.wouldEmitDecl(inner, depth) {
				p.omitNextSkeletonPathLine = true
			}
			p.dfs(child, lines, depth)
			break
		}
	}
	p.seen[node.ID()] = struct{}{}
}

func (p *skeletonParser) dfsFunction(node *sitter.Node, lines *[]string, depth int, indent, ntype string) {
	if ntype == "arrow_function" && p.isCallArgument(node) {
		p.seen[node.ID()] = struct{}{}
		return
	}
	if !p.matchesSymbol(node, depth) {
		p.seen[node.ID()] = struct{}{}
		return
	}
	if p.wouldEmitDecl(node, depth) {
		p.emitSignatureBlock(lines, node, indent, true)
	}
	p.seen[node.ID()] = struct{}{}
	for i := 0; i < int(node.ChildCount()); i++ {
		p.dfs(node.Child(i), lines, depth+1)
	}
}

func (p *skeletonParser) dfsClass(node *sitter.Node, lines *[]string, depth int, indent string) {
	if !p.matchesSymbol(node, depth) {
		p.seen[node.ID()] = struct{}{}
		return
	}
	if p.wouldEmitDecl(node, depth) {
		p.emitSignatureBlock(lines, node, indent, false)
	}
	p.seen[node.ID()] = struct{}{}
	for i := 0; i < int(node.ChildCount()); i++ {
		p.dfs(node.Child(i), lines, depth+1)
	}
}

func inList(set []string, s string) bool {
	for _, x := range set {
		if x == s {
			return true
		}
	}
	return false
}
