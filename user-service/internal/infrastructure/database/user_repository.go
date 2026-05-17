package database

import (
	"context"
	"errors"

	"user-service/internal/domain/entities"
	"user-service/internal/domain/services"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users(id, name, email, password_hash, phone, address, created_at)
		VALUES($1,$2,$3,$4,$5,$6,$7)
	`
	_, err := r.db.Exec(ctx, query,
		user.ID, user.Name, user.Email, user.PasswordHash,
		user.Phone, user.Address, user.CreatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return services.ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	query := `
		SELECT id, name, email, password_hash, phone, address, created_at
		FROM users WHERE id = $1
	`
	var user entities.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Phone, &user.Address, &user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, services.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT id, name, email, password_hash, phone, address, created_at
		FROM users WHERE email = $1
	`
	var user entities.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Phone, &user.Address, &user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, services.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users SET name = $2, phone = $3, address = $4
		WHERE id = $1
	`
	tag, err := r.db.Exec(ctx, query, user.ID, user.Name, user.Phone, user.Address)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return services.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context) ([]entities.User, error) {
	query := `
		SELECT id, name, email, password_hash, phone, address, created_at
		FROM users ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.PasswordHash,
			&user.Phone, &user.Address, &user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func isUniqueViolation(err error) bool {
	var pgErr interface{ Code() string }
	return errors.As(err, &pgErr) && pgErr.Code() == "23505"
}
