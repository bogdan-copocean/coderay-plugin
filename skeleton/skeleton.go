package skeleton

import (
	"github.com/bogdan-copocean/coderay-plugin/internal/parsing"
)

// ErrInvalidLineRange is returned when LineRange bounds are invalid.
var ErrInvalidLineRange = parsing.ErrInvalidLineRange

// Options configures ExtractSkeleton (MCP-aligned defaults: IncludeImports false).
type Options struct {
	IncludeImports bool
	Symbol         string
	LineRange      *LineRange
}

// ExtractSkeleton extracts class/function signatures and docstrings without bodies.
func ExtractSkeleton(path, content string, opts Options) (string, error) {
	opt := parsing.ExtractOptions{
		IncludeImports: opts.IncludeImports,
		Symbol:         opts.Symbol,
	}
	if opts.LineRange != nil {
		opt.LineRange = &[2]int{opts.LineRange.Start, opts.LineRange.End}
	}
	return parsing.ExtractSkeleton(path, content, opt)
}
