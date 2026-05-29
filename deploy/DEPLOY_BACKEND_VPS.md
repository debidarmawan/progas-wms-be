# Deploy Backend — VPS + Aiven MySQL (dev open access)

- **Repo:** https://github.com/debidarmawan/progas-wms-be.git  
- **API:** `http://103-169-206-116.domainesia.io/api/v1`  
- **DB:** Aiven MySQL — mode **dev / open access** (tanpa whitelist IP, tanpa `ca.pem`)

## Stack di VPS

| Service | Fungsi |
|---------|--------|
| **backend** | Go API |
| **nginx** | Port 80 → `/api/*` |

---

## 1. Aiven (dev open access)

Di Aiven Console untuk service MySQL dev Anda:

1. Pastikan **Public access** / open access aktif (tidak perlu tambah IP VPS)
2. Salin **Host**, **Port**, **User**, **Password**, **Database name**
3. **Tidak perlu** download CA untuk setup dev ini

---

## 2. Format `DB_URL`

Dari Aiven biasanya seperti:

`mysql://avnadmin:PASSWORD@mysql-xxxx.a.aivencloud.com:12345/defaultdb`

Ubah ke format Go:

```env
DB_URL=avnadmin:PASSWORD@tcp(mysql-xxxx.a.aivencloud.com:12345/defaultdb?charset=utf8mb4&parseTime=True&loc=UTC&tls=skip-verify
```

| Bagian | Nilai dev |
|--------|-----------|
| TLS | `tls=skip-verify` (tanpa ca.pem) |
| `DB_CA_CERT` | kosong |

> `tls=skip-verify` hanya untuk **development**. Production pakai `tls=aiven` + `ca.pem`.

---

## 3. Deploy di VPS

```bash
ssh root@103-169-206-116.domainesia.io

apt update && apt install -y ca-certificates curl git ufw
curl -fsSL https://get.docker.com | sh
apt install -y docker-compose-plugin
ufw allow OpenSSH && ufw allow 80/tcp && ufw enable

mkdir -p /opt/progas && cd /opt/progas
git clone https://github.com/debidarmawan/progas-wms-be.git
cd progas-wms-be/deploy

cp .env.example .env
nano .env   # isi DB_URL dari Aiven + JWT secret + bootstrap admin

chmod +x deploy.sh
./deploy.sh
```

Generate JWT secret:

```bash
openssl rand -base64 32
openssl rand -base64 32
```

---

## 4. Tes

```bash
curl http://103-169-206-116.domainesia.io/api/v1/health

curl -s -X POST http://103-169-206-116.domainesia.io/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@progas.local","password":"PASSWORD_BOOTSTRAP_ANDA"}'
```

Log:

```bash
docker compose logs -f backend
```

---

## Troubleshooting

| Masalah | Solusi |
|---------|--------|
| `connection refused` | Cek host/port Aiven, service MySQL running |
| `Access denied` | User/password/database salah |
| TLS error | Pastikan `tls=skip-verify` ada di `DB_URL` |
| Backend restart | `docker compose logs backend` |

---

## Production nanti

Ganti ke `tls=aiven`, whitelist IP VPS, pasang `ca.pem` — lihat `deploy/certs/README.md`.
