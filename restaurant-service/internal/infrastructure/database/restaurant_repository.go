package database

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aitu/food-delivery/restaurant-service/internal/domain/entities"
	"github.com/aitu/food-delivery/restaurant-service/internal/domain/services"
	"github.com/aitu/food-delivery/restaurant-service/internal/infrastructure/metrics"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RestaurantRepository struct {
	db      *pgxpool.Pool
	metrics *metrics.Metrics
}

func NewRestaurantRepository(db *pgxpool.Pool, m *metrics.Metrics) *RestaurantRepository {
	return &RestaurantRepository{db: db, metrics: m}
}

func (r *RestaurantRepository) Create(ctx context.Context, e *entities.Restaurant) error {
	started := time.Now()
	defer r.metrics.ObserveDB("restaurant_create", started)
	now := time.Now().UTC()
	e.ID, e.CreatedAt, e.UpdatedAt = uuid.New(), now, now
	_, err := r.db.Exec(ctx, `insert into restaurants
		(id, owner_id, name, description, cuisine_type, address, city, rating, total_reviews, image_url, is_open, created_at, updated_at)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		e.ID, e.OwnerID, e.Name, e.Description, e.CuisineType, e.Address, e.City, e.Rating, e.TotalReviews, e.ImageURL, e.IsOpen, e.CreatedAt, e.UpdatedAt)
	return err
}

func (r *RestaurantRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Restaurant, error) {
	started := time.Now()
	defer r.metrics.ObserveDB("restaurant_get", started)
	row := r.db.QueryRow(ctx, `select id, owner_id, name, description, cuisine_type, address, city, rating, total_reviews, image_url, is_open, created_at, updated_at from restaurants where id=$1`, id)
	return scanRestaurant(row)
}

func (r *RestaurantRepository) Update(ctx context.Context, e *entities.Restaurant) error {
	started := time.Now()
	defer r.metrics.ObserveDB("restaurant_update", started)
	e.UpdatedAt = time.Now().UTC()
	tag, err := r.db.Exec(ctx, `update restaurants set name=$2, description=$3, cuisine_type=$4, address=$5, city=$6, image_url=$7, is_open=$8, updated_at=$9 where id=$1`,
		e.ID, e.Name, e.Description, e.CuisineType, e.Address, e.City, e.ImageURL, e.IsOpen, e.UpdatedAt)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return services.ErrNotFound
	}
	return nil
}

func (r *RestaurantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	started := time.Now()
	defer r.metrics.ObserveDB("restaurant_delete", started)
	tag, err := r.db.Exec(ctx, `delete from restaurants where id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return services.ErrNotFound
	}
	return nil
}

func (r *RestaurantRepository) List(ctx context.Context, f entities.RestaurantFilter) ([]entities.Restaurant, int64, error) {
	started := time.Now()
	defer r.metrics.ObserveDB("restaurant_list", started)
	where, args := restaurantWhere(f)
	sortBy := allow(f.SortBy, map[string]string{"name": "name", "rating": "rating", "created_at": "created_at"}, "created_at")
	sortDir := "desc"
	if strings.EqualFold(f.SortDirection, "asc") {
		sortDir = "asc"
	}
	countSQL := "select count(*) from restaurants " + where
	var total int64
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	args = append(args, f.Limit, f.Offset)
	rows, err := r.db.Query(ctx, fmt.Sprintf(`select id, owner_id, name, description, cuisine_type, address, city, rating, total_reviews, image_url, is_open, created_at, updated_at
		from restaurants %s order by %s %s limit $%d offset $%d`, where, sortBy, sortDir, len(args)-1, len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	list, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.Restaurant])
	return list, total, err
}

func (r *RestaurantRepository) TopRated(ctx context.Context, city string, limit int) ([]entities.Restaurant, error) {
	started := time.Now()
	defer r.metrics.ObserveDB("restaurant_top_rated", started)
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	args := []any{limit}
	where := "where is_open = true"
	if city != "" {
		where += " and city = $2"
		args = append(args, city)
	}
	rows, err := r.db.Query(ctx, `select id, owner_id, name, description, cuisine_type, address, city, rating, total_reviews, image_url, is_open, created_at, updated_at
		from restaurants `+where+` order by rating desc, total_reviews desc limit $1`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.Restaurant])
}

func restaurantWhere(f entities.RestaurantFilter) (string, []any) {
	clauses := make([]string, 0, 4)
	args := make([]any, 0, 4)
	add := func(clause string, value any) {
		args = append(args, value)
		clauses = append(clauses, fmt.Sprintf(clause, len(args)))
	}
	if f.Query != "" {
		args = append(args, f.Query)
		idx := len(args)
		clauses = append(clauses, fmt.Sprintf("(name ilike '%%' || $%d || '%%' or description ilike '%%' || $%d || '%%')", idx, idx))
	}
	if f.CuisineType != "" {
		add("cuisine_type = $%d", f.CuisineType)
	}
	if f.City != "" {
		add("city = $%d", f.City)
	}
	if f.OpenOnly != nil {
		add("is_open = $%d", *f.OpenOnly)
	}
	if len(clauses) == 0 {
		return "", args
	}
	return "where " + strings.Join(clauses, " and "), args
}

func allow(input string, allowed map[string]string, fallback string) string {
	if v, ok := allowed[input]; ok {
		return v
	}
	return fallback
}

type restaurantScanner interface {
	Scan(dest ...any) error
}

func scanRestaurant(row restaurantScanner) (*entities.Restaurant, error) {
	var e entities.Restaurant
	err := row.Scan(&e.ID, &e.OwnerID, &e.Name, &e.Description, &e.CuisineType, &e.Address, &e.City, &e.Rating, &e.TotalReviews, &e.ImageURL, &e.IsOpen, &e.CreatedAt, &e.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, services.ErrNotFound
	}
	return &e, err
}
