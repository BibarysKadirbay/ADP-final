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

type MenuRepository struct {
	db      *pgxpool.Pool
	metrics *metrics.Metrics
}

func NewMenuRepository(db *pgxpool.Pool, m *metrics.Metrics) *MenuRepository {
	return &MenuRepository{db: db, metrics: m}
}

func (r *MenuRepository) CreateItem(ctx context.Context, item *entities.MenuItem) error {
	now := time.Now().UTC()
	item.ID, item.CreatedAt, item.UpdatedAt = uuid.New(), now, now
	_, err := r.db.Exec(ctx, `insert into menu_items (id, category_id, restaurant_id, name, description, price, image_url, is_available, created_at, updated_at)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, item.ID, item.CategoryID, item.RestaurantID, item.Name, item.Description, item.Price, item.ImageURL, item.IsAvailable, item.CreatedAt, item.UpdatedAt)
	return err
}

func (r *MenuRepository) GetItemByID(ctx context.Context, id uuid.UUID) (*entities.MenuItem, error) {
	row := r.db.QueryRow(ctx, `select id, category_id, restaurant_id, name, description, price, image_url, is_available, created_at, updated_at from menu_items where id=$1`, id)
	return scanMenuItem(row)
}

func (r *MenuRepository) UpdateItem(ctx context.Context, item *entities.MenuItem) error {
	item.UpdatedAt = time.Now().UTC()
	tag, err := r.db.Exec(ctx, `update menu_items set category_id=$2, name=$3, description=$4, price=$5, image_url=$6, is_available=$7, updated_at=$8 where id=$1`,
		item.ID, item.CategoryID, item.Name, item.Description, item.Price, item.ImageURL, item.IsAvailable, item.UpdatedAt)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return services.ErrNotFound
	}
	return nil
}

func (r *MenuRepository) DeleteItem(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `delete from menu_items where id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return services.ErrNotFound
	}
	return nil
}

func (r *MenuRepository) SetAvailability(ctx context.Context, id uuid.UUID, available bool) (*entities.MenuItem, error) {
	row := r.db.QueryRow(ctx, `update menu_items set is_available=$2, updated_at=now() where id=$1 returning id, category_id, restaurant_id, name, description, price, image_url, is_available, created_at, updated_at`, id, available)
	return scanMenuItem(row)
}

func (r *MenuRepository) GetByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]entities.MenuCategory, error) {
	cats, err := NewCategoryRepository(r.db, r.metrics).ListByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, `select id, category_id, restaurant_id, name, description, price, image_url, is_available, created_at, updated_at from menu_items where restaurant_id=$1 order by name asc`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.MenuItem])
	if err != nil {
		return nil, err
	}
	grouped := make(map[uuid.UUID][]entities.MenuItem, len(cats))
	for _, item := range items {
		grouped[item.CategoryID] = append(grouped[item.CategoryID], item)
	}
	out := make([]entities.MenuCategory, 0, len(cats))
	for _, cat := range cats {
		out = append(out, entities.MenuCategory{Category: cat, Items: grouped[cat.ID]})
	}
	return out, nil
}

func scanMenuItem(row restaurantScanner) (*entities.MenuItem, error) {
	var item entities.MenuItem
	err := row.Scan(&item.ID, &item.CategoryID, &item.RestaurantID, &item.Name, &item.Description, &item.Price, &item.ImageURL, &item.IsAvailable, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, services.ErrNotFound
	}
	return &item, err
}
