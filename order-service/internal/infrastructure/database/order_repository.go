package database

import (
	"context"

	"order-service/internal/domain/entities"
	"order-service/internal/domain/services"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *entities.Order) error {
	query := `
		INSERT INTO orders (
			id,
			user_id,
			restaurant_id,
			delivery_id,
			total_price,
			status,
			payment_status,
			address,
			comment,
			created_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		order.ID,
		order.UserID,
		order.RestaurantID,
		order.DeliveryID,
		order.TotalPrice,
		order.Status,
		order.PaymentStatus,
		order.Address,
		order.Comment,
		order.CreatedAt,
	)

	return err
}

func (r *OrderRepository) GetByID(ctx context.Context, id string) (*entities.Order, error) {
	query := `
		SELECT 
			id,
			user_id,
			restaurant_id,
			delivery_id,
			total_price,
			status,
			payment_status,
			address,
			comment,
			created_at
		FROM orders
		WHERE id = $1
	`

	var order entities.Order

	err := r.db.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.RestaurantID,
		&order.DeliveryID,
		&order.TotalPrice,
		&order.Status,
		&order.PaymentStatus,
		&order.Address,
		&order.Comment,
		&order.CreatedAt,
	)

	if err != nil {
		return nil, services.ErrOrderNotFound
	}

	return &order, nil
}

func (r *OrderRepository) List(ctx context.Context) ([]entities.Order, error) {
	query := `
		SELECT 
			id,
			user_id,
			restaurant_id,
			delivery_id,
			total_price,
			status,
			payment_status,
			address,
			comment,
			created_at
		FROM orders
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entities.Order

	for rows.Next() {
		var order entities.Order

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.RestaurantID,
			&order.DeliveryID,
			&order.TotalPrice,
			&order.Status,
			&order.PaymentStatus,
			&order.Address,
			&order.Comment,
			&order.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status string,
) error {
	query := `
		UPDATE orders
		SET status = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, status, id)

	return err
}

func (r *OrderRepository) UpdatePaymentStatus(
	ctx context.Context,
	id string,
	paymentStatus string,
) error {
	query := `
		UPDATE orders
		SET payment_status = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, paymentStatus, id)

	return err
}
