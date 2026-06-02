CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =====================================================
-- USERS
-- =====================================================

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(64) UNIQUE NOT NULL,
    email VARCHAR(128) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(128),
    status VARCHAR(20) NOT NULL DEFAULT 'inactive'
        CHECK (status IN (
            'active',
            'inactive',
            'suspended',
            'locked'
        )),

    last_login_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- ROLES
-- =====================================================

CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(64) UNIQUE NOT NULL,
    name VARCHAR(128) NOT NULL,
    description TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- PERMISSIONS
-- =====================================================

CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(128) UNIQUE NOT NULL,
    name VARCHAR(128) NOT NULL,
    description TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- USER ROLES (MANY TO MANY)
-- =====================================================

CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,

    PRIMARY KEY (user_id, role_id),

    CONSTRAINT fk_user_roles_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_user_roles_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE
);

-- =====================================================
-- ROLE PERMISSIONS (MANY TO MANY)
-- =====================================================

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,

    PRIMARY KEY (role_id, permission_id),

    CONSTRAINT fk_role_permissions_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_role_permissions_permission
        FOREIGN KEY (permission_id)
        REFERENCES permissions(id)
        ON DELETE CASCADE
);

-- =====================================================
-- USER PERMISSIONS
-- DIRECT PERMISSION OVERRIDE
-- =====================================================

CREATE TABLE IF NOT EXISTS user_permissions (
    user_id UUID NOT NULL,
    permission_id UUID NOT NULL,

    PRIMARY KEY (user_id, permission_id),

    CONSTRAINT fk_user_permissions_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_user_permissions_permission
        FOREIGN KEY (permission_id)
        REFERENCES permissions(id)
        ON DELETE CASCADE
);

-- =====================================================
-- OAUTH CLIENTS
-- =====================================================

CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id VARCHAR(128) UNIQUE NOT NULL,
    client_secret_hash VARCHAR(255) NOT NULL,
    name VARCHAR(128) NOT NULL,
    is_confidential BOOLEAN NOT NULL DEFAULT TRUE,
    status VARCHAR(20) NOT NULL DEFAULT 'inactive'
        CHECK (status IN (
            'active',
            'inactive'
        )),

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS client_grants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL,
    grant_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_client_grants_client
        FOREIGN KEY (client_id)
        REFERENCES clients(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_client_grants
        UNIQUE (client_id, grant_type)
);

CREATE TABLE IF NOT EXISTS client_redirect_uris (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL,
    redirect_uri TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_client_redirect_uris_client
        FOREIGN KEY (client_id)
        REFERENCES clients(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_client_redirect_uri
        UNIQUE (client_id, redirect_uri)
);

-- =====================================================
-- REFRESH TOKENS
-- ACCESS TOKEN = JWT (STATELESS)
-- =====================================================

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    client_id UUID,
    refresh_token_hash TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_refresh_tokens_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_refresh_tokens_client
        FOREIGN KEY (client_id)
        REFERENCES clients(id)
        ON DELETE CASCADE
);

-- =====================================================
-- INDEXES
-- =====================================================

CREATE INDEX IF NOT EXISTS idx_users_username
ON users(username);

CREATE INDEX IF NOT EXISTS idx_users_email
ON users(email);

CREATE INDEX IF NOT EXISTS idx_user_roles_user
ON user_roles(user_id);

CREATE INDEX IF NOT EXISTS idx_user_roles_role
ON user_roles(role_id);

CREATE INDEX IF NOT EXISTS idx_role_permissions_role
ON role_permissions(role_id);

CREATE INDEX IF NOT EXISTS idx_role_permissions_permission
ON role_permissions(permission_id);

CREATE INDEX IF NOT EXISTS idx_user_permissions_user
ON user_permissions(user_id);

CREATE INDEX IF NOT EXISTS idx_user_permissions_permission
ON user_permissions(permission_id);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user
ON refresh_tokens(user_id);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_client
ON refresh_tokens(client_id);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires
ON refresh_tokens(expires_at);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_revoked
ON refresh_tokens(is_revoked);

CREATE INDEX IF NOT EXISTS idx_clients_client_id
ON clients(client_id);

CREATE INDEX IF NOT EXISTS idx_client_grants_client
ON client_grants(client_id);

CREATE INDEX IF NOT EXISTS idx_client_grants_type
ON client_grants(grant_type);

CREATE INDEX IF NOT EXISTS idx_client_redirect_uris_client
ON client_redirect_uris(client_id);