-- Dummy seed data for auth service
-- This file includes user, client, role, permission, and token seed rows.

INSERT INTO users (id, username, email, password_hash, full_name, status) VALUES
('11111111-1111-1111-1111-111111111111', 'admin', 'admin@example.com', '$2a$10$hs7AlCQqfYSHtrQddAu3uuCzKy82liEz4Hbj9FgnMTbsek4iWFpuu', 'Admin User', 'active'),
('22222222-2222-2222-2222-111111111111', 'student', 'student@example.com', '$2a$10$hs7AlCQqfYSHtrQddAu3uuCzKy82liEz4Hbj9FgnMTbsek4iWFpuu', 'Student User', 'active'),
('dd44cbfc-92ab-4f4f-9ac0-6bed055a95f1', 'newadmin', 'newadmin@example.com', '$2a$10$hs7AlCQqfYSHtrQddAu3uuCzKy82liEz4Hbj9FgnMTbsek4iWFpuu', 'New Admin User', 'active');

INSERT INTO clients (id, client_id, name, client_secret_hash, redirect_uris, grants, is_confidential) VALUES
('22222222-2222-2222-2222-222222222222', 'default-client', 'Default Client', '$2a$10$TxRKJE0Vvto0XFhTYSBB3uPlO1GUEOhxRbNn6t..aBzN34d73Ta0.', 'http://localhost/callback', 'password,refresh_token', true);

INSERT INTO roles (id, code, name, description) VALUES
('33333333-3333-3333-3333-333333333333', 'admin', 'admin', 'Administrator role'),
('44444444-4444-4444-4444-333333333333', 'student', 'student', 'Student role');

INSERT INTO permissions (id, code, name, description) VALUES
('55555555-5555-5555-5555-555555555555', 'all:access', 'all:access', 'Full access permission'),
('66666666-6666-6666-6666-666666666666', 'student:read', 'student:read', 'Read student data'),
('77777777-7777-7777-7777-777777777777', 'student:update', 'student:update', 'Update student data'),
('88888888-8888-8888-8888-888888888888', 'profile:read', 'profile:read', 'Read user profile');

INSERT INTO user_roles (user_id, role_id) VALUES
('11111111-1111-1111-1111-111111111111', '33333333-3333-3333-3333-333333333333'),
('22222222-2222-2222-2222-111111111111', '44444444-4444-4444-4444-333333333333'),
('dd44cbfc-92ab-4f4f-9ac0-6bed055a95f1', '33333333-3333-3333-3333-333333333333');

INSERT INTO role_permissions (role_id, permission_id) VALUES
('33333333-3333-3333-3333-333333333333', '55555555-5555-5555-5555-555555555555'),
('44444444-4444-4444-4444-333333333333', '66666666-6666-6666-6666-666666666666'),
('44444444-4444-4444-4444-333333333333', '77777777-7777-7777-7777-777777777777'),
('44444444-4444-4444-4444-333333333333', '88888888-8888-8888-8888-888888888888');

INSERT INTO tokens (user_id, client_id, refresh_token_hash, expires_at) VALUES
('11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', 'dummy-refresh-token', NOW() + INTERVAL '1 day'),
('22222222-2222-2222-2222-111111111111', '22222222-2222-2222-2222-222222222222', 'dummy-refresh-token-2', NOW() + INTERVAL '1 day'),
('dd44cbfc-92ab-4f4f-9ac0-6bed055a95f1', '22222222-2222-2222-2222-222222222222', 'dummy-refresh-token-3', NOW() + INTERVAL '1 day');
