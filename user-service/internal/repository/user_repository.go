package repository

import (
	"database/sql"
	"user-service/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user domain.User) error {
	_, err := r.db.Exec(
		"INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)",
		user.ID,
		user.Name,
		user.Email,
		user.Password,
	)
	return err
}

func (r *UserRepository) GetByEmail(email string) (domain.User, error) {
	var user domain.User

	err := r.db.QueryRow(
		"SELECT id, name, email, password FROM users WHERE email=$1",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	return user, err
}

func (r *UserRepository) GetByID(id string) (domain.User, error) {
	var user domain.User

	err := r.db.QueryRow(
		"SELECT id, name, email, password FROM users WHERE id=$1",
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	return user, err
}

func (r *UserRepository) Update(user domain.User) error {
	_, err := r.db.Exec(
		"UPDATE users SET name=$1, email=$2 WHERE id=$3",
		user.Name,
		user.Email,
		user.ID,
	)
	return err
}
