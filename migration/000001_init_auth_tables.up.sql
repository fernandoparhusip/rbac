CREATE EXTENSION IF NOT EXISTS pgcrypto; -- Useful for uuid_generate_v4() if needed later

-- 1. Action (CREATE, READ, UPDATE, DELETE, dll)
CREATE TABLE auth_action (
    id uuid NOT NULL,
    code varchar(50) NOT NULL UNIQUE,
    "name" varchar(100) NOT NULL,
    description text,
    CONSTRAINT auth_action_pkey PRIMARY KEY (id)
);

-- 2. Resource (Users, Orders, Inventory, dll)
CREATE TABLE auth_resource (
    id uuid NOT NULL,
    code varchar(100) NOT NULL UNIQUE,
    "name" varchar(150) NOT NULL,
    description text,
    CONSTRAINT auth_resource_pkey PRIMARY KEY (id)
);

-- 3. Permission (Gabungan Resource & Action -> e.g. order:read)
CREATE TABLE auth_permission (
    id uuid NOT NULL,
    resource_id uuid NOT NULL,
    action_id uuid NOT NULL,
    permission_key varchar(200) NOT NULL UNIQUE, -- Contoh: 'order:read'
    description text,
    CONSTRAINT auth_permission_pkey PRIMARY KEY (id),
    CONSTRAINT fk_permission_action FOREIGN KEY (action_id) REFERENCES auth_action(id) ON DELETE CASCADE,
    CONSTRAINT fk_permission_resource FOREIGN KEY (resource_id) REFERENCES auth_resource(id) ON DELETE CASCADE
);

-- 4. Role (SuperAdmin, Admin, Staff)
CREATE TABLE auth_role (
    id uuid NOT NULL,
    code varchar(50) NOT NULL UNIQUE,
    "name" varchar(100) NOT NULL,
    description text,
    is_system boolean DEFAULT false,
    created_at timestamp DEFAULT current_timestamp,
    updated_at timestamp,
    CONSTRAINT auth_role_pkey PRIMARY KEY (id)
);

-- 5. Role Inheritance (SuperAdmin mewarisi izin Admin)
CREATE TABLE auth_role_inheritance (
    id uuid NOT NULL,
    parent_role_id uuid NOT NULL,
    child_role_id uuid NOT NULL,
    created_at timestamp DEFAULT current_timestamp,
    CONSTRAINT auth_role_inheritance_pkey PRIMARY KEY (id),
    CONSTRAINT uq_role_inheritance UNIQUE (parent_role_id, child_role_id),
    CONSTRAINT chk_role_inheritance CHECK (parent_role_id <> child_role_id),
    CONSTRAINT fk_parent_role FOREIGN KEY (parent_role_id) REFERENCES auth_role(id) ON DELETE CASCADE,
    CONSTRAINT fk_child_role FOREIGN KEY (child_role_id) REFERENCES auth_role(id) ON DELETE CASCADE
);

-- 6. Scope Type (Regional, Branch, Department)
CREATE TABLE auth_scope_type (
    id uuid NOT NULL,
    code varchar(50) NOT NULL UNIQUE,
    "name" varchar(100) NOT NULL,
    description text,
    CONSTRAINT auth_scope_type_pkey PRIMARY KEY (id)
);

-- 7. Scope (Data aktualnya: "Jawa Barat", "Cabang Bandung")
CREATE TABLE auth_scope (
    id uuid NOT NULL,
    scope_type_id uuid NOT NULL,
    parent_id uuid, -- Hierarchy: Regional -> Branch
    code varchar(100) NOT NULL UNIQUE,
    "name" varchar(150) NOT NULL,
    description text,
    CONSTRAINT auth_scope_pkey PRIMARY KEY (id),
    CONSTRAINT fk_scope_type FOREIGN KEY (scope_type_id) REFERENCES auth_scope_type(id) ON DELETE CASCADE,
    CONSTRAINT fk_scope_parent FOREIGN KEY (parent_id) REFERENCES auth_scope(id) ON DELETE CASCADE
);

