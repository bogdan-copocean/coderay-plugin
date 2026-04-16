package skeleton

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func testdataSkeletonDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(file), "..", "testdata", "skeleton")
}

func loadGolden(t *testing.T, fixture string, canonicalPath string) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(testdataSkeletonDir(t), "fixtures", fixture))
	if err != nil {
		t.Fatal(err)
	}
	return strings.ReplaceAll(string(raw), "__CANONICAL_PATH__", canonicalPath)
}

func TestGoldenPython(t *testing.T) {
	dir := testdataSkeletonDir(t)
	py := filepath.Join(dir, "canonical_concepts.py")
	content, err := os.ReadFile(py)
	if err != nil {
		t.Fatal(err)
	}
	abs, err := filepath.Abs(py)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("with_imports", func(t *testing.T) {
		got, err := ExtractSkeleton(py, string(content), Options{IncludeImports: true})
		if err != nil {
			t.Fatal(err)
		}
		want := loadGolden(t, "py_with_imports.expected", abs)
		if got != want {
			t.Fatal(diffHead(got, want))
		}
	})
	t.Run("without_imports", func(t *testing.T) {
		got, err := ExtractSkeleton(py, string(content), Options{IncludeImports: false})
		if err != nil {
			t.Fatal(err)
		}
		want := loadGolden(t, "py_without_imports.expected", abs)
		if got != want {
			t.Fatalf("mismatch without_imports:\n%s", diffHead(got, want))
		}
	})
	t.Run("unsupported_ext", func(t *testing.T) {
		_, c := string(content), string(content)
		for _, p := range []string{"test.xyz", "noext"} {
			got, err := ExtractSkeleton(p, c, Options{})
			if err != nil {
				t.Fatal(err)
			}
			if got != c {
				t.Fatalf("path %s: expected raw content", p)
			}
		}
	})
	t.Run("symbol_Repository", func(t *testing.T) {
		got, err := ExtractSkeleton(py, string(content), Options{Symbol: "Repository"})
		if err != nil {
			t.Fatal(err)
		}
		want := loadGolden(t, "py_symbol_Repository.expected", abs)
		if got != want {
			t.Fatal(diffHead(got, want))
		}
	})
	t.Run("symbol_Repository_get", func(t *testing.T) {
		got, err := ExtractSkeleton(py, string(content), Options{Symbol: "Repository.get"})
		if err != nil {
			t.Fatal(err)
		}
		want := loadGolden(t, "py_symbol_Repository_get.expected", abs)
		if got != want {
			t.Fatal(diffHead(got, want))
		}
	})
	t.Run("symbol_local_imports", func(t *testing.T) {
		got, err := ExtractSkeleton(py, string(content), Options{Symbol: "local_imports_example"})
		if err != nil {
			t.Fatal(err)
		}
		want := loadGolden(t, "py_symbol_local_imports_example.expected", abs)
		if got != want {
			t.Fatal(diffHead(got, want))
		}
	})
	t.Run("symbol_decorator", func(t *testing.T) {
		got, err := ExtractSkeleton(py, string(content), Options{Symbol: "decorator"})
		if err != nil {
			t.Fatal(err)
		}
		want := loadGolden(t, "py_symbol_decorator.expected", abs)
		if got != want {
			t.Fatal(diffHead(got, want))
		}
	})
	t.Run("unknown_symbol", func(t *testing.T) {
		for _, sym := range []string{"DoesNotExist", "FakeClass.method"} {
			got, err := ExtractSkeleton(py, string(content), Options{Symbol: sym})
			if err != nil {
				t.Fatal(err)
			}
			prefix := "# Symbol '" + sym + "' not found. Available symbols: "
			if !strings.HasPrefix(got, prefix) {
				t.Fatalf("sym %s: %q", sym, truncate(got, 200))
			}
		}
	})
}

