package database

import (
	"errors"
	"strconv"

	"github.com/aitu/food-delivery/delivery-service/internal/domain/services"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func itoa(v int) string {
	return strconv.Itoa(v)
}

func mapPgError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return services.ErrConflict
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return services.ErrNotFound
	}
	return err
}
