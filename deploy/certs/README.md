# Sertifikat Aiven (opsional)

**Dev open access:** tidak perlu folder ini. Pakai di `.env`:

```env
DB_URL=...@tcp(HOST:PORT)/DB?charset=utf8mb4&parseTime=True&loc=UTC&tls=skip-verify
DB_CA_CERT=
```

**Production:** download `ca.pem` dari Aiven, simpan di sini, lalu:

```env
DB_URL=...&tls=aiven
DB_CA_CERT=/app/certs/ca.pem
```

Tambahkan di `docker-compose.yml` pada service `backend`:

```yaml
volumes:
  - ./certs:/app/certs:ro
```
