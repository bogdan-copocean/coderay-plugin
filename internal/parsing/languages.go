package parsing

import (
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	tsx "github.com/smacker/go-tree-sitter/typescript/tsx"
	tsts "github.com/smacker/go-tree-sitter/typescript/typescript"
)

// CstDispatchConfig mirrors coderay Python CstDispatchConfig.
type CstDispatchConfig struct {
	ImportTypes           []string
	FunctionScopeTypes    []string
	ClassScopeTypes       []string
	DecoratorScopeTypes   []string
	CallTypes             []string
	AssignmentTypes       []string
	DecoratorTypes        []string
	WithTypes             []string
	ClassBodyTypes        []string
	TypedParamTypes       []string
}

// SkeletonConfig holds skeleton emission rules per language.
type SkeletonConfig struct {
	SymbolTypes         []string
	DocstringExprType   string
	TopLevelExprTypes   []string
	BodyBlockTypes      []string
}

// LanguageConfig is the full per-language setup for parsing.
type LanguageConfig struct {
	Name       string
	Extensions []string
	Cst        CstDispatchConfig
	Skeleton   SkeletonConfig
}

var languageRegistry = map[string]*LanguageConfig{
	"python":    pythonConfig(),
	"javascript": javascriptConfig(),
	"typescript": typescriptConfig(),
}

var extensionToLang map[string]string

func init() {
	extensionToLang = make(map[string]string)
	for name, cfg := range languageRegistry {
		for _, ext := range cfg.Extensions {
			extensionToLang[strings.ToLower(ext)] = name
		}
	}
}

// GetLanguageForFile returns config for path extension, or nil if unsupported.
func GetLanguageForFile(path string) *LanguageConfig {
	ext := strings.ToLower(filepath.Ext(path))
	langName := extensionToLang[ext]
	if langName == "" {
		return nil
	}
	return languageRegistry[langName]
}

func pythonConfig() *LanguageConfig {
	return &LanguageConfig{
		Name:       "python",
		Extensions: []string{".py", ".pyi"},
		Cst: CstDispatchConfig{
			ImportTypes: []string{
				"import_statement",
				"import_from_statement",
				"future_import_statement",
			},
			FunctionScopeTypes:  []string{"function_definition"},
			ClassScopeTypes:     []string{"class_definition"},
			DecoratorScopeTypes: []string{"decorated_definition"},
			CallTypes:           []string{"call"},
			AssignmentTypes:   []string{"assignment"},
			DecoratorTypes:      []string{"decorator"},
			WithTypes:           []string{"with_statement"},
			ClassBodyTypes:      []string{"block"},
			TypedParamTypes:     []string{"typed_parameter"},
		},
		Skeleton: SkeletonConfig{
			SymbolTypes: []string{
				"function_definition",
				"class_definition",
				"decorated_definition",
			},
			DocstringExprType: "expression_statement",
			TopLevelExprTypes: []string{"expression_statement"},
			BodyBlockTypes:    []string{"block"},
		},
	}
}

func javascriptConfig() *LanguageConfig {
	return &LanguageConfig{
		Name:       "javascript",
		Extensions: []string{".js", ".jsx", ".mjs", ".cjs"},
		Cst:        jsTsCst(),
		Skeleton:   jsTsSkeleton(),
	}
}

func typescriptConfig() *LanguageConfig {
	return &LanguageConfig{
		Name:       "typescript",
		Extensions: []string{".ts", ".tsx"},
		Cst:        jsTsCst(),
		Skeleton:   jsTsSkeleton(),
	}
}

func jsTsCst() CstDispatchConfig {
	return CstDispatchConfig{
		ImportTypes: []string{"import_statement"},
		FunctionScopeTypes: []string{
			"function_declaration",
			"method_definition",
			"arrow_function",
		},
		ClassScopeTypes: []string{
			"class_declaration",
			"interface_declaration",
			"type_alias_declaration",
			"type_declaration",
		},
		DecoratorScopeTypes: []string{},
		CallTypes:           []string{"call_expression"},
		AssignmentTypes:     []string{"assignment_expression", "variable_declarator"},
		DecoratorTypes:      []string{},
		WithTypes:           []string{},
		ClassBodyTypes:      []string{"block", "class_body"},
		TypedParamTypes: []string{
			"typed_parameter",
			"required_parameter",
			"optional_parameter",
		},
	}
}

func jsTsSkeleton() SkeletonConfig {
	return SkeletonConfig{
		SymbolTypes: []string{
			"function_declaration",
			"class_declaration",
			"method_definition",
			"arrow_function",
			"interface_declaration",
			"type_alias_declaration",
			"type_declaration",
		},
		DocstringExprType: "expression_statement",
		TopLevelExprTypes: []string{"expression_statement", "lexical_declaration"},
		BodyBlockTypes:    []string{"statement_block"},
	}
}

// LanguageForPath returns the tree-sitter language handle for a file path.
func LanguageForPath(path string) *sitter.Language {
	cfg := GetLanguageForFile(path)
	if cfg == nil {
		return nil
	}
	ext := strings.ToLower(filepath.Ext(path))
	switch cfg.Name {
	case "python":
		return python.GetLanguage()
	case "javascript":
		return javascript.GetLanguage()
	case "typescript":
		if ext == ".tsx" {
			return tsx.GetLanguage()
		}
		return tsts.GetLanguage()
	default:
		return nil
	}
}
