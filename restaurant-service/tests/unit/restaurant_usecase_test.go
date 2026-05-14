package unit

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aitu/food-delivery/restaurant-service/internal/domain/entities"
	"github.com/aitu/food-delivery/restaurant-service/internal/domain/services"
	"github.com/aitu/food-delivery/restaurant-service/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateRestaurantPublishesEventAndCachesInvalidated(t *testing.T) {
	ctx := context.Background()
	ownerID := uuid.New()
	repos := newFakes()
	uc := usecase.NewRestaurantUsecase(repos.restaurants, repos.categories, repos.menu, repos.cache, repos.events, time.Minute)

	created, err := uc.CreateRestaurant(ctx, &entities.Restaurant{
		OwnerID: ownerID, Name: "Sakura", CuisineType: "japanese", City: "Almaty", IsOpen: true,
	})

	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, created.ID)
	require.Contains(t, repos.events.subjects, usecase.EventRestaurantCreated)
	require.Contains(t, repos.cache.deletedPatterns, "restaurants:list:*")
}

func TestOwnerCannotUpdateAnotherOwnersRestaurant(t *testing.T) {
	ctx := context.Background()
	repos := newFakes()
	ownerID := uuid.New()
	otherOwnerID := uuid.New()
	existing := entities.Restaurant{ID: uuid.New(), OwnerID: ownerID, Name: "A", CuisineType: "korean", City: "Astana"}
	repos.restaurants.items[existing.ID] = existing
	uc := usecase.NewRestaurantUsecase(repos.restaurants, repos.categories, repos.menu, repos.cache, repos.events, time.Minute)

	_, err := uc.UpdateRestaurant(ctx, otherOwnerID, &entities.Restaurant{
		ID: existing.ID, OwnerID: otherOwnerID, Name: "B", CuisineType: "korean", City: "Astana",
	})

	require.ErrorIs(t, err, services.ErrForbidden)
}

func TestGetRestaurantUsesCacheOnSecondCall(t *testing.T) {
	ctx := context.Background()
	repos := newFakes()
	id := uuid.New()
	repos.restaurants.items[id] = entities.Restaurant{ID: id, OwnerID: uuid.New(), Name: "Cache Cafe", CuisineType: "local", City: "Atyrau"}
	uc := usecase.NewRestaurantUsecase(repos.restaurants, repos.categories, repos.menu, repos.cache, repos.events, time.Minute)

	_, err := uc.GetRestaurantByID(ctx, id)
	require.NoError(t, err)
	_, err = uc.GetRestaurantByID(ctx, id)
	require.NoError(t, err)

	require.Equal(t, 1, repos.restaurants.gets)
}

type fakes struct {
	restaurants *fakeRestaurantRepo
	categories  *fakeCategoryRepo
	menu        *fakeMenuRepo
	cache       *fakeCache
	events      *fakeEvents
}

func newFakes() fakes {
	return fakes{
		restaurants: &fakeRestaurantRepo{items: map[uuid.UUID]entities.Restaurant{}},
		categories:  &fakeCategoryRepo{items: map[uuid.UUID]entities.Category{}},
		menu:        &fakeMenuRepo{items: map[uuid.UUID]entities.MenuItem{}},
		cache:       &fakeCache{values: map[string][]byte{}},
		events:      &fakeEvents{},
	}
}

type fakeRestaurantRepo struct {
	mu    sync.Mutex
	items map[uuid.UUID]entities.Restaurant
	gets  int
}

