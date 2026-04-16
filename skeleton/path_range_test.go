package skeleton

import (
	"testing"
)

func TestParseSkeletonFileArg(t *testing.T) {
	t.Parallel()
	base, r, err := ParseSkeletonFileArg("foo.py", true)
	if err != nil || base != "foo.py" || r != nil {
		t.Fatalf("got %q %v %v", base, r, err)
	}

	base, r, err = ParseSkeletonFileArg("foo.py:1-10", true)
	if err != nil || base != "foo.py" || r.Start != 1 || r.End != 10 {
		t.Fatalf("got %q %+v %v", base, r, err)
	}

	_, _, err = ParseSkeletonFileArg(":1-2", true)
	if err == nil {
		t.Fatal("expected error for empty base")
	}

	_, _, err = ParseSkeletonFileArg("a.py:3-1", true)
	if err == nil {
		t.Fatal("expected error end < start")
	}
}

func TestParseFileLineRange(t *testing.T) {
	t.Parallel()
	r, err := ParseFileLineRange("  1-10 ")
	if err != nil || r.Start != 1 || r.End != 10 {
		t.Fatalf("%+v %v", r, err)
	}
	_, err = ParseFileLineRange("bad")
	if err == nil {
		t.Fatal("expected error")
	}
	_, err = ParseFileLineRange("3-1")
	if err == nil {
		t.Fatal("expected error")
	}
}
