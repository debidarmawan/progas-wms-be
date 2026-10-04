# .agent — cross-tool rule mirror

This folder exposes the project's AI rules to agents/tools that do not read `.clinerules/`.

The files here are **symlinks** to the canonical rules in `../.clinerules/`:

- `coding.md` → `../.clinerules/coding.md`
- `testing.md` → `../.clinerules/testing.md`
- `architecture.md` → `../.clinerules/architecture.md`

## Source of truth
Edit the files in `.clinerules/` — never edit the copies here (they are the same files via symlink).

## Which tool reads what
| Tool | Rule source |
|---|---|
| Cline | `.clinerules/` (read natively) |
| Cursor, Codex, Windsurf, Gemini, others | `AGENTS.md` (repo root), which links to these rules |
| Claude Code | `CLAUDE.md` → `@AGENTS.md` |

`AGENTS.md` at the repo root is the cross-tool entry point and points to all three rule files.
