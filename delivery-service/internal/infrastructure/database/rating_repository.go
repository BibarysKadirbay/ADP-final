package database

import (
	"context"
	"time"

	"github.com/aitu/food-delivery/delivery-service/internal/domain/entities"
	"github.com/aitu/food-delivery/delivery-service/internal/infrastructure/metrics"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RatingRepository struct {
	db      *pgxpool.Pool
	metrics *metrics.Metrics
}

func NewRatingRepository(db *pgxpool.Pool, m *metrics.Metrics) *RatingRepository {
	return &RatingRepository{db: db, metrics: m}
}

func (r *RatingRepository) Create(ctx context.Context, rating *entities.CourierRating) error {
	started := time.Now()
	defer r.metrics.ObserveDB("courier_rating_create", started)
	rating.ID, rating.CreatedAt = uuid.New(), time.Now().UTC()
	_, err := r.db.Exec(ctx, `insert into courier_ratings (id, courier_id, order_id, customer_id, rating, comment, created_at)
		values ($1,$2,$3,$4,$5,$6,$7)`, rating.ID, rating.CourierID, rating.OrderID, rating.CustomerID, rating.Rating, rating.Comment, rating.CreatedAt)
	return mapPgError(err)
}
