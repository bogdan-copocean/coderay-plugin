---
name: skeleton
description: >
  ALWAYS run `coderay-skeleton --file PATH` before reading any source file. NEVER do a
  full file read unless absolutely necessary. Activate whenever about to open a file for
  discovery, symbol lookup, or editing prep.
---

# Hard rules

**Rule 1 — skeleton before every file read.**
Before opening any source file, run `coderay-skeleton --file PATH`. Skip only when you
already have a skeleton-derived `path:start-end` for this file from earlier in this session.

**Rule 2 — read only the span skeleton gave you.**
Never read beyond the line range skeleton returned. The only exception is when no narrower
read can complete the task — not a tighter range, not `--symbol`, not another skeleton call.
That bar is high.

```
coderay-skeleton --file PATH  →  pick path:start-end  →  read only that span
```

# Reading skeleton output

Each entry is prefixed with an absolute path and 1-based inclusive line range:

```
/abs/path/to/file.py:38-53
class BaseService:
    """Base service with a simple interface."""
/abs/path/to/file.py:43-45
    def __init__(self, root: Path) -> None:
        """Initialize the base service with a root path."""
        ...
```

`path:start-end` is your read target. Use `--symbol NAME` to filter to one declaration.
Dotted path for nested: `--symbol Class.method`. If the symbol isn't found, skeleton lists
available names.

# Flags

| Flag | When to use |
|------|-------------|
| `--file PATH` | Required. Accepts optional `:START-END` suffix to pre-narrow to a line window. |
| `--symbol NAME` | Filter to one declaration. Dotted for nested (`Class.method`). |
| `--include-imports` | Include import statements (default: off). |
| `--file-line-range START-END` | Line window as separate flag — mutually exclusive with `:range` suffix. |

# Unsupported files

Unsupported extensions return raw content — still read only the span you need, never the whole file.

# When coderay-skeleton is unavailable — built-in fallback

The discipline stays the same — the tools change.

Use grep to locate the symbol, then read a bounded window from that line. Scan the window
for the closing delimiter to find the actual end. Read only that span. Never use "binary
not found" as a reason to read the whole file.
