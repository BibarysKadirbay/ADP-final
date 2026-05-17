package database

import (
	"context"
	"errors"

	"github.com/aitu/food-delivery/payment-service/internal/domain/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreateInTx(ctx context.Context, payment *entities.Payment) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO payments(id, order_id, user_id, amount, status, method, created_at)
		VALUES($1,$2,$3,$4,$5,$6,$7)
	`
	if _, err := tx.Exec(ctx, query,
		payment.ID, payment.OrderID, payment.UserID, payment.Amount,
		payment.Status, payment.Method, payment.CreatedAt,
	); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *PaymentRepository) GetByID(ctx context.Context, id string) (*entities.Payment, error) {
	query := `SELECT id, order_id, user_id, amount, status, method, created_at FROM payments WHERE id = $1`
	var p entities.Payment
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.OrderID, &p.UserID, &p.Amount, &p.Status, &p.Method, &p.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("payment not found")
	}
	return &p, err
}

func (r *PaymentRepository) ListByOrder(ctx context.Context, orderID string) ([]entities.Payment, error) {
	query := `SELECT id, order_id, user_id, amount, status, method, created_at FROM payments WHERE order_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entities.Payment
	for rows.Next() {
		var p entities.Payment
		if err := rows.Scan(&p.ID, &p.OrderID, &p.UserID, &p.Amount, &p.Status, &p.Method, &p.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, rows.Err()
}
