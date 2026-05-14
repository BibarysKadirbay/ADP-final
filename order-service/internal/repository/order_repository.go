package repository

import (
	"database/sql"
	"order-service/internal/domain"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order domain.Order) error {

	_, err := r.db.Exec(
		"INSERT INTO orders (id, user_id, restaurant_id, total_price, status) VALUES ($1,$2,$3,$4,$5)",
		order.ID,
		order.UserID,
		order.RestaurantID,
		order.TotalPrice,
		order.Status,
	)

	return err
}

func (r *OrderRepository) GetByID(id string) (domain.Order, error) {

	var order domain.Order

	err := r.db.QueryRow(
		"SELECT id,user_id,restaurant_id,total_price,status FROM orders WHERE id=$1",
		id,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.RestaurantID,
		&order.TotalPrice,
		&order.Status,
	)

	return order, err
}

func (r *OrderRepository) UpdateStatus(id string, status string) error {

	_, err := r.db.Exec(
		"UPDATE orders SET status=$1 WHERE id=$2",
		status,
		id,
	)

	return err
}
