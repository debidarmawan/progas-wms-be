# Architecture & Structural Decisions — progas-wms-be

System: REST API for Progas WMS (cylinder-gas warehouse management).
Stack: Go 1.26, Fiber v3, GORM + MySQL, JWT auth, RBAC.

## 1. Layers
```
handler/    HTTP boundary: bind + validate DTO, call usecase, write response
usecase/    business rules + orchestration + transactions
repository/ data access (GORM), one file per entity
model/      GORM structs (tables)
dto/        request/response contracts
mapper/     model -> dto conversion
enum/       strongly-typed statuses/values
helper/     cross-cutting utilities (TxManager, pagination, validators, domain guards)
constant/   permission keys, audit keys, role names, env keys
middleware/ RBAC authorize
global/     response envelope + typed ErrorResponse
config/     DB connection, AutoMigrate, seeds
server/     routes.go — manual DI wiring + route registration
```
Dependency direction is one-way: `handler → usecase → repository → model`.

## 2. Request lifecycle
`RequestId → RequestLogger → PanicHandler → CORS → VerifyAuthToken → Authorize(permission) → Handler → Usecase → Repository → Mapper → global.CreateResponse`.
The response envelope is `{ code, status, data, message, error_code }`.

## 3. Dependency wiring (single place)
All dependencies are constructed manually in `server/routes.go`:
`repository.NewXxxRepository(db)` → `usecase.NewXxxUsecase(...)` → `handler.NewXxxHandler(...)` → route group.
There is no DI container, service locator, or global singleton (except the `validate` instance in `helper/validator.go`). Wire every new module here.

## 4. Transactions
- Cross-table mutations use `helper.TxManager` (`helper/tx_manager.go`).
- Pattern: `tx := u.txManager.New(); defer tx.CheckPanic();` then `tx.Rollback()` on every early return and `tx.Commit()` at the end.
- Repository mutators take `tx helper.Tx`; `dbFromTx(tx)` falls back to the base `*gorm.DB` when nil.
- Audit logs are written AFTER commit, outside the transaction (best-effort).

## 5. Auth & RBAC
- JWT: short-lived access token (claims `user_id`, `role_id`) + long-lived refresh token. Secrets/expiry come from env (`constant.*`).
- `VerifyAuthToken` (`server/middleware.go`) parses the bearer token and stores `user_id` / `role_id` in `c.Locals`.
- `middleware.Authorize(rbacRepo, permKey)` guards each route group; `Superadmin` bypasses all checks.
- Permission flow: add `constant.PermX` → add a `permissionSeed` entry in `config/seed_rbac.go` → attach the route to a guarded group in `server/routes.go`.
- Roles are seeded in `config/seed_roles.go`; the canonical list is `constant/role.go`.

## 6. Database
- MySQL via GORM; `SingularTable: true`; UUID v7 PKs (`model.BaseModel.BeforeCreate`); soft delete via `gorm.DeletedAt`.
- Schema is managed by `AutoMigrate` in `config/migration.go` — register every new model there.
- Schema/seed run only when `RUN_MIGRATION_AND_SEED=true` (`config/database.go`); seeds are idempotent (`SeedRoles`, `SeedRBAC`, `SeedBootstrapAdmin`).
- A `localhost` `GO_ENV` opens an SSH tunnel to the DB (`config/database.go`).

## 7. Configuration
- Env is loaded with `godotenv` (`config.Init`) and read via `config.GetEnv(key)`.
- Env keys are constants in `constant/constant.go`; document new vars in `.env.example`.

## 8. Core domain invariants
- Cylinder lifecycle is the central invariant:
  `EMPTY → READY_TO_FILL → FILLED → READY → IN_TRANSIT → OUTSTANDING → EMPTY` (loop), with `MAINTENANCE`, `LOST`, `WRITE_OFF` as side states.
- Transitions are guarded by pure functions in `helper/outbound.go` (`CanIssueOnDO`, `ValidateDOCylinders`, `ValidateExchangeOutCylinders`, `ValidateExchangeInCylinders`) and recorded in `CylinderLedger`.
- Sales flow: no PO/SO-to-barcode reservation yet. A DO is issued directly from barcodes; `sales_order_id` is optional and, when set, validates item + remaining qty and updates `qty_delivered` / SO status (`PARTIAL` / `COMPLETED`).
- Invoice is 1:1 with DO and is created in the same transaction as Issue DO.

## 9. Documentation contract
- Any change to sales, delivery, or cylinder-status logic MUST be reflected in `../progas-docs/03-business-flow.md` and logged in `../progas-docs/CHANGELOG.md` and `../progas-docs/00-progress.md`.
- Keep `AGENTS.md` accurate — it is the entry point for agents.

## 10. Extension checklist — adding a module
1. `model/xxx.go` (embed `BaseModel`) + register in `config/migration.go`.
2. `enum/xxx_status.go` if it has a status (with `IsValid`/`Scan`/`Value`).
3. `repository/xxx.go` (interface + GORM impl; mutators take `tx helper.Tx`).
4. `dto/xxx.go` (Request/Response + concrete `PaginatedXxxList` if listed).
5. `mapper/xxx.go` (model → dto).
6. `usecase/xxx.go` (business rules, transactions, audit after commit).
7. `handler/xxx.go` (validate, call usecase, swagger annotations).
8. `server/routes.go` — wire repo/usecase/handler + guarded route group.
9. `constant/permission.go` + `config/seed_rbac.go` for new permissions.
10. `constant/audit.go` for new audit actions/objects.
11. Tests (see `testing.md`) and docs updates (see §9).

## 11. Deployment
- Multi-stage `Dockerfile` (static binary on `alpine`); `docker-compose.yml` runs the app with `network_mode: host` (MySQL is a native VPS service).
- CI: `.github/workflows/deploy-dev.yml` deploys on push to `development` via SSH. Branches: `main` (default), `development` (active).
