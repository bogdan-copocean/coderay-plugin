---
name: help
description: Quick-reference for coderay-skeleton CLI. Trigger when user asks "how do I use coderay", "coderay help", "what does skeleton do", or /coderay:help.
argument-hint: [topic]
allowed-tools: [Bash]
---

# coderay-skeleton

Returns signatures, docstrings, and top-level assignments for a source file — bodies replaced with `...`. Each entry includes a `path:start-end` line range so you read only the span you need.

**No index. No config. No network.**

## Usage

```bash
coderay-skeleton --file path/to/file.py
coderay-skeleton --file src/app.ts --symbol MyClass.method
coderay-skeleton --file foo.py:10-80
coderay-skeleton --file src/app.ts --include-imports
```

# Flags

| Flag | When to use |
|------|-------------|
| `--file PATH` | Required. Accepts optional `:START-END` suffix to pre-narrow to a line window. |
| `--symbol NAME` | Filter to one declaration. Dotted for nested (`Class.method`). |
| `--include-imports` | Include import statements (default: off). |
| `--file-line-range START-END` | Line window as separate flag — mutually exclusive with `:range` suffix. |

Full docs and flags: https://github.com/bogdan-copocean/coderay-plugin

Run `coderay-skeleton --help` for a flag summary.
