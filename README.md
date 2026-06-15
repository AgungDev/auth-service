# auth_service — Layanan Otentikasi untuk Sistem Informasi Akademik

> Auth service khusus buat Universitas Muhammadiyah Maluku Utara. Semua proses login / otorisasi lewat sini. Ringkas, aman (usahakan), dan modular.

---

## Ringkasan singkat

`auth_service` adalah microservice bertanggung jawab penuh atas:

* Otentikasi (login) pengguna
* Otorisasi (RBAC — role & permission)
* Mengeluarkan dan memverifikasi token (OAuth2 flows + JWT)
* Manajemen client application (API clients yang dipakai internal seperti `gateway`, `portal-mahasiswa`, `admin-panel`)
* Endpoint introspeksi/validasi token untuk service lain

Tech stack utama:

* Golang (minimal Go 1.20)
* Clean Architecture (domain/usecase/repository/controller)
* GORM (ORM, PostgreSQL)
* OAuth2 (go-oauth2 / ory hydra sebagai referensi — implementasi di service: Authorization Code + Password + Refresh)
* JWT (signed using RS256 — gunakan keypair)
* Docker + docker-compose untuk dev

---

## Tujuan & Ruang Lingkup

* Menyediakan central auth untuk seluruh microservices kampus
* Menegakkan kebijakan akses: siapa boleh akses apa
* Mudah di-integrasikan oleh API Gateway (mis. `api_gateway` akan memanggil endpoint introspect)
* Fase awal: internal clients + password grant untuk legacy. Fase produksi: pakai Authorization Code PKCE untuk UI web/mobile.

---

## Konsep domain (secara ringkas)

* **User** — entitas pengguna (student, lecturer, staff, superadmin)
* **Role** — kumpulan permissions (contoh: `student`, `dosen`, `admin_fakultas`, `superadmin`)
* **Permission** — aksi spesifik (contoh: `grades:read`, `grades:write`, `courses:manage`)
* **Client** — aplikasi yang memakai auth (gateway, mobile app)
* **Token** — access token (JWT), refresh token (opaque atau juga JWT)

---

## Struktur repository (saran)

```
auth_service/
├─ cmd/
│  └─ server/             # entrypoint
├─ internal/
│  ├─ domain/             # entity definitions
│  ├─ repository/         # GORM implementations
│  ├─ usecase/            # business logic
│  ├─ transport/
│  │  ├─ http/            # http handlers, middlewares
│  │  └─ grpc/ (opsional)
│  └─ config/             # env + config loader
├─ migrations/            # sql / goose / golang-migrate
├─ docker/
├─ scripts/
├─ pkg/                   # shared helpers (jwt, crypto)
└─ Makefile
```

---

## Database schema (singkat)

Tabel utama:

* `users` (id, username, email, password_hash, full_name, is_active, created_at, updated_at)
* `roles` (id, name, description)
* `permissions` (id, name, description)
* `role_permissions` (role_id, permission_id)
* `user_roles` (user_id, role_id)
* `clients` (id, client_id, client_secret_hash, redirect_uris, grants, is_confidential)
* `tokens` (id, user_id, client_id, refresh_token_hash, expires_at)

> Tips: jangan simpan password atau client secret secara plain — hash dengan argon2id / bcrypt.

---

## Alur Otentikasi & Otorisasi (singkat)

1. **Login (Password grant - dev/legacy)**

   * Client -> `POST /oauth/token` dengan grant_type=password, username, password, client_id, client_secret
   * Service memverifikasi kredensial -> keluarkan `access_token` (JWT) + `refresh_token`
2. **Authorization Code (production UI)**

   * Client (browser) redirect ke `/oauth/authorize` -> user login -> redirect kembali dengan `code`
   * Client tukarkan code ke `/oauth/token` -> dapat access + refresh
3. **Validasi token oleh service lain**

   * API Gateway/service memanggil `/introspect` atau memverifikasi JWT menggunakan public key (recommended)

Security notes: gunakan HTTPS, RS256 untuk JWT (private key aman di KMS / vault), expirasi singkat untuk access token (mis. 15m - 1h), refresh token lebih panjang tapi revokable.

---

## RBAC: contoh implementasi

Permissions: string namespace-style, contoh:

* `courses:read`
* `courses:write`
* `grades:read`
* `grades:write`
* `users:manage`

Role -> Permission mapping disimpan di `role_permissions`.

Saat membuat token, sertakan `roles` (atau `permissions`) di claim JWT untuk mempermudah pengecekan.

Contoh claims minimal:

