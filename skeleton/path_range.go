package skeleton

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// LineRange is a 1-based inclusive file line window (same as Python skeleton tests).
type LineRange struct {
	Start, End int
}

var suffixRe = regexp.MustCompile(`:(\d+)-(\d+)$`)

// ParseSkeletonFileArg returns the filesystem path and optional 1-based inclusive line range.
// If parseSuffix is false, returns path unchanged and no range.
func ParseSkeletonFileArg(path string, parseSuffix bool) (base string, lineRange *LineRange, err error) {
	if !parseSuffix {
		return path, nil, nil
	}
	sub := suffixRe.FindStringSubmatch(path)
	if sub == nil {
		return path, nil, nil
	}
	start, err1 := strconv.Atoi(sub[1])
	end, err2 := strconv.Atoi(sub[2])
	if err1 != nil || err2 != nil {
		return "", nil, fmt.Errorf("invalid file line range suffix")
	}
	if end < start {
		return "", nil, fmt.Errorf("file line range end must be >= start")
	}
	base = strings.TrimSuffix(path, sub[0])
	if base == "" {
		return "", nil, fmt.Errorf("empty path before file line range")
	}
	return base, &LineRange{Start: start, End: end}, nil
}

var fullLineRangeRe = regexp.MustCompile(`^(\d+)-(\d+)$`)

// ParseFileLineRange parses START-END (1-based inclusive).
func ParseFileLineRange(s string) (LineRange, error) {
	s = strings.TrimSpace(s)
	sub := fullLineRangeRe.FindStringSubmatch(s)
	if sub == nil {
		return LineRange{}, fmt.Errorf("expected file line range START-END")
	}
	start, _ := strconv.Atoi(sub[1])
	end, _ := strconv.Atoi(sub[2])
	if end < start {
		return LineRange{}, fmt.Errorf("file line range end must be >= start")
	}
	return LineRange{Start: start, End: end}, nil
}
