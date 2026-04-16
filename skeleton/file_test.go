package skeleton

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadFileSkeleton(t *testing.T) {
	tmp := t.TempDir()
	py := filepath.Join(tmp, "hello.py")
	if err := os.WriteFile(py, []byte("def greet(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	absPy, err := filepath.Abs(py)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("missing_file", func(t *testing.T) {
		_, err := ReadFileSkeleton(filepath.Join(tmp, "nope.py"), Options{}, nil)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("real_file", func(t *testing.T) {
		out, err := ReadFileSkeleton(py, Options{}, nil)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "def greet") || !strings.Contains(out, absPy) {
			t.Fatal(out)
		}
	})

	t.Run("file_line_range", func(t *testing.T) {
		s := "1-1"
		out, err := ReadFileSkeleton(py, Options{}, &s)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "def greet") {
			t.Fatal(out)
		}
	})

	t.Run("path_suffix", func(t *testing.T) {
		out, err := ReadFileSkeleton(py+":1-1", Options{}, nil)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "def greet") || !strings.Contains(out, absPy) {
			t.Fatal(out)
		}
	})

	t.Run("dual_range", func(t *testing.T) {
		s := "2-2"
		_, err := ReadFileSkeleton(py+":1-1", Options{}, &s)
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
