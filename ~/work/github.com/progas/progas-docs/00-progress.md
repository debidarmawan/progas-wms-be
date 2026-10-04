# Progress Log — Progas WMS

> Sumber kebenaran bersama: status implementasi lintas modul, prioritas pengerjaan, dan riwayat progres.
> Dibuat untuk memastikan progres dapat dilacak dan dilanjutkan tanpa analisis ulang.

**Workspace:**
- `progas-wms-be` — Backend (Go 1.26, Fiber v3, GORM, MySQL)
- `progas-wms` — Web frontend (Next.js)
- `progas-wms-mobile` — Mobile app (React Native, Expo SDK 57, NativeWind v4)
- `progas-docs` — Dokumentasi (architecture, modules-api, business-flow, roadmap, changelog)

---

## Status Implementasi Saat Ini (As-Is)

| Modul | Status | Keterangan |
|-------|--------|------------|
| Backend: Auth, User, Role, RBAC, Audit Log | ✅ Selesai | Login JWT, RBAC middleware, seed role, audit trail |
| Backend: Master Item, Cylinder, Customer, Customer Item Price | ✅ Selesai | CRUD + validasi barcode duplikat, ownership, hydrotest |
| Backend: Customer PO (DRAFT/CONFIRMED) + Sales Order (DRAFT/CONFIRMED/PARTIAL/COMPLETED) | ✅ Selesai | Level master item + qty, `sales_order_id` opsional pada DO |
| Backend: DO Issue dengan support SO-linked | ✅ Selesai | Validasi item + sisa qty SO, update `qty_delivered` & status SO |
| Backend: Invoice & Payment (1:1 per DO) | ✅ Selesai | Otomatis saat Issue DO, tracking piutang pelanggan |
| Backend: Inbound, Filling Batch, QC | ✅ Selesai | EMPTY → READY → IN_TRANSIT lifecycle |
| Backend: Exchange, Fleet, Driver | ✅ Selesai | Outbound gate swap + outstanding |
| Backend: Maintenance (WO, hydrotest), Inventory, Dashboard, Reports | ✅ Selesai | Virtual warehouse, stock opname, ledger, turnaround |
| Frontend: Halaman PO Pelanggan, Sales Order, Surat Jalan (list) | ✅ Selesai | |
| Frontend: "Buat Surat Jalan" | ⚠️ Parsial | Sudah bisa pilih SO opsional + scan barcode, tapi **tidak memisahkan per baris SO qty** dan tidak ada validasi sisa qty per line di form |
| Frontend: Trip/Pengiriman Armada | ❌ Belum | Entitas & halaman belum ada |
| Mobile: Design system + navigation shell | ✅ Selesai (shell) | Expo SDK 57, NativeWind v4, react-native-reusables, clean-architecture scaffolding — **masih kosong fitur bisnis** |
| Mobile: Role `driver` di backend | ❌ Belum | Role driver belum ada (RBAC role.go hanya punya 4 role) |
| Mobile: API scan-bongkar | ❌ Belum | API `POST /mobile/delivery-orders/:doId/scan-deliver` belum ada |

---

## Prioritas Pengerjaan (Urutan)

### P1 — Frontend: Perbaiki "Buat Surat Jalan" untuk mode SO-linked (TIDAK bergantung mobile)
- Tampilkan baris SO yang dipilih + qty per baris + sisa qty
- User bisa memilih baris item yang mau dikirim & qty per baris
- Validasi: barang tidak melebihi sisa qty SO line
- `qty_delivered` SO line bertambah, SO jadi PARTIAL/COMPLETED
- **File terdampak (frontend):** `app/dashboard/outbound/delivery-orders/new/page.tsx`, `lib/api/outbound.ts`
- **Dokumentasi:** update `03-business-flow.md` §3, entry `CHANGELOG.md`

### P2 — Backend: Trip/Pengiriman Armada (Fase 2 roadmap)
- Model `DeliveryTrip` + `TripCylinder`, repository, usecase, handler, route (`/logistics/trips`)
- Status tabung baru: `RESERVED`/`LOADED`/`DELIVERED` di `enum.CylinderStatus`
- API: `POST /logistics/trips` (buat trip, assign fleet+driver+daftar DO), `POST /logistics/trips/:id/load` (scan barcode → reserve), `POST /logistics/trips/:id/depart` (berangkat → IN_TRANSIT)
- `Issue DO` menerima `trip_id` opsional → jika ada, tabung jadi `RESERVED` (bukan `IN_TRANSIT`)
- **File terdampak (backend):** `model/delivery_trip.go`, `repository/delivery_trip.go`, `usecase/delivery_trip.go`, `handler/delivery_trip.go`, `server/routes.go`, `helper/validate.go`, `constant/permission.go`, `config/migration.go`, `helper/pagination.go`, `mapper/delivery_trip.go`, `dto/delivery_trip.go`
- **Dokumentasi:** update `03-business-flow.md` §3 & §7, `04-roadmap.md` → pindah ke as-is, entry `CHANGELOG.md`

### P3 — Backend: Role `driver` + permission `mobile.*` (pre-requisite mobile)
- Tambah role `driver` di `constant/role.go` + seed RBAC
- Update `01-architecture.md` jika perlu
- **Dokumentasi:** entry `CHANGELOG.md`

### P4 — Backend: API Mobile scan-bongkar (Fase 3 roadmap)
- `POST /mobile/delivery-orders/:doId/scan-deliver` (body: `{ barcodes: string[] }`)
- Validasi barcode terdaftar di `TripCylinder` trip aktif, status `RESERVED`/`LOADED`, item sesuai DO, qty tidak melebihi sisa
- Update `TripCylinder.status → DELIVERED`, update `DeliveryOrderDetail.qty` terlink, update Trip status → COMPLETED
- **File terdampak (backend):** `handler/delivery_order.go` (method baru), `usecase/delivery_order.go`, `repository/delivery_order.go`, `server/routes.go`, `constant/permission.go`
- **Dokumentasi:** update `03-business-flow.md` §7, entry `CHANGELOG.md`

### P5 — Mobile: Implementasi fitur scan-bongkar (setelah P3 & P4 selesai)
- Login driver, list DO untuk driver, scan barcode, submit delivery (realtime + offline queue)
- **File terdampak (mobile):** `src/screens/*`, `src/feature/delivery/*`, `src/packages/mobile/...`
- **Dokumentasi:** entry `CHANGELOG.md`

---

## Cara Menggunakan Log Ini

1. Setiap kali pengerjaan selesai, tandai item **✅ Selesai** dan update status "Status Implementasi Saat Ini".
2. Setiap perubahan kode/flow **WAJIB** dicatat di `CHANGELOG.md` (sesuai peraturan).
3. Referensi file terdampak dicantumkan agar next agent bisa langsung melanjutkan tanpa eksplorasi ulang.
