#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")"

if [[ ! -f .env ]]; then
  echo "Buat .env dulu: cp .env.example .env && nano .env"
  exit 1
fi

if grep -q '^AUTH_TOKEN_SECRET_KEY=$' .env || grep -q '^REFRESH_TOKEN_SECRET_KEY=$' .env; then
  echo "Isi AUTH_TOKEN_SECRET_KEY dan REFRESH_TOKEN_SECRET_KEY di .env"
  echo "Generate: openssl rand -base64 32"
  exit 1
fi

if ! grep -q '^DB_URL=.\+' .env; then
  echo "Isi DB_URL (MySQL Aiven) di .env"
  exit 1
fi

if grep -q 'tls=aiven' .env && [[ ! -f certs/ca.pem ]]; then
  echo "DB_URL memakai tls=aiven — letakkan ca.pem di deploy/certs/ atau pakai tls=skip-verify untuk dev open access"
  exit 1
fi

echo "==> Build & start backend stack"
docker compose up -d --build

echo ""
docker compose ps

echo ""
echo "Menunggu backend..."
for i in {1..30}; do
  if curl -sf "http://127.0.0.1/api/v1/health" >/dev/null 2>&1; then
    echo "OK — backend siap"
    curl -s "http://127.0.0.1/api/v1/health"
    echo ""
    HOST=$(grep -E '^PUBLIC_HOST=' .env | cut -d= -f2- || echo "103-169-206-116.domainesia.io")
    echo "Akses dari luar: http://${HOST}/api/v1/health"
    exit 0
  fi
  sleep 2
done

echo "Backend belum merespons. Cek log:"
echo "  docker compose logs backend"
exit 1
