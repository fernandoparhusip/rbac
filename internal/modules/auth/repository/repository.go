package repository

import (
	"context"
	"errors"
	"time"

	"rbac/internal/modules/auth/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) AuthRepository {
	return &repository{
		db: db,
	}
}

// FindByEmail mencari user berdasarkan email dan menggabungkan datanya dengan role yang dimiliki
func (r *repository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User

	query := `
		SELECT
			U.id,
			U.username,
			U.email,
			U.password,
			U.full_name,
			COALESCE(R.name, 'User') AS role,
			U.status,
			U.created_at,
			U.updated_at
		FROM auth_user U
		LEFT JOIN auth_user_role UR ON U.id = UR.user_id
		LEFT JOIN auth_role R ON UR.role_id = R.id
		WHERE U.email = $1
		LIMIT 1
	`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FullName,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

// CreateUser memasukkan data user baru dan mengaitkannya ke tabel auth_user_role dalam satu transaksi
func (r *repository) CreateUser(ctx context.Context, user domain.User, roleID *uuid.UUID) (domain.User, error) {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	if user.Status == 0 {
		user.Status = 1
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	// Gunakan role default (User) jika roleID tidak ditentukan
	if roleID == nil || *roleID == uuid.Nil {
		defaultID, err := r.GetDefaultRoleID(ctx)
		if err != nil {
			return user, errors.New("gagal mendapatkan default role: " + err.Error())
		}
		roleID = &defaultID
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return user, err
	}
	defer tx.Rollback(ctx)

	// 1. Simpan user ke tabel auth_user
	queryUser := `
		INSERT INTO auth_user (
			id, username, email, password, full_name, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, username, email, password, full_name, status, created_at, updated_at
	`

	err = tx.QueryRow(
		ctx,
		queryUser,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.FullName,
		user.Status,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FullName,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}

	// 2. Hubungkan user dengan role di auth_user_role
	queryUserRole := `
		INSERT INTO auth_user_role (id, user_id, role_id)
		VALUES ($1, $2, $3)
	`
	userRoleID := uuid.New()
	_, err = tx.Exec(ctx, queryUserRole, userRoleID, user.ID, *roleID)
	if err != nil {
		return user, err
	}

	// 3. Ambil nama role yang baru saja dipasangkan
	queryRoleName := `SELECT name FROM auth_role WHERE id = $1`
	err = tx.QueryRow(ctx, queryRoleName, *roleID).Scan(&user.Role)
	if err != nil {
		return user, err
	}

	if err := tx.Commit(ctx); err != nil {
		return user, err
	}

	return user, nil
}

// AssignRoleToUser menambahkan penugasan role ke user
func (r *repository) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) (domain.User, error) {
	queryInsert := `
		INSERT INTO auth_user_role (id, user_id, role_id)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(ctx, queryInsert, uuid.New(), userID, roleID)
	if err != nil {
		return domain.User{}, err
	}

	// Ambil kembali user dengan role yang baru
	var user domain.User
	queryUser := `
		SELECT
			U.id, U.username, U.email, U.password, U.full_name, COALESCE(R.name, 'User') AS role,
			U.status, U.created_at, U.updated_at
		FROM auth_user U
		LEFT JOIN auth_user_role UR ON U.id = UR.user_id
		LEFT JOIN auth_role R ON UR.role_id = R.id
		WHERE U.id = $1
		LIMIT 1
	`
	err = r.db.QueryRow(ctx, queryUser, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FullName,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

// GetDefaultRoleID mengambil ID dari role default 'User'
func (r *repository) GetDefaultRoleID(ctx context.Context) (uuid.UUID, error) {
	var roleID uuid.UUID
	query := `SELECT id FROM auth_role WHERE name = 'User' LIMIT 1`

	err := r.db.QueryRow(ctx, query).Scan(&roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Fallback: jika belum ada data 'User', coba buatkan
			roleID = uuid.New()
			_, errInsert := r.db.Exec(ctx, "INSERT INTO auth_role (id, name) VALUES ($1, 'User')", roleID)
			if errInsert == nil {
				return roleID, nil
			}
		}
		return uuid.Nil, err
	}

	return roleID, nil
}

// GetRoleByName mencari role berdasarkan nama
func (r *repository) GetRoleByName(ctx context.Context, name string) (domain.Role, error) {
	var role domain.Role
	query := `SELECT id, name FROM auth_role WHERE LOWER(name) = LOWER($1) LIMIT 1`

	err := r.db.QueryRow(ctx, query, name).Scan(&role.ID, &role.Name)
	if err != nil {
		return role, err
	}

	return role, nil
}

// GetRoleByID mencari role berdasarkan UUID ID
func (r *repository) GetRoleByID(ctx context.Context, id uuid.UUID) (domain.Role, error) {
	var role domain.Role
	query := `SELECT id, name FROM auth_role WHERE id = $1 LIMIT 1`

	err := r.db.QueryRow(ctx, query, id).Scan(&role.ID, &role.Name)
	if err != nil {
		return role, err
	}

	return role, nil
}
