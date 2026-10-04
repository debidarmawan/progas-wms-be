# progas-wms-be — Agent Notes

Backend API for Progas WMS (cylinder gas warehouse management). Go 1.26, Fiber v3, GORM/MySQL.

## Project rules (read before coding)

The authoritative coding standards live in `.clinerules/` (Cline reads this directory natively).
They are mirrored under `.agent/` for other AI agents/tools:

- `.clinerules/coding.md` — Coding standards (layering, naming, errors, DTOs, models, enums, constants, transactions).
- `.clinerules/testing.md` — Test requirements (what to test, patterns, definition of done).
- `.clinerules/architecture.md` — Structural decisions (layering, wiring, auth/RBAC, DB, domain invariants, extension checklist).

When adding or modifying code, follow these rules. When you change a rule, edit the file in `.clinerules/` (the `.agent/` copies are symlinks to it).

**Read `../progas-docs/` before making non-trivial changes.** It contains the full architecture, module inventory, and business flow so you don't need to explore this entire codebase to get oriented:

- `../progas-docs/01-architecture.md` — layering (handler → usecase → repository → model), transaction pattern, RBAC, auth
- `../progas-docs/02-modules-api.md` — every module, model fields, and API endpoint with required permission
- `../progas-docs/03-business-flow.md` — current (as-is) business flow, including the cylinder status lifecycle
- `../progas-docs/04-roadmap.md` — planned PO → SO → DO → Trip flow (NOT yet implemented)

## Critical facts (avoid re-deriving these by exploring)

- Sales flow today: Delivery Order (DO) is created directly from a list of cylinder barcodes — there is NO Purchase Order or Sales Order yet. Don't assume PO/SO models exist.
- Multi-table mutations use `helper.TxManager` (explicit transaction, manual `tx.Rollback()` on every early return, `tx.Commit()` at the end).
- Every write usecase logs to `AuditLogRepository` after commit succeeds (best-effort, outside the transaction).
- Permissions are enforced per route group in `server/routes.go` via `constant.Perm*` keys — a new endpoint needs a new permission key plus an RBAC seed entry.
- Cylinder status is the core domain invariant: `EMPTY → READY_TO_FILL → FILLED → READY → IN_TRANSIT → OUTSTANDING → EMPTY` (loop), with `MAINTENANCE`/`LOST`/`WRITE_OFF` as side states.

## Business flow changes

Any change to sales, delivery, or cylinder-status business logic MUST be reflected in `../progas-docs/03-business-flow.md` and logged in `../progas-docs/CHANGELOG.md`.
