-- ==========================================
-- 1. TABEL USER (Induk)
-- ==========================================
CREATE TABLE auth_user (
    id uuid NOT NULL,
    username varchar(150) NOT NULL UNIQUE,
    email varchar(254) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    full_name varchar(150) NOT NULL,
    
    status int2 NOT NULL DEFAULT 1, -- 1: Aktif, 0: Suspend
    
    created_at timestamp DEFAULT current_timestamp,
    updated_at timestamp,
    
    -- Menjadikan kolom 'id' sebagai Primary Key
    CONSTRAINT auth_user_pkey PRIMARY KEY (id)
);


-- ==========================================
-- 2. TABEL SESSION (Anak)
-- ==========================================
CREATE TABLE auth_session (
    id uuid NOT NULL,
    session_uuid uuid NOT NULL UNIQUE,
    user_id uuid NOT NULL,    -- Kolom ini akan dihubungkan ke Tabel User
    
    access_jti varchar(255) NOT NULL UNIQUE,
    refresh_jti varchar(255) NOT NULL UNIQUE,
    refresh_hash text NOT NULL,
    
    expired_at timestamp NOT NULL,
    logout_at timestamp,
    revoked boolean DEFAULT false,
    created_at timestamp DEFAULT current_timestamp,
    
    -- Menjadikan kolom 'id' sebagai Primary Key
    CONSTRAINT auth_session_pkey PRIMARY KEY (id),
    
    -- [ INILAH FOREIGN KEY-NYA ]
    -- Menghubungkan 'user_id' di tabel ini dengan 'id' di tabel auth_user.
    -- ON DELETE CASCADE: Artinya jika User dihapus, semua riwayat sesinya ikut terhapus otomatis!
    CONSTRAINT fk_session_user FOREIGN KEY (user_id) REFERENCES auth_user(id) ON DELETE CASCADE
);

-- ==========================================
-- 3. TABEL ROLE (Induk Peran)
-- ==========================================
CREATE TABLE auth_role (
    id uuid NOT NULL,
    name varchar(100) NOT NULL UNIQUE,
    CONSTRAINT auth_role_pkey PRIMARY KEY (id)
);

-- ==========================================
-- 4. TABEL USER ROLE (Relasi Penghubung)
-- ==========================================
CREATE TABLE auth_user_role (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    role_id uuid NOT NULL,
    CONSTRAINT auth_user_role_pkey PRIMARY KEY (id),
    CONSTRAINT fk_user_role_user FOREIGN KEY (user_id) REFERENCES auth_user(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_role_role FOREIGN KEY (role_id) REFERENCES auth_role(id) ON DELETE CASCADE
);

-- Seed data role bawaan
INSERT INTO auth_role (id, name) VALUES 
    ('a0000000-0000-0000-0000-000000000001', 'Admin'),
    ('a0000000-0000-0000-0000-000000000002', 'User')
ON CONFLICT (name) DO NOTHING;

-- Index Kecepatan Tinggi untuk pencarian Middleware
CREATE INDEX idx_session_access_jti ON auth_session(access_jti);
