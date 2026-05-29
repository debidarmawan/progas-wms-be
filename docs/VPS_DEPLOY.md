# Deploy di VPS (Docker + Nginx, port 3131)

Arsitektur:

```
Internet :3131
    │
    ▼
┌─────────────┐     ┌──────────────────┐
│   nginx     │────▶│  backend (Go)    │
│  :3131      │     │  :3131 (internal)│
└─────────────┘     └────────┬─────────┘
                             │
                             ▼
                      MySQL (VPS / remote)
```

Backend **tidak** dipublish ke host; hanya Nginx yang membuka port **3131**.

---

## 1. Prasyarat VPS

- Ubuntu 22.04+ / Debian 12+ (disarankan)
- Docker Engine + Docker Compose plugin
- MySQL 8 (di VPS yang sama, atau managed DB)
- Port **3131** dibuka di firewall

```bash
# Contoh Ubuntu
sudo apt update && sudo apt install -y docker.io docker-compose-v2 git
sudo usermod -aG docker $USER
# logout/login ulang
```

---

## 2. Clone & environment

```bash
cd /opt
sudo git clone https://github.com/progas/progas-wms-be.git
cd progas-wms-be
sudo chown -R $USER:$USER .

cp .env.example .env
nano .env
```

Isi minimal `.env` untuk production:

```env
GO_ENV=production
PORT=3131
DB_MAX_POOL=10

AUTH_TOKEN_EXPIRED_IN_MINUTES=15
REFRESH_TOKEN_EXPIRED_IN_DAYS=7
AUTH_TOKEN_SECRET_KEY=<random-min-32-chars>
REFRESH_TOKEN_SECRET_KEY=<random-min-32-chars>

# MySQL di VPS yang sama (service name host.docker.internal tidak ada di Linux — pakai IP host)
DB_URL=progas_user:STRONG_PASSWORD@tcp(172.17.0.1:3306)/progas_wms?charset=utf8mb4&parseTime=True&loc=Local
```

> **MySQL di host Linux:** dari container, akses host via gateway bridge `172.17.0.1` atau bind MySQL ke `0.0.0.0` dan gunakan IP publik/private VPS. Alternatif: tambahkan service `mysql` di `docker-compose.yml` (tidak termasuk default repo).

Pastikan MySQL mengizinkan koneksi dari Docker network dan database `progas_wms` sudah dibuat.

---

## 3. Jalankan

```bash
docker compose up -d --build
docker compose ps
docker compose logs -f
```

---

## 4. Verifikasi

```bash
curl http://127.0.0.1:3131/api/v1/health
# ok

curl -X POST http://127.0.0.1:3131/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"...","password":"..."}'
```

Dari luar VPS:

```text
http://IP_VPS:3131/api/v1/...
```

---

## 5. Firewall

```bash
# UFW
sudo ufw allow 3131/tcp
sudo ufw allow OpenSSH
sudo ufw enable
```

---

## 6. Update deploy

```bash
cd /opt/progas-wms-be
git pull
docker compose up -d --build
```

---

## 7. HTTPS (opsional)

Port 3131 + HTTPS tidak standar untuk Let's Encrypt (biasanya 80/443). Opsi:

1. **TLS di Nginx** — pasang sertifikat manual di `nginx.conf` (`listen 3131 ssl`) dan mount `fullchain.pem` / `privkey.pem`.
2. **Reverse proxy di host** — Nginx/Caddy di host pada `:443` → `http://127.0.0.1:3131` (dua layer nginx; cukup untuk production dengan domain).

Contoh Caddy di host (port 443 → 3131):

```caddy
api.example.com {
    reverse_proxy 127.0.0.1:3131
}
```

---

## 8. Troubleshooting

| Masalah | Cek |
|---------|-----|
| Connection refused | `docker compose ps`, port 3131 listen `ss -tlnp \| grep 3131` |
| DB error saat start | `docker compose logs backend`, DSN & MySQL bind address |
| 502 dari nginx | backend belum ready: `docker compose logs backend` |
| Swagger 404 | `GO_ENV` harus `development` |

```bash
docker compose logs backend --tail=100
docker compose logs nginx --tail=50
```

---

## File terkait

| File | Fungsi |
|------|--------|
| `docker-compose.yml` | backend + nginx |
| `nginx.conf` | Reverse proxy port 3131 |
| `Dockerfile` | Image backend |
| `.env` | Secret & DSN (jangan di-commit) |

Folder `deploy/` diabaikan — gunakan file di root repo seperti di atas.