```json
{
  "sub": "user:123",
  "iss": "auth.ummu-malu",
  "exp": 1712345678,
  "roles": [
    {
        "role": "student",
        "permissions": ["courses:read","profile:read"]
    },
  ],
}
```

---

## Endpoint API (draft)

### Public / OAuth endpoints

* `POST /oauth/token` — tukar kredensial / grant code => token
* `GET  /oauth/authorize` — endpoint authorize (user login + consent)
* `POST /oauth/revoke` — revocation
* `POST /oauth/introspect` — introspeksi token (opsional jika service lain tidak mau/ bisa verifikasi JWT)

### Management (butuh admin role)

* `GET  /api/v1/users` — list users
* `GET  /api/v1/users/{id}`
* `POST /api/v1/users` — create user
* `PUT  /api/v1/users/{id}` — update user
* `DELETE /api/v1/users/{id}` — deactivate user
* `GET  /api/v1/roles`
* `POST /api/v1/roles`
* `POST /api/v1/roles/{id}/permissions` — assign permission
* `GET  /api/v1/clients`
* `POST /api/v1/clients` — register client (for internal apps)
* `POST /api/v1/tokens/revoke`

> Semua endpoint management harus diproteksi oleh role `admin` / `superadmin`.

---

## Contoh .env (development)

```
APP_ENV=development
APP_PORT=8080
DATABASE_URL=postgres://auth:password@db:5432/authdb?sslmode=disable
JWT_PRIVATE_KEY_PATH=/run/secrets/jwt_private.pem
JWT_PUBLIC_KEY_PATH=/run/secrets/jwt_public.pem
ACCESS_TOKEN_EXPIRE_MINUTES=15
REFRESH_TOKEN_EXPIRE_DAYS=30
OAUTH_ALLOW_PASSWORD_GRANT=true
SIG_KEY_ALGO=RS256
```

---

## Docker & docker-compose (saran singkat)

* `docker-compose.yml` untuk dev: postgres, migration, auth_service
* Gunakan secrets untuk private key dan client secrets

---

## Swagger / OpenAPI Documentation

* Install Swagger CLI:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

* Install Gin Swagger packages:

```bash
go get github.com/swaggo/gin-swagger github.com/swaggo/files
```

* Generate Swagger docs:

```bash
make swagger
```

* On Windows PowerShell, use:

```powershell
.\swagger.ps1
```

* Swagger UI URL:

```text
http://localhost:9001/swagger/index.html
```

* Swagger route:

```text
/swagger/*any
```

---

## Migrasi & Seed

* Gunakan `golang-migrate` atau `goose`.
* Seed awal: roles (`student`, `dosen`, `admin_fakultas`, `superadmin`), beberapa permissions dasar, default `superadmin` user (random password — minta perubahan saat pertama login).

---

## Contoh alur pembuatan token (curl)

```
# Password grant (dev)
curl -X POST http://localhost:8080/oauth/token \
  -d "grant_type=password&username=alice&password=secret&client_id=try-to-schedule&client_secret=xxx"

# Response: {"access_token":"...","token_type":"Bearer","expires_in":900,"refresh_token":"..."}
```

---

## Tips implementasi - hal yang sering salah

1. **Jangan** sign JWT dengan HMAC using shared secret untuk production — gunakan RS256 dan simpan private key aman.
2. **Jangan** masukkan permissions terlalu banyak di token jika sensitif; cukup roles, dan gunakan introspect untuk permission-full check bila perlu.
3. **Jangan** simpan refresh token tanpa cara revoke; simpan memungkinkan blacklist/ revocation list.
4. **Selalu** gunakan HTTPS di production.

---

## Testing

* Unit test untuk usecase bisnis
* Integration test: run Postgres (docker-compose), jalankan migration, tes endpoint oauth/token, introspect
* Contract tests untuk API Gateway (walau nanti)

---

## CI / CD (singkat)

* Build image, run lint, run unit tests
* Scan vuln image
* Deploy ke staging, jalankan smoke tests

---

## Roadmap (fase)

1. **Fase 1 (MVP)**: Password grant, JWT RS256, user/role management, client management, introspect
2. **Fase 2**: Authorization Code + PKCE, consent screen, OIDC discovery (`/.well-known/openid-configuration`)
3. **Fase 3**: Integrasi KMS (key rotation), audit logs, 2FA (TOTP)

---

## References & Komponen cocok

* `github.com/go-oauth2/oauth2`
* `gorm.io/gorm`
* `github.com/golang-jwt/jwt/v5`
* `github.com/alexedwards/argon2id` (atau `bcrypt`)

