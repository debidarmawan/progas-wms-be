# Deploy di VPS — Backend saja (Frontend di Vercel)

VPS Ubuntu 24.04 (~2 GB RAM) menjalankan **MySQL native + API (Docker) + Nginx (Docker)**. Frontend Next.js/React di-host di **Vercel** dan memanggil API lewat `https://...vercel.app` → `http://IP_VPS:3131/api/v1`.

## Arsitektur

```
┌─────────────────────┐         HTTP :3131          ┌──────────────────────────────┐
│  Vercel (Frontend)  │ ──────────────────────────▶ │  VPS                         │
│  NEXT_PUBLIC_API_   │      CORS_ORIGINS           │  nginx (Docker) :3131        │
│  URL → VPS:3131     │                             │         │                    │
└─────────────────────┘                             │         ▼                    │
                                                    │  backend (Docker) :3131      │
         DBeaver (opsional) ── SSH tunnel :3306 ──▶ │         │                    │
                                                    │         ▼                    │
                                                    │  MySQL native (127.0.0.1)    │
                                                    └──────────────────────────────┘
```

| Komponen | Lokasi | Port publik |
|----------|--------|-------------|
| Frontend | Vercel | 443 (HTTPS) |
| API | VPS (Docker) | **3131** |
| MySQL | VPS (native) | **127.0.0.1:3306** (disarankan) |

Detail konfigurasi frontend: [`VERCEL_FRONTEND.md`](./VERCEL_FRONTEND.md).

---

## 1. Spesifikasi VPS

Disarankan: **Ubuntu 24.04**, min. **2 GB RAM**, **2 vCPU**, **30 GB** disk.

| Service | `mem_limit` | Catatan |
|---------|-------------|---------|
| backend | 256m | `DB_MAX_POOL=5` |
| nginx | 64m | reverse proxy saja |
| MySQL | — | Install native (`apt`), tune sesuai RAM VPS |

Opsional: tambah swap 1–2 GB jika OOM saat build pertama.

---

## 2. MySQL native di VPS

Install dan buat database/user (sesuaikan password):

```bash
sudo apt update && sudo apt install -y mysql-server
sudo mysql -e "
CREATE DATABASE progas_wms CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'progas_app'@'localhost' IDENTIFIED BY '<password-app-kuat>';
CREATE USER 'progas_app'@'172.17.%' IDENTIFIED BY '<password-app-kuat>';
GRANT ALL PRIVILEGES ON progas_wms.* TO 'progas_app'@'localhost';
GRANT ALL PRIVILEGES ON progas_wms.* TO 'progas_app'@'172.17.%';
FLUSH PRIVILEGES;
"
```

User `172.17.%` diperlukan agar container backend (Docker bridge) bisa connect ke MySQL host.

Pastikan MySQL bisa menerima koneksi dari Docker bridge. Di `/etc/mysql/mysql.conf.d/mysqld.cnf`:

```ini
bind-address = 0.0.0.0
```

Lalu restart: `sudo systemctl restart mysql`

**Jangan** buka port 3306 ke internet — lihat firewall di bawah.

---

## 3. File `.env` di VPS

Salin dari `.env.example` dan isi (ganti semua password):

```env
GO_ENV=production
PORT=3131
DB_MAX_POOL=5

CORS_ORIGINS=https://progas-wms.vercel.app
CORS_ALLOW_VERCEL_PREVIEW=true

AUTH_TOKEN_EXPIRED_IN_MINUTES=15
REFRESH_TOKEN_EXPIRED_IN_DAYS=7
AUTH_TOKEN_SECRET_KEY=<random-32+>
REFRESH_TOKEN_SECRET_KEY=<random-32+>

MYSQL_DATABASE=progas_wms
MYSQL_USER=progas_app
MYSQL_PASSWORD=<password-app-kuat>
MYSQL_PORT=3306

# host.docker.internal = MySQL native di VPS (bukan localhost / mysql)
DB_URL=progas_app:<password-app-kuat>@tcp(host.docker.internal:3306)/progas_wms?charset=utf8mb4&parseTime=True&loc=Local
```

| Variabel | Fungsi |
|----------|--------|
| `CORS_ORIGINS` | Origin frontend yang boleh memanggil API (HTTPS Vercel) |
| `CORS_ALLOW_VERCEL_PREVIEW` | `true` = izinkan preview `*.vercel.app` |
| `DB_URL` … `@tcp(host.docker.internal:3306)` | Backend container connect ke MySQL native di host |

---

## 4. Jalankan stack

```bash
docker compose up -d --build
docker compose ps
```

Migrate + seed otomatis saat backend start.

```bash
curl http://127.0.0.1:3131/api/v1/health
```

Dari browser (setelah firewall): `http://IP_VPS:3131/api/v1/health`

---

## 5. Firewall VPS

```bash
sudo ufw allow 3131/tcp    # API untuk Vercel & klien
sudo ufw allow OpenSSH
sudo ufw enable
sudo ufw deny 3306         # MySQL tidak publik
```

---

## 6. HTTPS di VPS (opsional, disarankan)

Browser di Vercel (HTTPS) **bisa** memanggil API HTTP (`http://IP:3131`) — tidak mixed-content untuk XHR/fetch. Namun untuk keamanan dan cookie `Secure`, pertimbangkan:

- Subdomain API + Let's Encrypt (Caddy/nginx di depan), atau
- Cloudflare Tunnel ke backend

Tanpa HTTPS, pastikan frontend hanya menyimpan JWT di memory / httpOnly cookie dengan strategi yang aman.

---

## 7. DBeaver dari PC Anda

### Opsi A — SSH tunnel (disarankan)

```bash
ssh -L 3306:127.0.0.1:3306 user@IP_VPS_ANDA
```

DBeaver → MySQL: Host `127.0.0.1`, Port `3306`, DB `progas_wms`, user `progas_app`.

---

## 8. Backup data MySQL

```bash
mysqldump -h 127.0.0.1 -u progas_app -p progas_wms > backup.sql
```

Restore:

```bash
mysql -h 127.0.0.1 -u progas_app -p progas_wms < backup.sql
```

Atau via Makefile: `make mysql-cli` (butuh client `mysql` terinstall di VPS).

---

## 9. Update & troubleshooting

```bash
git pull
docker compose up -d --build
```

| Masalah | Solusi |
|---------|--------|
| CORS error di browser | Set `CORS_ORIGINS` = URL exact Vercel (dengan `https://`) |
| Backend restart loop | `docker compose logs backend` — cek `DB_URL` harus `host.docker.internal:3306` |
| Access denied for user | Pastikan user MySQL punya grant untuk `172.17.%` (Docker bridge) |
| Connection refused ke MySQL | Cek `bind-address` MySQL & `sudo systemctl status mysql` |
| Frontend tidak connect | Cek `NEXT_PUBLIC_API_URL` di Vercel & firewall 3131 |

```bash
docker compose logs backend --tail=50
sudo tail -50 /var/log/mysql/error.log
```

---

## 10. Dev lokal (DB cloud / Aiven)

Untuk development laptop dengan DB eksternal, cukup set `DB_URL` ke host cloud di `.env`:

```bash
docker compose up -d --build
```

---

## File terkait

| File | Fungsi |
|------|--------|
| `docker-compose.yml` | backend + nginx (MySQL di host) |
| `server/cors.go` | CORS dari env |
| `nginx.conf` | Reverse proxy :3131 |
| `.env.example` | Template env |
| `docs/VERCEL_FRONTEND.md` | Env & setup frontend Vercel |
