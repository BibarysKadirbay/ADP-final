package database

import (
	"context"

	"user-service/internal/domain/entities"
	"user-service/internal/domain/services"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(
	db *pgxpool.Pool,
) *UserRepository {

	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *entities.User,
) error {

	query := `
		INSERT INTO users(
			id,
			name,
			email,
			phone,
			address,
			created_at
		)
		VALUES($1,$2,$3,$4,$5,$6)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Phone,
		user.Address,
		user.CreatedAt,
	)

	return err
}

func (r *UserRepository) GetByID(
	ctx context.Context,
	id string,
) (*entities.User, error) {

	query := `
		SELECT
			id,
			name,
			email,
			phone,
			address,
			created_at
		FROM users
		WHERE id = $1
	`

	var user entities.User

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Address,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, services.ErrUserNotFound
	}

	return &user, nil
}

func (r *UserRepository) List(
	ctx context.Context,
) ([]entities.User, error) {

	query := `
		SELECT
			id,
			name,
			email,
			phone,
			address,
			created_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []entities.User

	for rows.Next() {

		var user entities.User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Address,
			&user.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
