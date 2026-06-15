# Auth Service API Documentation

Base URL: `/api/v1`

## Authentication Endpoints

### Register User
POST `/auth/register`

Body:
```json
{
  "username": "admin",
  "email": "admin@example.com",
  "password": "Password123!",
  "full_name": "Admin User"
}
```

Success response:
```json
{
  "success": true,
  "message": "user registered",
  "data": {
    "id": "uuid",
    "username": "admin",
    "email": "admin@example.com",
    "full_name": "Admin User",
    "status": "active"
  }
}
```

### Login User
POST `/auth/login`

Body:
```json
{
  "username": "admin",
  "password": "Password123!",
}
```

`client_id` is optional. If provided, the login token will be linked to the client.

Success response:
```json
{
  "success": true,
  "message": "login success",
  "data": {
    "access_token": "...",
    "token_type": "Bearer",
    "expires_in": 900,
    "refresh_token": "...",
    "user": {
      "id": "uuid",
      "username": "admin",
      "email": "admin@example.com",
      "full_name": "Admin User",
      "status": "active"
    }
  }
}
```

### Refresh Token
POST `/auth/refresh`

Body:
```json
{
  "refresh_token": "..."
}
```

Success response:
```json
{
  "success": true,
  "message": "token refreshed",
  "data": {
    "access_token": "...",
    "token_type": "Bearer",
    "expires_in": 900,
    "refresh_token": "..."
  }
}
```

### Introspect Token
POST `/auth/introspect`

Body:
```json
{
  "token": "..."
}
```

This endpoint checks the access token (JWT) validity, not the refresh token.
If you refresh a token, the old refresh token becomes invalid, but the existing access token remains active until it expires.

Success response:
```json
{
  "success": true,
  "message": "token introspected",
  "data": {
    "active": true,
    "user": {
      "sub": "uuid",
      "iss": "app-name"
    },
    "roles": ["user"],
    "permissions": []
  }
}
```

### Logout
POST `/auth/logout`

Headers:
- `Authorization: Bearer <access_token>`

Success response:
```json
{
  "success": true,
  "message": "logout success"
}
```

### Get User Permissions
GET `/auth/permissions`

Headers:
- `Authorization: Bearer <access_token>`

Success response:
```json
{
  "success": true,
  "message": "user permissions retrieved",
  "data": {
    "permissions": ["all:access"]
  }
}
```

### Authorize
POST `/auth/authorize`

Check if the authenticated user has a specific permission. This endpoint validates permissions from the database based on user roles.

**Important:** User must have a role assignment in the `user_roles` table for permission lookup to work. If a user has no roles, they will have no permissions.

Headers:
- `Authorization: Bearer <access_token>`

Body:
```json
{
  "user_id": "uuid",
  "permission": "student:read"
}
```

Success response (authorized):
```json
{
  "success": true,
  "message": "authorization result",
  "data": {
    "authorized": true,
    "user_id": "uuid",
    "permission": "student:read"
  }
}
```

Error response (not authorized):
```json
{
  "success": true,
  "message": "authorization result",
  "data": {
    "authorized": false,
    "user_id": "uuid",
    "permission": "student:read"
  }
}
```

**Note:** The authorization check queries the database to verify permissions assigned via user roles, ensuring permission changes take effect immediately without waiting for token expiration. To assign permissions to a user:

1. Create or use an existing `role`
2. Link `permission` codes to the `role` via `role_permissions` table
3. Link the `user` to the `role` via `user_roles` table

**Example data flow:**
- User `student:read` permission comes from:
  - User → `user_roles` → Role "student" → `role_permissions` → Permission `student:read`

## OAuth Token Endpoint

### Request Token
POST `/oauth/token`

Body (password grant):
```json
{
  "grant_type": "password",
  "username": "admin",
  "password": "Password123!",
  "client_id": "default-client",
  "client_secret": "secret123"
}
```

Body (refresh_token grant):
```json
{
  "grant_type": "refresh_token",
  "refresh_token": "...",
  "client_id": "default-client",
  "client_secret": "secret123"
}
```

Success response:
```json
{
  "success": true,
  "message": "token granted",
  "data": {
    "access_token": "...",
    "token_type": "Bearer",
    "expires_in": 900,
    "refresh_token": "..."
  }
}
```

## User Endpoints

### Get User by ID
GET `/users/{id}`

Headers:
- `Authorization: Bearer <access_token>`

### Get All Users
GET `/users`

Headers:
- `Authorization: Bearer <access_token>`

### Update User
PUT `/users/{id}`

Headers:
- `Authorization: Bearer <access_token>`

Body:
```json
{
  "full_name": "New Name",
  "email": "new@example.com"
}
```

### Delete User
DELETE `/users/{id}`

Headers:
- `Authorization: Bearer <access_token>`

## Role Endpoints

### Create Role
POST `/roles`

Headers:
- `Authorization: Bearer <access_token>`

Body:
```json
{
  "code": "admin",
  "name": "admin",
  "description": "Administrator role"
}
```

### Get Role by ID
GET `/roles/{id}`

Headers:
- `Authorization: Bearer <access_token>`

### Get All Roles
GET `/roles`

Headers:
- `Authorization: Bearer <access_token>`

### Update Role
PUT `/roles/{id}`

Headers:
- `Authorization: Bearer <access_token>`

Body:
```json
{
  "code": "admin",
  "name": "admin",
  "description": "Updated description"
}
```

### Delete Role
DELETE `/roles/{id}`

Headers:
- `Authorization: Bearer <access_token>`

## Permission Endpoints

### Create Permission
POST `/permissions`

Headers:
- `Authorization: Bearer <access_token>`

Body:
```json
{
  "code": "all:access",
  "name": "all:access",
  "description": "Full access permission"
}
```

### Get Permission by ID
GET `/permissions/{id}`

Headers:
- `Authorization: Bearer <access_token>`

### Get All Permissions
GET `/permissions`

Headers:
- `Authorization: Bearer <access_token>`

### Update Permission
PUT `/permissions/{id}`

Headers:
- `Authorization: Bearer <access_token>`

Body:
```json
{
  "code": "all:access",
  "name": "all:access",
  "description": "Updated description"
}
```

### Delete Permission
DELETE `/permissions/{id}`

Headers:
- `Authorization: Bearer <access_token>`

## Client Endpoints

### Create Client
POST `/clients`

Headers:
- `Authorization: Bearer <access_token>`

Body:
```json
{
  "client_id": "default-client",
  "name": "Default Client",
  "client_secret": "secret123",
  "redirect_uris": "http://localhost/callback",
  "grants": "password,refresh_token",
  "is_confidential": true
}
```

### Get Client by ID
GET `/clients/{id}`

Headers:
- `Authorization: Bearer <access_token>`

### Get All Clients
GET `/clients`

Headers:
- `Authorization: Bearer <access_token>`

### Update Client
PUT `/clients/{id}`

Headers:
- `Authorization: Bearer <access_token>`

Body:
```json
{
  "client_id": "default-client",
  "name": "Default Client",
  "client_secret": "secret123",
  "redirect_uris": "http://localhost/callback",
  "grants": "password,refresh_token",
  "is_confidential": true
}
```

### Delete Client
DELETE `/clients/{id}`

Headers:
- `Authorization: Bearer <access_token>`
