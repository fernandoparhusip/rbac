-- Menambahkan kolom role ke tabel auth_user untuk mendukung Role-Based Access Control (RBAC)
ALTER TABLE auth_user ADD COLUMN IF NOT EXISTS role varchar(50) NOT NULL DEFAULT 'User';
