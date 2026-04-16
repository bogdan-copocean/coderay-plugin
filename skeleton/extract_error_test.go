package skeleton

import (
	"errors"
	"testing"
)

func TestExtractSkeletonInvalidLineRange(t *testing.T) {
	lr := LineRange{Start: 2, End: 1}
	_, err := ExtractSkeleton("x.py", "a=1", Options{LineRange: &lr})
	if !errors.Is(err, ErrInvalidLineRange) {
		t.Fatalf("got %v", err)
	}
}
