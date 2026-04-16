package skeleton

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadFileSkeleton reads a source file from disk and runs ExtractSkeleton.
// filePathArg may include an optional :START-END suffix (same as ParseSkeletonFileArg).
// The filesystem path is resolved with filepath.Abs (relative paths use the process working directory).
// UTF-8 bytes are read with invalid sequences replaced (U+FFFD), matching the CodeRay MCP behavior.
func ReadFileSkeleton(filePathArg string, opts Options, fileLineRangeStr *string) (string, error) {
	pathStr, rngSuffix, err := ParseSkeletonFileArg(filePathArg, true)
	if err != nil {
		return "", err
	}
	var lineRange *LineRange
	if rngSuffix != nil {
		lineRange = rngSuffix
	}
	if fileLineRangeStr != nil && *fileLineRangeStr != "" {
		if lineRange != nil {
			return "", fmt.Errorf("Use either file_path :START-END suffix or file_line_range, not both.")
		}
		lr, err := ParseFileLineRange(*fileLineRangeStr)
		if err != nil {
			return "", err
		}
		lineRange = &lr
	}

	candAbs, err := filepath.Abs(filepath.Clean(filepath.FromSlash(pathStr)))
	if err != nil {
		return "", err
	}

	st, err := os.Stat(candAbs)
	if err != nil {
		return "", fmt.Errorf("File not found: %s", pathStr)
	}
	if st.IsDir() {
		return "", fmt.Errorf("File not found: %s", pathStr)
	}

	raw, err := os.ReadFile(candAbs)
	if err != nil {
		return "", err
	}
	content := strings.ToValidUTF8(string(raw), "\uFFFD")

	opts2 := opts
	if lineRange != nil {
		opts2.LineRange = lineRange
	}
	return ExtractSkeleton(candAbs, content, opts2)
}