-- 8. Role Permission (Menghubungkan Role ke Permission)
CREATE TABLE auth_role_permission (
    id uuid NOT NULL,
    role_id uuid NOT NULL,
    permission_id uuid NOT NULL,
    granted_at timestamp DEFAULT current_timestamp,
    CONSTRAINT auth_role_permission_pkey PRIMARY KEY (id),
    CONSTRAINT uq_role_permission UNIQUE (role_id, permission_id),
    CONSTRAINT fk_role_permission_role FOREIGN KEY (role_id) REFERENCES auth_role(id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permission_perm FOREIGN KEY (permission_id) REFERENCES auth_permission(id) ON DELETE CASCADE
);

-- 9. Role Permission Scope (Menghubungkan hak akses dengan scope spesifik)
-- Misal: Role X punya hak Y, tapi HANYA di Scope Z
CREATE TABLE auth_role_permission_scope (
    id uuid NOT NULL,
    role_permission_id uuid NOT NULL,
    scope_id uuid NOT NULL,
    CONSTRAINT auth_role_permission_scope_pkey PRIMARY KEY (id),
    CONSTRAINT uq_role_permission_scope UNIQUE (role_permission_id, scope_id),
    CONSTRAINT fk_role_perm_scope_rp FOREIGN KEY (role_permission_id) REFERENCES auth_role_permission(id) ON DELETE CASCADE,
    CONSTRAINT fk_role_perm_scope_scope FOREIGN KEY (scope_id) REFERENCES auth_scope(id) ON DELETE CASCADE
);

-- 10. User 
CREATE TABLE auth_user (
    id uuid NOT NULL,
    username varchar(150) NOT NULL UNIQUE,
    email varchar(254) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    full_name varchar(150) NOT NULL,
    status int2 NOT NULL DEFAULT 1,
    created_at timestamp DEFAULT current_timestamp,
    updated_at timestamp,
    CONSTRAINT auth_user_pkey PRIMARY KEY (id)
);

-- 11. User Role (Menghubungkan User ke Role)
CREATE TABLE auth_user_role (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    role_id uuid NOT NULL,
    assigned_at timestamp DEFAULT current_timestamp,
    CONSTRAINT auth_user_role_pkey PRIMARY KEY (id),
    CONSTRAINT uq_auth_user_role UNIQUE (user_id, role_id),
    CONSTRAINT fk_user_role_user FOREIGN KEY (user_id) REFERENCES auth_user(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_role_role FOREIGN KEY (role_id) REFERENCES auth_role(id) ON DELETE CASCADE
);

-- 12. User Permission (Pengecualian: User dikasih hak akses tanpa lewat Role)
CREATE TABLE auth_user_permission (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    permission_id uuid NOT NULL,
    is_granted boolean DEFAULT true, -- Bisa diset false untuk me-REVOKE hak akses tertentu dari seorang user
    granted_at timestamp DEFAULT current_timestamp,
    CONSTRAINT auth_user_permission_pkey PRIMARY KEY (id),
    CONSTRAINT uq_user_permission UNIQUE (user_id, permission_id),
    CONSTRAINT fk_user_perm_user FOREIGN KEY (user_id) REFERENCES auth_user(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_perm_perm FOREIGN KEY (permission_id) REFERENCES auth_permission(id) ON DELETE CASCADE
);

-- 13. Menu (Untuk Frontend)
CREATE TABLE auth_menu (
    id uuid NOT NULL,
    parent_id uuid,
    code varchar(100) NOT NULL UNIQUE,
    "name" varchar(150) NOT NULL,
    route varchar(255),
    icon varchar(100),
    sort_no int4 DEFAULT 0,
    CONSTRAINT auth_menu_pkey PRIMARY KEY (id),
    CONSTRAINT fk_menu_parent FOREIGN KEY (parent_id) REFERENCES auth_menu(id) ON DELETE CASCADE
);

-- 14. Role Menu (Menu apa saja yang bisa dilihat oleh suatu Role)
CREATE TABLE auth_role_menu (
    id uuid NOT NULL,
    role_id uuid NOT NULL,
    menu_id uuid NOT NULL,
    CONSTRAINT auth_role_menu_pkey PRIMARY KEY (id),
    CONSTRAINT uq_role_menu UNIQUE (role_id, menu_id),
    CONSTRAINT fk_role_menu_role FOREIGN KEY (role_id) REFERENCES auth_role(id) ON DELETE CASCADE,
    CONSTRAINT fk_role_menu_menu FOREIGN KEY (menu_id) REFERENCES auth_menu(id) ON DELETE CASCADE
);

-- 15. API (Pemetaan Endpoint ke Resource & Action)
CREATE TABLE auth_api (
    id uuid NOT NULL,
    "method" varchar(10) NOT NULL,
    "path" varchar(255) NOT NULL,
    resource_id uuid NOT NULL,
    action_id uuid NOT NULL,
    CONSTRAINT auth_api_pkey PRIMARY KEY (id),
    CONSTRAINT uq_api UNIQUE ("method", "path"),
    CONSTRAINT fk_api_action FOREIGN KEY (action_id) REFERENCES auth_action(id) ON DELETE CASCADE,
    CONSTRAINT fk_api_resource FOREIGN KEY (resource_id) REFERENCES auth_resource(id) ON DELETE CASCADE
);

-- 16. Session (Manajemen Token & Audit)
CREATE TABLE auth_session (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    access_token_hash text NOT NULL,
    refresh_token_hash text NOT NULL,
    ip_address varchar(100),
    user_agent text,
    expired_at timestamp NOT NULL,
    is_revoked boolean DEFAULT false,
    created_at timestamp DEFAULT current_timestamp,
    CONSTRAINT auth_session_pkey PRIMARY KEY (id),
    CONSTRAINT fk_session_user FOREIGN KEY (user_id) REFERENCES auth_user(id) ON DELETE CASCADE
);
