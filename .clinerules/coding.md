# Coding Standards — progas-wms-be

Scope: all Go code in module `progas-wms-be` (Go 1.26, Fiber v3, GORM + MySQL).
These rules are always active. When unsure, mirror the closest existing file in the same layer.

## 1. Golden rules
- `go build ./...` and `go vet ./...` MUST pass before a task is considered done.
- Format with `gofmt`/`goimports` only (tabs, no manual whitespace alignment): `go fmt ./...`.
- Business rules live ONLY in `usecase/` (or pure helpers in `helper/`). `handler/` and `repository/` stay thin.
- Never open a DB handle or write raw SQL in `handler/`, `mapper/`, or `usecase/`; go through a repository.
- Reuse `constant.*`, `enum.*`, `helper.*`, `global.*` before inventing new ones.
- No new third-party dependency without a strong reason; the stack is fixed in `go.mod`.

## 2. Layering and allowed dependencies
| Layer | May import | Must NOT import |
|---|---|---|
| `handler/` | `dto`, `usecase`, `global`, `helper`, fiber | repositories, models, gorm |
| `usecase/` | `dto`, `model`, `mapper`, `repository`, `enum`, `helper`, `constant`, `global` | fiber, gorm directly |
| `repository/` | `model`, `global`, `helper`, `enum`, gorm | usecase, dto, handler |
| `mapper/` | `model`, `dto`, `enum` | repositories, gorm |
| `model/` | `enum` | app-specific packages |

Keep the dependency direction one-way: `handler → usecase → repository → model`.

## 3. Naming
- IDs use `Id`, never `ID`: `UserId`, `RoleId`, `CreatedBy`, `FindById`, `FindByBarcode`.
- Repository verbs match existing style: `FindAll`, `FindById`, `FindBy...`, `Create`, `Update`, `Delete`, `Count...`.
- Constructors: `NewXxxRepository` / `NewXxxUsecase` / `NewXxxHandler`. Interfaces are named for the role (`CylinderRepository`, `DeliveryOrderUsecase`).
- DTOs: `XxxRequest` / `XxxResponse`; list queries use the shared `dto.ListQuery`.
- Booleans read as predicates: `IsActive`, `HasSearch`, `CanIssueOnDO`.

## 4. Errors
- Use `global.ErrorResponse` as the LAST return value of repository/usecase methods. Handlers return `error`.
- Map failures to the right constructor:
  - validation / domain rule: `global.BadRequestError(msg)` or `global.BadRequestErrorWithData(msg, data, constant.ErrValidationError)`
  - missing record: `global.NotFoundError("X not found")`
  - unexpected / infra: `global.InternalServerError(err)` (never surface raw infra errors as 4xx)
  - auth: `global.ForbiddenError()` / `global.UnauthorizedError()`
- Do not `panic` for expected failures (only `PanicHandler` recovers panics).
- Return `nil` error response on success.

## 5. DTOs, validation, Swagger
- Add `json:"..."` tags to every field; add `validate:"..."` tags to request fields.
- Validate in handlers via `helper.ValidateBody / ValidateQuery / ValidateParam / ValidateHeader`.
- Use `validate:"required"`, `min=1`, and `dive` for slices.
- Every exported handler method gets swaggo annotations (`@Summary`, `@Tags`, `@Param`, `@Success`, `@Router`) and answers with `global.CreateResponse(res, fiber.StatusOK, c)`.
- For list responses in Swagger use the concrete `dto.PaginatedXxxList` types (swag cannot render nested generics).

## 6. Models
- Embed `model.BaseModel` (`Id` varchar(36) PK, timestamps, soft delete, `CreatedBy/UpdatedBy/DeletedBy`). IDs are auto-generated UUID v7 in `BeforeCreate` — never set them manually.
- Use `enum.*` types for status columns with `gorm:"type:varchar(N)"`.
- GORM uses `SingularTable: true` (`config/database.go`), so tables are singular (`Customer`, `SalesOrder`). In raw SQL, backtick PascalCase table names, e.g. `` `Customer`.name ``.
- Register every new model in `config/migration.go` (`AutoMigrate`).

## 7. Enums
- String-based types in `enum/`, with `IsValid()`, `Scan()`, `Value()`.
- When adding a value: add the const, extend `IsValid()`, and keep `Scan`/`Value` round-tripping.

## 8. Constants (never hardcode)
- Permission keys: `constant.Perm*` (`constant/permission.go`).
- Audit actions/objects and ledger/movement actions: `constant.Audit*` / `constant.Ledger*` / `constant.Movement*` (`constant/audit.go`).
- Role names: `constant/role.go`. Env keys: `constant/constant.go`.

## 9. Transactions, logging, audit
- Multi-table writes use `helper.TxManager`: `tx := u.txManager.New()`, `defer tx.CheckPanic()`, `tx.Rollback()` on EVERY early return, `tx.Commit()` last.
- Repository mutators accept `tx helper.Tx` and resolve the handle via `dbFromTx(tx)`.
- Audit logging is best-effort and runs AFTER a successful commit: `_ = u.auditLogRepo.Log(actor, action, object, id, map[string]any{...})`. Never inside the transaction; never fail the request because audit failed.
- Never log secrets, tokens, or password hashes.

## 10. Style
- Keep functions small and single-purpose; extract pure logic into `helper/` (e.g. `helper/finance.go`, `helper/outbound.go`).
- Prefer guard clauses and early returns over deep nesting.
- No dead or commented-out code; no unused imports.
- Comments explain WHY, not WHAT.

## 11. Commits
- Conventional Commits as used historically: `feat: ...`, `fix: ...`, `chore: ...`, `refactor: ...`, `docs: ...`.