func TestGoldenTypeScript(t *testing.T) {
	dir := testdataSkeletonDir(t)
	ts := filepath.Join(dir, "canonical_concepts.ts")
	content, err := os.ReadFile(ts)
	if err != nil {
		t.Fatal(err)
	}
	abs, err := filepath.Abs(ts)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("with_imports", func(t *testing.T) {
		got, err := ExtractSkeleton(ts, string(content), Options{IncludeImports: true})
		if err != nil {
			t.Fatal(err)
		}
		want := loadGolden(t, "ts_with_imports.expected", abs)
		if got != want {
			t.Fatal(diffHead(got, want))
		}
	})
	t.Run("without_imports", func(t *testing.T) {
		got, err := ExtractSkeleton(ts, string(content), Options{IncludeImports: false})
		if err != nil {
			t.Fatal(err)
		}
		want := loadGolden(t, "ts_without_imports.expected", abs)
		if got != want {
			t.Fatal(diffHead(got, want))
		}
	})
	t.Run("symbol_CoreService", func(t *testing.T) {
		got, err := ExtractSkeleton(ts, string(content), Options{Symbol: "CoreService"})
		if err != nil {
			t.Fatal(err)
		}
		want := loadGolden(t, "ts_symbol_CoreService.expected", abs)
		if got != want {
			t.Fatal(diffHead(got, want))
		}
	})
	t.Run("symbol_withClosure", func(t *testing.T) {
		got, err := ExtractSkeleton(ts, string(content), Options{Symbol: "CoreService.withClosure"})
		if err != nil {
			t.Fatal(err)
		}
		want := loadGolden(t, "ts_symbol_CoreService_withClosure.expected", abs)
		if got != want {
			t.Fatal(diffHead(got, want))
		}
	})
	t.Run("symbol_buildProfileLabel", func(t *testing.T) {
		got, err := ExtractSkeleton(ts, string(content), Options{Symbol: "buildProfileLabel"})
		if err != nil {
			t.Fatal(err)
		}
		want := loadGolden(t, "ts_symbol_buildProfileLabel.expected", abs)
		if got != want {
			t.Fatal(diffHead(got, want))
		}
	})
	t.Run("unknown", func(t *testing.T) {
		got, err := ExtractSkeleton(ts, string(content), Options{Symbol: "DoesNotExist"})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasPrefix(got, "# Symbol 'DoesNotExist' not found. Available symbols: ") {
			t.Fatal(got)
		}
	})
}

func TestSkeletonLineRangeColumn(t *testing.T) {
	dir := testdataSkeletonDir(t)
	py := filepath.Join(dir, "canonical_concepts.py")
	content, err := os.ReadFile(py)
	if err != nil {
		t.Fatal(err)
	}
	abs, err := filepath.Abs(py)
	if err != nil {
		t.Fatal(err)
	}
	sk, err := ExtractSkeleton(py, string(content), Options{IncludeImports: false})
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(sk, "\n")
	var i int
	for i = range lines {
		if strings.Contains(lines[i], "def top_level_helper") {
			break
		}
	}
	if i == 0 || i >= len(lines) {
		t.Fatal("top_level_helper not found")
	}
	prev := strings.TrimSpace(lines[i-1])
	if !strings.Contains(prev, abs) || !strings.Contains(prev, ":27-29") {
		t.Fatalf("bad path line: %q", prev)
	}
	if strings.HasPrefix(strings.TrimLeft(lines[i], " "), "27-") {
		t.Fatal("line should not start with bare range")
	}
}

func TestFileLineRangeEmptyMessage(t *testing.T) {
	dir := testdataSkeletonDir(t)
	py := filepath.Join(dir, "canonical_concepts.py")
	content, err := os.ReadFile(py)
	if err != nil {
		t.Fatal(err)
	}
	lr := LineRange{Start: 99998, End: 99999}
	got, err := ExtractSkeleton(py, string(content), Options{LineRange: &lr})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "file line range 99998-99999") {
		t.Fatal(got)
	}
}

func TestFileLineRangeWindow(t *testing.T) {
	dir := testdataSkeletonDir(t)
	py := filepath.Join(dir, "canonical_concepts.py")
	content, err := os.ReadFile(py)
	if err != nil {
		t.Fatal(err)
	}
	lr := LineRange{Start: 27, End: 35}
	got, err := ExtractSkeleton(py, string(content), Options{IncludeImports: false, LineRange: &lr})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "def top_level_helper") {
		t.Fatal(got)
	}
	if !strings.Contains(got, "async def async_helper") {
		t.Fatal(got)
	}
	if strings.Contains(got, "class BaseService") {
		t.Fatal("should not include BaseService")
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func diffHead(got, want string) string {
	g := strings.Split(got, "\n")
	w := strings.Split(want, "\n")
	n := len(g)
	if len(w) < n {
		n = len(w)
	}
	for i := 0; i < n; i++ {
		if g[i] != w[i] {
			return fmt.Sprintf("first diff line %d:\nwant: %s\ngot:  %s", i+1, w[i], g[i])
		}
	}
	if len(g) != len(w) {
		return fmt.Sprintf("length mismatch: got %d lines want %d", len(g), len(w))
	}
	return "unknown mismatch"
}
