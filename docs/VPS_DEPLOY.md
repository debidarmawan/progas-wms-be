# Deploy di VPS (Docker + MySQL + Nginx, port 3131)

Arsitektur:

```
Internet :3131                    PC Anda (DBeaver)
    │                                  │
    ▼                                  │ SSH tunnel atau :3306
┌─────────────┐     ┌──────────────┐   │ (jika MYSQL_BIND=0.0.0.0)
│   nginx     │────▶│   backend    │   │
│  :3131      │     │   :3131      │   │
└─────────────┘     └──────┬───────┘   │
                           │           │
                           ▼           ▼
                    ┌─────────────────────────┐
                    │  mysql (Docker) :3306   │
                    │  volume: mysql_data     │
                    └─────────────────────────┘
```

---

## 1. File `.env` di VPS

Salin dari `.env.example` dan isi (ganti semua password):

```env
GO_ENV=production
PORT=3131
DB_MAX_POOL=10

AUTH_TOKEN_EXPIRED_IN_MINUTES=15
REFRESH_TOKEN_EXPIRED_IN_DAYS=7
AUTH_TOKEN_SECRET_KEY=<random-32+>
REFRESH_TOKEN_SECRET_KEY=<random-32+>

MYSQL_ROOT_PASSWORD=<root-password-kuat>
MYSQL_DATABASE=progas_wms
MYSQL_USER=progas_app
MYSQL_PASSWORD=<password-app-kuat>
MYSQL_PORT=3306
MYSQL_BIND=127.0.0.1

# Penting: host "mysql" = nama service Docker, BUKAN localhost
DB_URL=progas_app:<password-app-kuat>@tcp(mysql:3306)/progas_wms?charset=utf8mb4&parseTime=True&loc=Local
```

| Variabel | Fungsi |
|----------|--------|
| `MYSQL_BIND=127.0.0.1` | MySQL hanya di VPS (disarankan) + SSH tunnel ke DBeaver |
| `MYSQL_BIND=0.0.0.0` | MySQL bisa diakses langsung dari internet (batasi IP di firewall) |
| `DB_URL` … `@tcp(mysql:3306)` | Backend connect ke container MySQL |

---

## 2. Jalankan stack

```bash
docker compose up -d --build
docker compose ps
```

Tunggu `progas-mysql` status **healthy**, lalu backend jalan (migrate + seed otomatis).

```bash
curl http://127.0.0.1:3131/api/v1/health
```

---

## 3. DBeaver dari PC Anda

### Opsi A — SSH tunnel (disarankan, `MYSQL_BIND=127.0.0.1`)

Di PC (terminal), biarkan jendela ini terbuka:

```bash
ssh -L 3306:127.0.0.1:3306 user@IP_VPS_ANDA
```

Di **DBeaver** → New Connection → MySQL:

| Field | Nilai |
|-------|--------|
| Host | `127.0.0.1` |
| Port | `3306` |
| Database | `progas_wms` |
| Username | `progas_app` |
| Password | sama dengan `MYSQL_PASSWORD` di `.env` |

Root (opsional): user `root`, password `MYSQL_ROOT_PASSWORD`.

### Opsi B — Langsung ke IP VPS (`MYSQL_BIND=0.0.0.0`)

1. Di `.env` VPS: `MYSQL_BIND=0.0.0.0`
2. `docker compose up -d`
3. Firewall **hanya IP Anda**:

```bash
sudo ufw allow from IP_PC_ANDA to any port 3306
sudo ufw deny 3306
```

DBeaver:

| Field | Nilai |
|-------|--------|
| Host | `IP_VPS` |
| Port | `3306` |
| Database | `progas_wms` |
| Username | `progas_app` |

> Jangan buka port 3306 ke `0.0.0.0/0` tanpa batasan — risiko brute-force.

---

## 4. Firewall VPS

```bash
sudo ufw allow 3131/tcp    # API
sudo ufw allow OpenSSH
# Port 3306: hanya jika pakai Opsi B + allow from IP tertentu
sudo ufw enable
```

---

## 5. Backup data MySQL

Data ada di volume Docker `mysql_data`:

```bash
docker compose exec mysql mysqldump -u root -p"${MYSQL_ROOT_PASSWORD}" progas_wms > backup.sql
```

Restore:

```bash
docker compose exec -T mysql mysql -u root -p"${MYSQL_ROOT_PASSWORD}" progas_wms < backup.sql
```

---

## 6. Update & troubleshooting

```bash
git pull
docker compose up -d --build
```

| Masalah | Solusi |
|---------|--------|
| Backend restart loop | `docker compose logs backend` — cek `DB_URL` harus `mysql:3306` |
| MySQL unhealthy | `docker compose logs mysql` — cek `MYSQL_ROOT_PASSWORD` |
| DBeaver connection refused (Opsi A) | Pastikan SSH tunnel aktif |
| DBeaver timeout (Opsi B) | `MYSQL_BIND`, firewall, `ss -tlnp \| grep 3306` |

```bash
docker compose logs mysql --tail=50
docker compose logs backend --tail=50
```

---

## 7. Masih pakai Aiven di laptop lokal

Di mesin dev (bukan VPS), Anda bisa tetap pakai Aiven di `.env` dan jalankan **tanpa** MySQL container:

```bash
# .env lokal — DB_URL ke Aiven
docker compose up -d --build backend nginx
```

Service `mysql` tidak ikut jalan; pastikan `DB_URL` tidak memakai host `mysql`.

---

## File terkait

| File | Fungsi |
|------|--------|
| `docker-compose.yml` | mysql + backend + nginx |
| `nginx.conf` | Reverse proxy :3131 |
| `.env.example` | Template env |
