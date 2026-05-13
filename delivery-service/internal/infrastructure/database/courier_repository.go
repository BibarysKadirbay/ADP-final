package database

import (
	"context"
	"errors"
	"time"

	"github.com/aitu/food-delivery/delivery-service/internal/domain/entities"
	"github.com/aitu/food-delivery/delivery-service/internal/domain/services"
	"github.com/aitu/food-delivery/delivery-service/internal/infrastructure/metrics"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CourierRepository struct {
	db      *pgxpool.Pool
	metrics *metrics.Metrics
}

func NewCourierRepository(db *pgxpool.Pool, m *metrics.Metrics) *CourierRepository {
	return &CourierRepository{db: db, metrics: m}
}

func (r *CourierRepository) Create(ctx context.Context, c *entities.Courier) error {
	started := time.Now()
	defer r.metrics.ObserveDB("courier_create", started)
	now := time.Now().UTC()
	c.ID, c.CreatedAt, c.UpdatedAt = uuid.New(), now, now
	c.Rating = 5
	_, err := r.db.Exec(ctx, `insert into couriers
		(id, user_id, full_name, phone, vehicle_type, rating, total_deliveries, is_available, created_at, updated_at)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		c.ID, c.UserID, c.FullName, c.Phone, c.VehicleType, c.Rating, c.TotalDeliveries, c.IsAvailable, c.CreatedAt, c.UpdatedAt)
	return mapPgError(err)
}

func (r *CourierRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Courier, error) {
	started := time.Now()
	defer r.metrics.ObserveDB("courier_get", started)
	row := r.db.QueryRow(ctx, `select id, user_id, full_name, phone, vehicle_type, rating, total_deliveries, is_available, created_at, updated_at from couriers where id=$1`, id)
	return scanCourier(row)
}

func (r *CourierRepository) UpdateAvailability(ctx context.Context, id uuid.UUID, available bool) (*entities.Courier, error) {
	started := time.Now()
	defer r.metrics.ObserveDB("courier_availability", started)
	row := r.db.QueryRow(ctx, `update couriers set is_available=$2, updated_at=now() where id=$1
		returning id, user_id, full_name, phone, vehicle_type, rating, total_deliveries, is_available, created_at, updated_at`, id, available)
	return scanCourier(row)
}

func (r *CourierRepository) ListAvailable(ctx context.Context, f entities.CourierFilter) ([]entities.Courier, int64, error) {
	started := time.Now()
	defer r.metrics.ObserveDB("courier_list_available", started)
	where := "where is_available = true"
	args := []any{}
	if f.VehicleType != "" {
		args = append(args, f.VehicleType)
		where += " and vehicle_type = $1"
	}
	var total int64
	if err := r.db.QueryRow(ctx, "select count(*) from couriers "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	args = append(args, f.Limit, f.Offset)
	rows, err := r.db.Query(ctx, `select id, user_id, full_name, phone, vehicle_type, rating, total_deliveries, is_available, created_at, updated_at
		from couriers `+where+` order by rating desc, total_deliveries asc limit $`+itoa(len(args)-1)+` offset $`+itoa(len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.Courier])
	return items, total, err
}

func (r *CourierRepository) RecalculateRating(ctx context.Context, courierID uuid.UUID) error {
	started := time.Now()
	defer r.metrics.ObserveDB("courier_rating_recalculate", started)
	tag, err := r.db.Exec(ctx, `update couriers set rating = coalesce((select avg(rating)::numeric(3,2) from courier_ratings where courier_id=$1), rating), updated_at=now() where id=$1`, courierID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return services.ErrNotFound
	}
	return nil
}

type courierScanner interface {
	Scan(dest ...any) error
}

func scanCourier(row courierScanner) (*entities.Courier, error) {
	var c entities.Courier
	err := row.Scan(&c.ID, &c.UserID, &c.FullName, &c.Phone, &c.VehicleType, &c.Rating, &c.TotalDeliveries, &c.IsAvailable, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, services.ErrNotFound
	}
	return &c, err
}
