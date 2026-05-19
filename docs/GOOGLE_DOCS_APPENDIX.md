# Appendix — Salin ke Google Docs PRD

> Paste section di bawah ini ke dokumen [PRD Utama WMS Gas Industri](https://docs.google.com/document/d/1J1qM9D6UIt6z3YojPyuSHvUUYelAtX2vLUisuW44NdA/edit?tab=t.0) sebagai **Section 7 — Implementation Task Breakdown**.

---

## 7. IMPLEMENTATION TASK BREAKDOWN (BACKEND)

**Repo:** progas-wms-be | **Dokumen lengkap:** `docs/PRD_TASK_BREAKDOWN.md`

### 7.1 Role Sistem

| Role (Database) | Pemetaan PRD |
|-----------------|--------------|
| Superadmin | Super Admin — akses penuh |
| Warehouse Admin | Admin Gudang & QC |
| Logistic Admin | Admin Logistik & Distribusi |
| Manager | Manajer — read-only |

### 7.2 Fase Implementasi

| Fase | Scope | Status |
|------|-------|--------|
| 0 | Auth, User, Role, Docker | Selesai |
| 1 | RBAC middleware, Audit Log, seed permission | Selesai |
| 2 | Master Data: Item, Cylinder, Customer | Planned |
| 3 | Inbound & Filling Batch | Planned |
| 4 | DO, Cylinder Swap, Fleet | Planned |
| 5 | Work Order, Dashboard, Laporan | Planned |

### 7.3 Fase 1 — Permission API (Aktif)

| Permission | Endpoint | Superadmin | Warehouse | Logistic | Manager |
|------------|----------|:----------:|:---------:|:--------:|:-------:|
| auth.logout | POST /logout | Ya | Ya | Ya | Ya |
| role.read | GET /roles | Ya | Ya | Ya | Ya |
| user.create | POST /users | Ya | Tidak | Tidak | Tidak |

### 7.4 Audit Log (Aktif)

Setiap aksi krusial dicatat: User ID, Action, Object Type, Object ID, Details, Waktu.

| Action | Trigger |
|--------|---------|
| USER_CREATE | User baru dibuat |
| USER_LOGIN | Login berhasil |

### 7.5 Backlog Epik PRD → Fase

- **Epik 1 (Master Data):** Registrasi tabung barcode, item spare part, pelanggan & kuota → Fase 2
- **Epik 2 (Produksi):** Filling batch, cross-gas validation → Fase 3
- **Epik 3 (Logistik):** DO, cylinder swap, outstanding formula → Fase 4
- **Dashboard & Laporan:** Fase 5

---

*Detail task ID per endpoint: lihat file `docs/PRD_TASK_BREAKDOWN.md` di repository backend.*
