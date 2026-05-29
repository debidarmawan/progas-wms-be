# Deploy Backend Progas WMS — VPS Ubuntu 24.04 (tanpa domain)

Hostname Domainesia Anda:

`103-169-206-116.domainesia.io`

API setelah deploy:

`http://103-169-206-116.domainesia.io/api/v1`

---

## Ringkasan stack

| Service | Fungsi |
|---------|--------|
| **mysql** | Database |
| **backend** | Go API (port internal 3131) |
| **nginx** | Port **80** publik → proxy ke backend |

---

## Langkah 1 — SSH ke VPS

```bash
ssh root@103-169-206-116.domainesia.io
# atau: ssh root@103.169.206.116
```

---

## Langkah 2 — Install Docker (Ubuntu 24.04)

```bash
apt update && apt upgrade -y
apt install -y ca-certificates curl git ufw

curl -fsSL https://get.docker.com | sh
systemctl enable docker
systemctl start docker
apt install -y docker-compose-plugin

docker --version
docker compose version
```

Firewall:

```bash
ufw allow OpenSSH
ufw allow 80/tcp
ufw allow 443/tcp
ufw enable
```

---

## Langkah 3 — Clone backend

```bash
mkdir -p /opt/progas && cd /opt/progas

git clone https://github.com/debidarmawan/progas-wms-be.git
cd progas-wms-be/deploy
```

---

## Langkah 4 — File `.env`

```bash
cp .env.example .env
nano .env
```

Isi password & secret (wajib):

```bash
# Generate 2 secret JWT
openssl rand -base64 32
openssl rand -base64 32
```

Contoh isian penting:

```env
PUBLIC_HOST=103-169-206-116.domainesia.io
API_PUBLIC_URL=http://103-169-206-116.domainesia.io/api/v1

MYSQL_ROOT_PASSWORD=...
MYSQL_PASSWORD=...
DB_URL=progas:SAMA_DENGAN_MYSQL_PASSWORD@tcp(mysql:3306)/progas_wms?charset=utf8mb4&parseTime=True&loc=Local

AUTH_TOKEN_SECRET_KEY=...   # dari openssl
REFRESH_TOKEN_SECRET_KEY=...

BOOTSTRAP_ADMIN_EMAIL=admin@progas.local
BOOTSTRAP_ADMIN_NAME=Superadmin
BOOTSTRAP_ADMIN_PASSWORD=PasswordKuat123!
```

Simpan (`Ctrl+O`, `Enter`, `Ctrl+X`).

> Setelah login pertama berhasil, **hapus** baris `BOOTSTRAP_ADMIN_*` dari `.env` lalu `docker compose up -d --build backend`.

---

## Langkah 5 — Deploy

```bash
chmod +x deploy.sh
./deploy.sh
```

Atau manual:

```bash
docker compose up -d --build
docker compose ps
```

---

## Langkah 6 — Tes

Dari VPS:

```bash
curl http://127.0.0.1/api/v1/health
# ok
```

Dari laptop/browser:

```
http://103-169-206-116.domainesia.io/api/v1/health
```

Login (setelah ada user di DB / seed):

```bash
curl -X POST http://103-169-206-116.domainesia.io/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"..."}'
```

---

## User admin pertama

Saat startup, backend otomatis:

1. Migrasi tabel
2. Seed role (Superadmin, Warehouse Admin, …)
3. Seed RBAC
4. Buat user admin jika `BOOTSTRAP_ADMIN_EMAIL` + `BOOTSTRAP_ADMIN_PASSWORD` diisi dan belum ada user

Cek log:

```bash
docker compose logs backend | tail -20
# Harus ada: Role seed completed, RBAC seed completed, Bootstrap admin created: ...
```

Tes login:

```bash
curl -s -X POST http://103-169-206-116.domainesia.io/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@progas.local","password":"PasswordKuat123!"}'
```

---

## Perintah berguna

```bash
cd /opt/progas/progas-wms-be/deploy

docker compose logs -f backend
docker compose logs -f mysql
docker compose restart backend
docker compose up -d --build   # setelah git pull

# Backup DB
docker exec progas-mysql mysqldump -u root -p"ROOT_PASSWORD" progas_wms > backup.sql
```

---

## Nanti saat pasang frontend

Di `.env` frontend (laptop/VPS), set:

```env
NEXT_PUBLIC_API_BASE_URL=http://103-169-206-116.domainesia.io/api/v1
```

Kalau nanti pakai domain + HTTPS, ganti ke `https://domain-anda.com/api/v1` lalu rebuild frontend.

---

## Troubleshooting

| Masalah | Solusi |
|---------|--------|
| Tidak bisa akses dari browser | Cek firewall Domainesia panel + `ufw status` |
| Backend restart loop | `docker compose logs backend` — biasanya `DB_URL` salah |
| MySQL tidak healthy | `docker compose logs mysql` — tunggu ~1 menit pertama kali |
| Build gagal OOM | Pastikan swap aktif: `free -h` |

---

## Update kode

```bash
cd /opt/progas/progas-wms-be
git pull
cd deploy
docker compose up -d --build
```
