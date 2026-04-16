# coderay-skeleton

Agents tend to read too much code by default — sometimes guessing how much, other times swallowing whole files — most of it noise.

A **skeleton** is a stripped view of a source file — class and function signatures, docstrings, and top-level assignments, with all bodies replaced by `...`. It tells you *what* a file contains and *exactly where* (start and end line), without the noise of implementation detail.

This plugin enforces **skeleton before almost all reads**: run `coderay-skeleton` on a file first, get signatures and exact line ranges, then read only that slice.

|                    |                                                                           |
|--------------------|---------------------------------------------------------------------------|
| **Pre-requisites** | None. Works on any supported file immediately.                            |
| **Output**         | Signatures, docstrings, top-level assignments; bodies replaced with `...` |
| **Scope**          | Optional symbol filter, optional line range window                        |

## Why not just use grep?

`grep` finds where a symbol *starts*. It doesn't know where it *ends*. Without the end line, an agent has to guess a window size or read the whole file.

`coderay-skeleton` gives both — precise start **and** end for every function and class. That one number is the difference between a 10-line read and a 500-line dump.

```bash
# grep finds the start — but where does the function end?
grep -n "def process_event" src/handler.py
# → 142:def process_event(self, event: Event) -> None:
# Read 20 lines? 80? The whole file?

# skeleton gives the full span immediately
coderay-skeleton --file src/handler.py --symbol process_event
# → src/handler.py:142-187
#   def process_event(self, event: Event) -> None:
#       """Handle incoming events and dispatch to registered processors."""
#       ...
# Read exactly lines 142-187. Done.
```

## Agent integration

The skill lives in [`SKILL.md`](./SKILL.md). Integrate with your agent:

### Claude Code (plugin — recommended)

```shell
claude plugin marketplace add bogdan-copocean/coderay-plugin && claude plugin install coderay@coderay-plugin
```

Binaries for macOS arm64 and Linux amd64/arm64 are published automatically on each release via GitHub Actions. On first use the wrapper downloads the right binary from `releases/latest` and caches it. No manual setup needed on supported platforms.

The `skeleton` skill auto-activates. `/coderay:help` shows a quick-reference.

### Other agents — manual setup (for now)

Native plugins for Cursor, Codex, Copilot, and Windsurf are planned. Until then, wire the skill manually:

| Agent | What to do |
|-------|------------|
| Claude Code (standalone) | Build binary, add to `PATH`, append `SKILL.md` to `CLAUDE.md` |
| Codex / OpenAI Agents | Append `SKILL.md` to `AGENTS.md` |
| Cursor | Copy `SKILL.md` to `.cursor/rules/coderay-skeleton.mdc` |
| Copilot | Append `SKILL.md` to `.github/copilot-instructions.md` |
| Windsurf | Append `SKILL.md` to `.windsurfrules` |

The binary must be on `PATH` for the skill to have effect. Build it with `make build` and move it to a directory in your `PATH`.

---

## Output format

```
/abs/path/to/file.py:38-53
class BaseService:
    """Base service with a simple interface."""
/abs/path/to/file.py:43-45
    def __init__(self, root: Path) -> None:
        """Initialize the base service with a root path."""
        ...
```

## CLI flags

| Flag | Meaning |
|------|---------|
| `--file` | Source file path; optional `:START-END` line window suffix |
| `--symbol` | Filter to one class or function; dotted for nested (`Class.method`) |
| `--include-imports` | Include import statements (default: off) |
| `--file-line-range` | Line window as separate flag — mutually exclusive with `:range` suffix |

**Supported languages:** Python (`.py`, `.pyi`), JavaScript (`.js`, `.jsx`, `.mjs`, `.cjs`), TypeScript (`.ts`, `.tsx`). Unsupported extensions return raw content unchanged.

## Build

Requires Go + CGO (tree-sitter uses C):

```bash
make build                        # builds plugin/bin/coderay-skeleton-{os}-{arch}
make test
```

## Relation to CodeRay

The full Python project at [bogdan-copocean/coderay](https://github.com/bogdan-copocean/coderay) provides an MCP server with `get_file_skeleton`, `semantic_search`, and `get_impact_radius`. This Go module implements **skeleton only** — no Python runtime, no index, no config required.
