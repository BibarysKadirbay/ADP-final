package database

import (
	"context"
	"errors"
	"time"

	"github.com/aitu/food-delivery/restaurant-service/internal/domain/entities"
	"github.com/aitu/food-delivery/restaurant-service/internal/domain/services"
	"github.com/aitu/food-delivery/restaurant-service/internal/infrastructure/metrics"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	db      *pgxpool.Pool
	metrics *metrics.Metrics
}

func NewCategoryRepository(db *pgxpool.Pool, m *metrics.Metrics) *CategoryRepository {
	return &CategoryRepository{db: db, metrics: m}
}

func (r *CategoryRepository) Create(ctx context.Context, c *entities.Category) error {
	c.ID, c.CreatedAt = uuid.New(), time.Now().UTC()
	_, err := r.db.Exec(ctx, `insert into menu_categories (id, restaurant_id, name, created_at) values ($1,$2,$3,$4)`, c.ID, c.RestaurantID, c.Name, c.CreatedAt)
	return err
}

func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Category, error) {
	row := r.db.QueryRow(ctx, `select id, restaurant_id, name, created_at from menu_categories where id=$1`, id)
	var c entities.Category
	err := row.Scan(&c.ID, &c.RestaurantID, &c.Name, &c.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, services.ErrNotFound
	}
	return &c, err
}

func (r *CategoryRepository) Update(ctx context.Context, c *entities.Category) error {
	tag, err := r.db.Exec(ctx, `update menu_categories set name=$2 where id=$1`, c.ID, c.Name)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return services.ErrNotFound
	}
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `delete from menu_categories where id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return services.ErrNotFound
	}
	return nil
}

func (r *CategoryRepository) ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]entities.Category, error) {
	rows, err := r.db.Query(ctx, `select id, restaurant_id, name, created_at from menu_categories where restaurant_id=$1 order by created_at asc`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.Category])
}
