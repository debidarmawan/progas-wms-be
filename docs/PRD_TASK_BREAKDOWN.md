# PRD Implementation Task Breakdown
## Warehouse Management System (WMS) — Gas Industri

> **Dokumen sumber PRD:** [PRD Utama - WMS Gas Industri (Google Docs)](https://docs.google.com/document/d/1J1qM9D6UIt6z3YojPyuSHvUUYelAtX2vLUisuW44NdA/edit?tab=t.0)  
> **Backend repo:** `progas-wms-be`  
> **Terakhir diperbarui:** 2026-05-19

---

## Role yang Terdaftar di Database

| Role (DB) | Pemetaan PRD | Deskripsi Singkat |
|-----------|--------------|-------------------|
| **Superadmin** | Super Admin | Akses penuh, manajemen user, audit |
| **Warehouse Admin** | Admin Gudang & QC | Inbound, filling, inventory gudang, work order |
| **Logistic Admin** | Admin Logistik & Distribusi | DO, armada, cylinder exchange |
| **Manager** | Manajer / Direksi | Read-only dashboard & laporan |

---

## Ringkasan Fase

| Fase | Nama | Status | Target |
|------|------|--------|--------|
| **0** | Fondasi Auth & User | ✅ Selesai | Login JWT, user, role read, Docker |
| **1** | RBAC + Audit Log | ✅ Selesai | Permission middleware, seed, audit trail |
| **2** | Master Data | ⬜ Planned | Item, Cylinder, Customer, Warehouse |
| **3** | Inbound & Produksi | ⬜ Planned | Empty receiving, filling batch, QC |
| **4** | Outbound & Logistik | ⬜ Planned | DO, swap, outstanding, fleet |
| **5** | Maintenance & Laporan | ⬜ Planned | Work order, spare part, ledger, dashboard |

---

## Fase 0 — Fondasi (Selesai)

| ID | Task | Endpoint / Artefak | Status |
|----|------|---------------------|--------|
| 0.1 | Setup project Go + Fiber + GORM | `main.go`, `go.mod` | ✅ |
| 0.2 | Koneksi MySQL + AutoMigrate | `config/database.go` | ✅ |
| 0.3 | Login + JWT access/refresh | `POST /api/v1/login` | ✅ |
| 0.4 | Refresh token | `POST /api/v1/refresh-token` | ✅ |
| 0.5 | Logout (stateless) | `POST /api/v1/logout` | ✅ |
| 0.6 | Create user + bcrypt | `POST /api/v1/users` | ✅ |
| 0.7 | List & detail role | `GET /api/v1/roles`, `GET /api/v1/roles/:id` | ✅ |
| 0.8 | JWT middleware protected routes | `VerifyAuthToken` | ✅ |
| 0.9 | Docker + Compose | `Dockerfile`, `docker-compose.yml` | ✅ |
| 0.10 | Swagger (development) | `/swagger/*` | ✅ |

---

## Fase 1 — RBAC + Audit Log (Selesai)

| ID | Task | Permission Key | Role yang Diizinkan | Status |
|----|------|----------------|---------------------|--------|
| 1.1 | Model `AuditLog` + migrate | — | — | ✅ |
| 1.2 | Repository audit log | — | — | ✅ |
| 1.3 | Repository RBAC (`HasPermission`, `IsSuperAdmin`) | — | — | ✅ |
| 1.4 | Seed `role_key` + `role_key_mapping` (idempotent) | — | — | ✅ |
| 1.5 | Middleware `Authorize(permissionKey)` | — | — | ✅ |
| 1.6 | Pasang RBAC ke route existing | lihat tabel bawah | — | ✅ |
| 1.7 | Audit: create user | `USER_CREATE` | — | ✅ |
| 1.8 | Audit: login | `USER_LOGIN` | — | ✅ |
| 1.9 | Forbidden response konsisten | — | — | ✅ |

### Permission — API Saat Ini (Fase 1)

| Permission Key | Method | Path | Superadmin | Warehouse Admin | Logistic Admin | Manager |
|----------------|--------|------|:----------:|:---------------:|:--------------:|:-------:|
| `auth.logout` | POST | `/api/v1/logout` | ✅ | ✅ | ✅ | ✅ |
| `role.read` | GET | `/api/v1/roles/*` | ✅ | ✅ | ✅ | ✅ |
| `user.create` | POST | `/api/v1/users` | ✅ | ❌ | ❌ | ❌ |

> **Superadmin** bypass penuh di middleware (selaras PRD: akses tanpa batas).

---

## Fase 2 — Master Data (Epik 1 PRD)

| ID | Task | Permission Key (rencana) | User Story PRD |
|----|------|--------------------------|----------------|
| 2.1 | Model `Item` (gas, spare part, `IsSerialized`, berat) | `item.*` | 1.2 |
| 2.2 | Model `Cylinder` (barcode, ownership, status, hydrotest) | `cylinder.*` | 1.1 |
| 2.3 | Model `Customer` (kuota tabung, outstanding) | `customer.*` | §4.1 |
| 2.4 | Model `Warehouse` / lokasi rak | `warehouse.*` | §3 |
| 2.5 | API registrasi tabung (validasi SN unik) | `cylinder.create` | AC 1.1 |
| 2.6 | API master item non-serialized + min stock | `item.create` | AC 1.2 |
| 2.7 | API CRUD pelanggan + kuota | `customer.manage` | §4.1 |
| 2.8 | State machine validasi status tabung | — | §2.1 |
| 2.9 | Seed permission Fase 2 ke RBAC | — | — |

**Acceptance Criteria (dari PRD):**
- Tolak barcode duplikat
- Dropdown ownership: COMPANY, CUSTOMER, VENDOR
- Validasi `LastHydrotestDate`
- Flag `IsSerialized=false` abaikan barcode per unit
- Min stock alert

---

## Fase 3 — Inbound & Produksi (Epik 2 PRD)

| ID | Task | Permission Key | User Story |
|----|------|----------------|------------|
| 3.1 | Empty receiving | `inbound.empty_receive` | §3 Inbound |
| 3.2 | Filling batch log | `production.filling_batch` | 2.1 |
| 3.3 | Cross-gas validation | — | AC 2.1 |
| 3.4 | Batch rollback jika status tabung invalid | — | AC 2.1 |
| 3.5 | QC pre/post filling | `production.qc` | §3 |
| 3.6 | Auto status → READY setelah batch sukses | — | AC 2.1 |

**Role akses (PRD §5):** Warehouse Admin ✅ | Logistic Admin ❌ | Manager read-only

---

## Fase 4 — Outbound & Logistik (Epik 3 PRD)

| ID | Task | Permission Key | User Story |
|----|------|----------------|------------|
| 4.1 | Model Delivery Order (DO) | `do.*` | 3.1 |
| 4.2 | Manifest barcode + hitung berat | `do.create` | AC 3.1 |
| 4.3 | Overload protection (`MaxWeightKg`) | — | §4.4, AC 3.1 |
| 4.4 | Status tabung → IN_TRANSIT | — | AC 3.1 |
| 4.5 | Cylinder swap / gate in-out | `exchange.process` | 3.2 |
| 4.6 | Rumus outstanding: `Lama + OUT - IN` | — | §4.1 |
| 4.7 | Abaikan CUSTOMER ownership dari outstanding | — | AC 3.2 |
| 4.8 | Blokir / approval jika melebihi kuota | `exchange.approve` | AC 3.2 |
| 4.9 | Alert tabung tertukar antar pelanggan | — | §4.1 |
| 4.10 | Fleet management | `fleet.*` | §3 modul 5 |

**Role akses:** Logistic Admin ✅ | Warehouse Admin ❌ (DO) | Manager read-only

---

## Fase 5 — Maintenance, Dashboard & Laporan

| ID | Task | Permission Key | Modul PRD |
|----|------|----------------|-----------|
| 5.1 | Work order + kurangi spare part | `workorder.*` | §6 |
| 5.2 | Stok opname spare part | `inventory.stockopname` | §6 |
| 5.3 | Jadwal hydrotest / servis | `cylinder.hydrotest` | §3 modul 4 |
| 5.4 | Dashboard API (stok, outstanding, alert) | `dashboard.read` | §3 modul 1 |
| 5.5 | Stock ledger per barcode | `report.ledger` | §3 modul 8 |
| 5.6 | Turn-around rate report | `report.turnaround` | §3 modul 8 |
| 5.7 | Virtual warehouse (outstanding customer) | `inventory.virtual` | §3 modul 4 |

---

## Matriks RBAC Lengkap (Target Akhir — Semua Fase)

| Permission | Superadmin | Warehouse Admin | Logistic Admin | Manager |
|------------|:----------:|:---------------:|:--------------:|:-------:|
| `auth.logout` | ✅ | ✅ | ✅ | ✅ |
| `role.read` | ✅ | ✅ | ✅ | ✅ |
| `user.create` | ✅ | ❌ | ❌ | ❌ |
| `user.manage` | ✅ | ❌ | ❌ | ❌ |
| `item.*` | ✅ | ✅ read/write | ❌ | ✅ read |
| `cylinder.*` | ✅ | ✅ | ❌ read | ✅ read |
| `customer.*` | ✅ | ✅ read | ✅ read | ✅ read |
| `inbound.*` | ✅ | ✅ | ❌ | ✅ read |
| `production.*` | ✅ | ✅ | ❌ | ✅ read |
| `inventory.warehouse` | ✅ | ✅ | ❌ | ✅ read |
| `do.*` | ✅ | ❌ | ✅ | ✅ read |
| `exchange.*` | ✅ | ❌ | ✅ | ✅ read |
| `fleet.*` | ✅ | ❌ | ✅ | ✅ read |
| `workorder.*` | ✅ | ✅ | ❌ | ✅ read |
| `dashboard.read` | ✅ | ✅ | ✅ | ✅ |
| `report.*` | ✅ | ❌ | ❌ | ✅ |
| `audit.read` | ✅ | ❌ | ❌ | ✅ |

---

## Audit Log — Aksi yang Wajib Dicatat (PRD §5)

| Action Code | Trigger | Fase |
|-------------|---------|------|
| `USER_CREATE` | Create user berhasil | 1 |
| `USER_LOGIN` | Login berhasil | 1 (opsional) |
| `CYLINDER_CREATE` | Registrasi tabung baru | 2 |
| `CYLINDER_OWNERSHIP_CHANGE` | Ubah kepemilikan tabung | 2 |
| `CYLINDER_DELETE` | Hapus nomor seri | 2 |
| `TRANSACTION_UPDATE` | Ubah data transaksi | 3–4 |
| `DO_ISSUE` | Terbitkan surat jalan | 4 |
| `EXCHANGE_COMPLETE` | Selesai cylinder swap | 4 |

**Field audit:** `user_id`, `action`, `object_type`, `object_id`, `details` (JSON), `created_at`

---

## Catatan untuk Google Docs

Salin section **Fase 0–5** dan **Matriks RBAC** ke dokumen PRD utama sebagai **Section 7 — Implementation Task Breakdown & Backend Roadmap**, agar product & engineering satu sumber kebenaran.

---

## Referensi Teknis Backend

- JWT claims: `user_id`, `role_id` → `dto.JWTClaims`
- RBAC tables: `role`, `role_key`, `role_key_mapping`
- Middleware chain: `VerifyAuthToken` → `Authorize(permissionKey)`
- Seed RBAC: `config/seed_rbac.go` (idempotent, jalankan saat startup migrate)
