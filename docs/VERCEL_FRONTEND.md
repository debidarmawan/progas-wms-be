# Frontend di Vercel — koneksi ke API VPS

Backend PROGAS WMS berjalan di VPS (port **3131**). Frontend (Next.js atau SPA) di-deploy di **Vercel** dan hanya perlu tahu URL API publik.

## Environment di Vercel

Project → **Settings** → **Environment Variables**:

| Variable | Contoh | Environment |
|----------|--------|-------------|
| `NEXT_PUBLIC_API_URL` | `http://203.0.113.10:3131/api/v1` | Production, Preview, Development |

Ganti IP dengan IP/domain VPS Anda. **Tanpa** trailing slash setelah `v1`.

Jika nanti API pakai HTTPS + domain:

```
NEXT_PUBLIC_API_URL=https://api.progas.co.id/api/v1
```

Update juga `CORS_ORIGINS` di `.env` VPS agar match origin Vercel.

## CORS di backend (VPS)

Di `.env` VPS:

```env
# Production — ganti dengan URL project Vercel Anda
CORS_ORIGINS=https://progas-wms.vercel.app

# Preview deploy Vercel (branch PR)
CORS_ALLOW_VERCEL_PREVIEW=true
```

Beberapa domain (production + custom):

```env
CORS_ORIGINS=https://app.progas.co.id,https://progas-wms.vercel.app
```

Setelah ubah `.env`:

```bash
docker compose up -d --build backend
```

## Contoh pemanggilan API (Next.js)

```typescript
const base = process.env.NEXT_PUBLIC_API_URL;

export async function login(email: string, password: string) {
  const res = await fetch(`${base}/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
  return res.json();
}
```

Request authenticated:

```typescript
headers: {
  Authorization: `Bearer ${accessToken}`,
  "Content-Type": "application/json",
}
```

## HTTP vs HTTPS

- Vercel = **HTTPS**
- API VPS default = **HTTP** `:3131`

Fetch dari browser HTTPS ke HTTP API **diperbolehkan** (bukan mixed content untuk API call). Data tetap terenkripsi di lapisan TLS Vercel↔user; leg VPS↔Vercel edge tidak TLS kecuali Anda pasang HTTPS di VPS.

Untuk production jangka panjang, pertimbangkan:

1. Domain API + Let's Encrypt di VPS, atau  
2. Cloudflare proxy / tunnel ke port 3131

Lalu set `NEXT_PUBLIC_API_URL=https://api.domain.com/api/v1`.

## Checklist deploy

- [ ] VPS: `docker compose up -d --build` + health OK
- [ ] VPS: firewall buka **3131**
- [ ] VPS: `CORS_ORIGINS` = URL Vercel production
- [ ] Vercel: `NEXT_PUBLIC_API_URL` = `http://IP_VPS:3131/api/v1`
- [ ] Test login dari URL Vercel (DevTools → Network, tidak ada CORS error)

## Swagger

Swagger hanya aktif jika `GO_ENV=development`. Di production VPS biarkan `GO_ENV=production`; dokumentasi API via repo atau generate client dari `docs/swagger.yaml`.