func (r *fakeRestaurantRepo) Create(_ context.Context, e *entities.Restaurant) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	e.ID = uuid.New()
	r.items[e.ID] = *e
	return nil
}
func (r *fakeRestaurantRepo) GetByID(_ context.Context, id uuid.UUID) (*entities.Restaurant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.gets++
	item, ok := r.items[id]
	if !ok {
		return nil, services.ErrNotFound
	}
	return &item, nil
}
func (r *fakeRestaurantRepo) Update(_ context.Context, e *entities.Restaurant) error {
	r.items[e.ID] = *e
	return nil
}
func (r *fakeRestaurantRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.items, id)
	return nil
}
func (r *fakeRestaurantRepo) List(_ context.Context, f entities.RestaurantFilter) ([]entities.Restaurant, int64, error) {
	var out []entities.Restaurant
	for _, item := range r.items {
		if f.Query == "" || strings.Contains(strings.ToLower(item.Name), strings.ToLower(f.Query)) {
			out = append(out, item)
		}
	}
	return out, int64(len(out)), nil
}
func (r *fakeRestaurantRepo) TopRated(context.Context, string, int) ([]entities.Restaurant, error) {
	return nil, nil
}

type fakeCategoryRepo struct {
	items map[uuid.UUID]entities.Category
}

func (r *fakeCategoryRepo) Create(_ context.Context, c *entities.Category) error {
	c.ID = uuid.New()
	r.items[c.ID] = *c
	return nil
}
func (r *fakeCategoryRepo) GetByID(_ context.Context, id uuid.UUID) (*entities.Category, error) {
	c, ok := r.items[id]
	if !ok {
		return nil, services.ErrNotFound
	}
	return &c, nil
}
func (r *fakeCategoryRepo) Update(_ context.Context, c *entities.Category) error {
	r.items[c.ID] = *c
	return nil
}
func (r *fakeCategoryRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.items, id)
	return nil
}
func (r *fakeCategoryRepo) ListByRestaurant(_ context.Context, restaurantID uuid.UUID) ([]entities.Category, error) {
	var out []entities.Category
	for _, c := range r.items {
		if c.RestaurantID == restaurantID {
			out = append(out, c)
		}
	}
	return out, nil
}

type fakeMenuRepo struct {
	items map[uuid.UUID]entities.MenuItem
}

func (r *fakeMenuRepo) CreateItem(_ context.Context, item *entities.MenuItem) error {
	item.ID = uuid.New()
	r.items[item.ID] = *item
	return nil
}
func (r *fakeMenuRepo) GetItemByID(_ context.Context, id uuid.UUID) (*entities.MenuItem, error) {
	item, ok := r.items[id]
	if !ok {
		return nil, services.ErrNotFound
	}
	return &item, nil
}
func (r *fakeMenuRepo) UpdateItem(_ context.Context, item *entities.MenuItem) error {
	r.items[item.ID] = *item
	return nil
}
func (r *fakeMenuRepo) DeleteItem(_ context.Context, id uuid.UUID) error {
	delete(r.items, id)
	return nil
}
func (r *fakeMenuRepo) GetByRestaurant(_ context.Context, restaurantID uuid.UUID) ([]entities.MenuCategory, error) {
	return []entities.MenuCategory{}, nil
}
func (r *fakeMenuRepo) SetAvailability(_ context.Context, id uuid.UUID, available bool) (*entities.MenuItem, error) {
	item := r.items[id]
	item.IsAvailable = available
	r.items[id] = item
	return &item, nil
}

type fakeCache struct {
	values          map[string][]byte
	deletedPatterns []string
}

func (c *fakeCache) Get(_ context.Context, key string, dest any) (bool, error) {
	raw, ok := c.values[key]
	if !ok {
		return false, nil
	}
	return true, json.Unmarshal(raw, dest)
}
func (c *fakeCache) Set(_ context.Context, key string, value any, _ time.Duration) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	c.values[key] = raw
	return nil
}
func (c *fakeCache) Delete(_ context.Context, keys ...string) error {
	for _, key := range keys {
		delete(c.values, key)
	}
	return nil
}
func (c *fakeCache) DeletePattern(_ context.Context, pattern string) error {
	c.deletedPatterns = append(c.deletedPatterns, pattern)
	return nil
}

type fakeEvents struct{ subjects []string }

func (e *fakeEvents) Publish(_ context.Context, subject string, _ any) error {
	e.subjects = append(e.subjects, subject)
	return nil
}
func (e *fakeEvents) Close() error { return nil }
