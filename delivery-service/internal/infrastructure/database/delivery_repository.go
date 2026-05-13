package database

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aitu/food-delivery/delivery-service/internal/domain/entities"
	"github.com/aitu/food-delivery/delivery-service/internal/domain/services"
	"github.com/aitu/food-delivery/delivery-service/internal/infrastructure/metrics"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeliveryRepository struct {
	db      *pgxpool.Pool
	metrics *metrics.Metrics
}

func NewDeliveryRepository(db *pgxpool.Pool, m *metrics.Metrics) *DeliveryRepository {
	return &DeliveryRepository{db: db, metrics: m}
}

func (r *DeliveryRepository) Assign(ctx context.Context, d *entities.Delivery) error {
	started := time.Now()
	defer r.metrics.ObserveDB("delivery_assign", started)
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	now := time.Now().UTC()
	d.ID, d.CreatedAt, d.UpdatedAt = uuid.New(), now, now
	_, err = tx.Exec(ctx, `insert into deliveries
		(id, order_id, courier_id, restaurant_id, customer_id, status, pickup_address, delivery_address, estimated_eta_minutes, route_distance_km, created_at, updated_at)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		d.ID, d.OrderID, d.CourierID, d.RestaurantID, d.CustomerID, d.Status, d.PickupAddress, d.DeliveryAddress, d.EstimatedETAMinutes, d.RouteDistanceKM, d.CreatedAt, d.UpdatedAt)
	if err != nil {
		return mapPgError(err)
	}
	_, err = tx.Exec(ctx, `insert into delivery_status_history (id, delivery_id, old_status, new_status, changed_at) values ($1,$2,$3,$4,$5)`,
		uuid.New(), d.ID, nil, d.Status, now)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `update couriers set is_available=false, updated_at=now() where id=$1`, d.CourierID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *DeliveryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Delivery, error) {
	started := time.Now()
	defer r.metrics.ObserveDB("delivery_get", started)
	row := r.db.QueryRow(ctx, baseDeliverySelect()+" where id=$1", id)
	return scanDelivery(row)
}

func (r *DeliveryRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*entities.Delivery, error) {
	row := r.db.QueryRow(ctx, baseDeliverySelect()+" where order_id=$1 order by created_at desc limit 1", orderID)
	return scanDelivery(row)
}

func (r *DeliveryRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entities.DeliveryStatus) (*entities.Delivery, error) {
	started := time.Now()
	defer r.metrics.ObserveDB("delivery_update_status", started)
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	current, err := scanDelivery(tx.QueryRow(ctx, baseDeliverySelect()+" where id=$1 for update", id))
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	pickupSet := current.PickupTime
	deliveredSet := current.DeliveredTime
	if status == entities.StatusPickedUp && pickupSet == nil {
		pickupSet = &now
	}
	if status == entities.StatusDelivered && deliveredSet == nil {
		deliveredSet = &now
	}
	updated, err := scanDelivery(tx.QueryRow(ctx, `update deliveries set status=$2, pickup_time=$3, delivered_time=$4, updated_at=$5 where id=$1
		returning id, order_id, courier_id, restaurant_id, customer_id, status, pickup_address, delivery_address, estimated_eta_minutes, pickup_time, delivered_time, route_distance_km, created_at, updated_at`,
		id, status, pickupSet, deliveredSet, now))
	if err != nil {
		return nil, err
	}
	oldStatus := current.Status
	_, err = tx.Exec(ctx, `insert into delivery_status_history (id, delivery_id, old_status, new_status, changed_at) values ($1,$2,$3,$4,$5)`,
		uuid.New(), id, oldStatus, status, now)
	if err != nil {
		return nil, err
	}
	if services.IsTerminal(status) {
		_, err = tx.Exec(ctx, `update couriers set total_deliveries = total_deliveries + case when $2 then 1 else 0 end, is_available=true, updated_at=now() where id=$1`,
			current.CourierID, status == entities.StatusDelivered)
		if err != nil {
			return nil, err
		}
	}
	return updated, tx.Commit(ctx)
}

func (r *DeliveryRepository) ListByCourier(ctx context.Context, courierID uuid.UUID, f entities.DeliveryFilter) ([]entities.Delivery, int64, error) {
	started := time.Now()
	defer r.metrics.ObserveDB("delivery_list_courier", started)
	where := "where courier_id=$1"
	args := []any{courierID}
	if f.Status != "" {
		args = append(args, f.Status)
		where += fmt.Sprintf(" and status=$%d", len(args))
	}
	var total int64
	if err := r.db.QueryRow(ctx, "select count(*) from deliveries "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	sortBy := allow(f.SortBy, map[string]string{"created_at": "created_at", "status": "status", "updated_at": "updated_at"}, "created_at")
	sortDir := "desc"
	if strings.EqualFold(f.SortDirection, "asc") {
		sortDir = "asc"
	}
	args = append(args, f.Limit, f.Offset)
	rows, err := r.db.Query(ctx, fmt.Sprintf(`%s %s order by %s %s limit $%d offset $%d`, baseDeliverySelect(), where, sortBy, sortDir, len(args)-1, len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.Delivery])
	return items, total, err
}

func (r *DeliveryRepository) ListByOrder(ctx context.Context, orderID uuid.UUID) ([]entities.Delivery, error) {
	rows, err := r.db.Query(ctx, baseDeliverySelect()+" where order_id=$1 order by created_at desc", orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.Delivery])
}

func (r *DeliveryRepository) History(ctx context.Context, deliveryID uuid.UUID) ([]entities.DeliveryStatusHistory, error) {
	rows, err := r.db.Query(ctx, `select id, delivery_id, old_status, new_status, changed_at from delivery_status_history where delivery_id=$1 order by changed_at asc`, deliveryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.DeliveryStatusHistory])
}

func (r *DeliveryRepository) Stats(ctx context.Context, courierID uuid.UUID) (*entities.DeliveryStats, error) {
	var s entities.DeliveryStats
	err := r.db.QueryRow(ctx, `select
		count(*),
		count(*) filter (where status in ('assigned','picked_up','on_the_way')),
		count(*) filter (where status='delivered'),
		count(*) filter (where status='cancelled'),
		coalesce(avg(estimated_eta_minutes),0),
		coalesce(avg(route_distance_km),0)
		from deliveries where courier_id=$1`, courierID).
		Scan(&s.TotalDeliveries, &s.ActiveDeliveries, &s.CompletedDeliveries, &s.CancelledDeliveries, &s.AverageETAMinutes, &s.AverageDistanceKM)
	return &s, err
}

func (r *DeliveryRepository) ActiveCountByCourier(ctx context.Context, courierID uuid.UUID) (int64, error) {
	var total int64
	err := r.db.QueryRow(ctx, `select count(*) from deliveries where courier_id=$1 and status in ('assigned','picked_up','on_the_way')`, courierID).Scan(&total)
	return total, err
}

func (r *DeliveryRepository) CancelByOrder(ctx context.Context, orderID uuid.UUID) (*entities.Delivery, error) {
	d, err := r.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if services.IsTerminal(d.Status) {
		return d, nil
	}
	return r.UpdateStatus(ctx, d.ID, entities.StatusCancelled)
}

type deliveryScanner interface {
	Scan(dest ...any) error
}

func scanDelivery(row deliveryScanner) (*entities.Delivery, error) {
	var d entities.Delivery
	err := row.Scan(&d.ID, &d.OrderID, &d.CourierID, &d.RestaurantID, &d.CustomerID, &d.Status, &d.PickupAddress, &d.DeliveryAddress, &d.EstimatedETAMinutes, &d.PickupTime, &d.DeliveredTime, &d.RouteDistanceKM, &d.CreatedAt, &d.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, services.ErrNotFound
	}
	return &d, err
}

func baseDeliverySelect() string {
	return `select id, order_id, courier_id, restaurant_id, customer_id, status, pickup_address, delivery_address, estimated_eta_minutes, pickup_time, delivered_time, route_distance_km, created_at, updated_at from deliveries`
}

func allow(input string, allowed map[string]string, fallback string) string {
	if v, ok := allowed[input]; ok {
		return v
	}
	return fallback
}
