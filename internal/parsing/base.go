package parsing

import (
	sitter "github.com/smacker/go-tree-sitter"
)

// BaseTreeSitterParser parses once and exposes node text helpers.
type BaseTreeSitterParser struct {
	ctx          *ParseContext
	sourceBytes  []byte
	parser       *sitter.Parser
	tree         *sitter.Tree
}

// NewBaseTreeSitterParser constructs a parser for the given context.
func NewBaseTreeSitterParser(ctx *ParseContext) *BaseTreeSitterParser {
	return &BaseTreeSitterParser{
		ctx:         ctx,
		sourceBytes: []byte(ctx.Content),
	}
}

// GetTree parses (cached) and returns the syntax tree.
func (b *BaseTreeSitterParser) GetTree() *sitter.Tree {
	if b.tree != nil {
		return b.tree
	}
	lang := LanguageForPath(b.ctx.FilePath)
	if lang == nil {
		return nil
	}
	if b.parser == nil {
		b.parser = sitter.NewParser()
		b.parser.SetLanguage(lang)
	}
	b.tree = b.parser.Parse(nil, b.sourceBytes)
	return b.tree
}

// NodeText returns source text for a node.
func (b *BaseTreeSitterParser) NodeText(n *sitter.Node) string {
	if n == nil || n.IsNull() {
		return ""
	}
	return n.Content(b.sourceBytes)
}
