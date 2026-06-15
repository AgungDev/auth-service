Tambahkan Swagger/OpenAPI Documentation ke project ummu-auth-service.

Project Information:
- Language: Go
- Framework: Gin
- Architecture: Clean Architecture
- Module Name: ummu-auth-service
- Database: PostgreSQL
- Authentication: JWT Access Token & Refresh Token

Requirements:

1. Gunakan library:
   - github.com/swaggo/swag
   - github.com/swaggo/gin-swagger
   - github.com/swaggo/files

2. Install dan konfigurasi Swagger untuk Gin.

3. Buat endpoint dokumentasi:
   GET /swagger/*any

4. Tambahkan informasi API:
   - Title: UMMU Auth Service API
   - Version: 1.0.0
   - Description: Authentication and Authorization Service for UMMU Microservices
   - Contact Name: UMMU Development Team

5. Tambahkan anotasi Swagger pada seluruh endpoint Auth:
   - Register
   - Login
   - Refresh Token
   - Logout
   - Profile
   - Verify Token
   - Check Permission

6. Dokumentasikan:
   - Request Body
   - Response Body
   - Success Response
   - Error Response
   - Authorization Header

7. Tambahkan Bearer Authentication Scheme JWT pada Swagger.

8. Buat model Swagger untuk:
   - RegisterRequest
   - LoginRequest
   - RefreshTokenRequest
   - AuthResponse
   - ErrorResponse
   - PermissionResponse

9. Pastikan Swagger dapat di-generate menggunakan:

   swag init -g cmd/server/main.go

10. Update README.md dengan:
    - Cara install Swagger
    - Cara generate docs
    - URL Swagger UI

Expected Result:
- Swagger UI dapat diakses pada:
  http://localhost:9001/swagger/index.html

- Semua endpoint Auth muncul pada Swagger UI.
- JWT Authorization dapat diuji langsung dari Swagger UI menggunakan tombol Authorize.
- Kode mengikuti struktur Clean Architecture yang sudah ada tanpa mengubah business logic.

Analisis seluruh struktur project ummu-auth-service terlebih dahulu.

Tambahkan Swagger/OpenAPI tanpa mengubah arsitektur yang ada.

Letakkan:
- konfigurasi swagger di cmd/server/main.go
- anotasi endpoint pada layer transport/http
- model dokumentasi di folder docs/swagger_models

Jangan memindahkan file yang sudah ada.
Jangan mengubah business logic.
Jangan mengubah repository dan usecase.

Pastikan seluruh route Gin yang sudah terdaftar otomatis terdokumentasi di Swagger.

Setelah selesai:
1. Generate docs Swagger.
2. Tambahkan command Makefile:
   make swagger
3. Update README.md.