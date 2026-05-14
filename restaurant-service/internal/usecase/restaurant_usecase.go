package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aitu/food-delivery/restaurant-service/internal/domain/entities"
	"github.com/aitu/food-delivery/restaurant-service/internal/domain/repositories"
	"github.com/aitu/food-delivery/restaurant-service/internal/domain/services"
	"github.com/google/uuid"
)

const (
	EventRestaurantCreated = "restaurant.created"
	EventRestaurantUpdated = "restaurant.updated"
	EventRestaurantDeleted = "restaurant.deleted"
	EventMenuItemCreated   = "menu.item.created"
	EventMenuItemUpdated   = "menu.item.updated"
	EventMenuItemDeleted   = "menu.item.deleted"
	EventMenuAvailability  = "menu.item.availability_changed"
)

type RestaurantUsecase struct {
	restaurants repositories.RestaurantRepository
	categories  repositories.CategoryRepository
	menu        repositories.MenuRepository
	cache       repositories.CacheRepository
	events      repositories.EventPublisher
	cacheTTL    time.Duration
}

func NewRestaurantUsecase(
	restaurants repositories.RestaurantRepository,
	categories repositories.CategoryRepository,
	menu repositories.MenuRepository,
	cache repositories.CacheRepository,
	events repositories.EventPublisher,
	cacheTTL time.Duration,
) *RestaurantUsecase {
	return &RestaurantUsecase{restaurants: restaurants, categories: categories, menu: menu, cache: cache, events: events, cacheTTL: cacheTTL}
}

func (u *RestaurantUsecase) CreateRestaurant(ctx context.Context, r *entities.Restaurant) (*entities.Restaurant, error) {
	if err := validateRestaurant(r); err != nil {
		return nil, err
	}
	r.Name, r.City, r.CuisineType = strings.TrimSpace(r.Name), strings.TrimSpace(r.City), strings.TrimSpace(r.CuisineType)
	if err := u.restaurants.Create(ctx, r); err != nil {
		return nil, err
	}
	u.invalidateRestaurantLists(ctx)
	_ = u.events.Publish(ctx, EventRestaurantCreated, r)
	return r, nil
}

func (u *RestaurantUsecase) GetRestaurantByID(ctx context.Context, id uuid.UUID) (*entities.Restaurant, error) {
	key := restaurantKey(id)
	var cached entities.Restaurant
	if ok, err := u.cache.Get(ctx, key, &cached); err == nil && ok {
		return &cached, nil
	}
	r, err := u.restaurants.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = u.cache.Set(ctx, key, r, u.cacheTTL)
	return r, nil
}

func (u *RestaurantUsecase) UpdateRestaurant(ctx context.Context, ownerID uuid.UUID, updated *entities.Restaurant) (*entities.Restaurant, error) {
	if err := validateRestaurant(updated); err != nil {
		return nil, err
	}
	current, err := u.restaurants.GetByID(ctx, updated.ID)
	if err != nil {
		return nil, err
	}
	if current.OwnerID != ownerID {
		return nil, services.ErrForbidden
	}
	updated.OwnerID = current.OwnerID
	updated.Rating = current.Rating
	updated.TotalReviews = current.TotalReviews
	updated.CreatedAt = current.CreatedAt
	if err := u.restaurants.Update(ctx, updated); err != nil {
		return nil, err
	}
	u.invalidateRestaurant(ctx, updated.ID)
	_ = u.events.Publish(ctx, EventRestaurantUpdated, updated)
	return updated, nil
}

func (u *RestaurantUsecase) DeleteRestaurant(ctx context.Context, id, ownerID uuid.UUID) error {
	current, err := u.restaurants.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if current.OwnerID != ownerID {
		return services.ErrForbidden
	}
	if err := u.restaurants.Delete(ctx, id); err != nil {
		return err
	}
	u.invalidateRestaurant(ctx, id)
	_ = u.events.Publish(ctx, EventRestaurantDeleted, map[string]string{"id": id.String()})
	return nil
}

func (u *RestaurantUsecase) ListRestaurants(ctx context.Context, filter entities.RestaurantFilter, page, pageSize int) ([]entities.Restaurant, entities.PageMeta, error) {
	page, pageSize = normalizePage(page, pageSize)
	filter.Limit, filter.Offset = pageSize, (page-1)*pageSize
	key := fmt.Sprintf("restaurants:list:%s:%s:%s:%v:%s:%s:%d:%d", filter.Query, filter.CuisineType, filter.City, filter.OpenOnly, filter.SortBy, filter.SortDirection, page, pageSize)
	var cached struct {
		Items []entities.Restaurant `json:"items"`
		Meta  entities.PageMeta     `json:"meta"`
	}
	if ok, err := u.cache.Get(ctx, key, &cached); err == nil && ok {
		return cached.Items, cached.Meta, nil
	}
	items, total, err := u.restaurants.List(ctx, filter)
	if err != nil {
		return nil, entities.PageMeta{}, err
	}
	meta := entities.PageMeta{Page: page, PageSize: pageSize, Total: total}
	_ = u.cache.Set(ctx, key, struct {
		Items []entities.Restaurant `json:"items"`
		Meta  entities.PageMeta     `json:"meta"`
	}{items, meta}, u.cacheTTL)
	return items, meta, nil
}

func (u *RestaurantUsecase) SearchRestaurants(ctx context.Context, query string, filter entities.RestaurantFilter, page, pageSize int) ([]entities.Restaurant, entities.PageMeta, error) {
	filter.Query = query
	return u.ListRestaurants(ctx, filter, page, pageSize)
}

func (u *RestaurantUsecase) TopRated(ctx context.Context, city string, limit int) ([]entities.Restaurant, error) {
	key := fmt.Sprintf("restaurants:top:%s:%d", city, limit)
	var cached []entities.Restaurant
	if ok, err := u.cache.Get(ctx, key, &cached); err == nil && ok {
		return cached, nil
	}
	items, err := u.restaurants.TopRated(ctx, city, limit)
	if err != nil {
		return nil, err
	}
	_ = u.cache.Set(ctx, key, items, u.cacheTTL)
	return items, nil
}

func (u *RestaurantUsecase) CreateCategory(ctx context.Context, ownerID uuid.UUID, c *entities.Category) (*entities.Category, error) {
	if strings.TrimSpace(c.Name) == "" {
		return nil, services.ErrInvalidInput
	}
	if err := u.ensureOwner(ctx, c.RestaurantID, ownerID); err != nil {
		return nil, err
	}
	if err := u.categories.Create(ctx, c); err != nil {
		return nil, err
	}
	u.invalidateMenu(ctx, c.RestaurantID)
	return c, nil
}

func (u *RestaurantUsecase) UpdateCategory(ctx context.Context, ownerID uuid.UUID, c *entities.Category) (*entities.Category, error) {
	current, err := u.categories.GetByID(ctx, c.ID)
	if err != nil {
		return nil, err
	}
	if err := u.ensureOwner(ctx, current.RestaurantID, ownerID); err != nil {
		return nil, err
	}
	current.Name = strings.TrimSpace(c.Name)
	if current.Name == "" {
		return nil, services.ErrInvalidInput
	}
	if err := u.categories.Update(ctx, current); err != nil {
		return nil, err
	}
	u.invalidateMenu(ctx, current.RestaurantID)
	return current, nil
}

func (u *RestaurantUsecase) DeleteCategory(ctx context.Context, id, ownerID uuid.UUID) error {
	current, err := u.categories.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := u.ensureOwner(ctx, current.RestaurantID, ownerID); err != nil {
		return err
	}
	if err := u.categories.Delete(ctx, id); err != nil {
		return err
	}
	u.invalidateMenu(ctx, current.RestaurantID)
	return nil
}

func (u *RestaurantUsecase) ListCategories(ctx context.Context, restaurantID uuid.UUID) ([]entities.Category, error) {
	return u.categories.ListByRestaurant(ctx, restaurantID)
}

func (u *RestaurantUsecase) CreateMenuItem(ctx context.Context, ownerID uuid.UUID, item *entities.MenuItem) (*entities.MenuItem, error) {
	if err := validateMenuItem(item); err != nil {
		return nil, err
	}
	if err := u.ensureOwner(ctx, item.RestaurantID, ownerID); err != nil {
		return nil, err
	}
	if err := u.menu.CreateItem(ctx, item); err != nil {
		return nil, err
	}
	u.invalidateMenu(ctx, item.RestaurantID)
	_ = u.events.Publish(ctx, EventMenuItemCreated, item)
	return item, nil
}

func (u *RestaurantUsecase) UpdateMenuItem(ctx context.Context, ownerID uuid.UUID, item *entities.MenuItem) (*entities.MenuItem, error) {
	if err := validateMenuItem(item); err != nil {
		return nil, err
	}
	current, err := u.menu.GetItemByID(ctx, item.ID)
	if err != nil {
		return nil, err
	}
	if err := u.ensureOwner(ctx, current.RestaurantID, ownerID); err != nil {
		return nil, err
	}
	item.RestaurantID = current.RestaurantID
	item.CreatedAt = current.CreatedAt
	if err := u.menu.UpdateItem(ctx, item); err != nil {
		return nil, err
	}
	u.invalidateMenu(ctx, item.RestaurantID)
	_ = u.events.Publish(ctx, EventMenuItemUpdated, item)
	return item, nil
}

func (u *RestaurantUsecase) DeleteMenuItem(ctx context.Context, id, ownerID uuid.UUID) error {
	current, err := u.menu.GetItemByID(ctx, id)
	if err != nil {
		return err
	}
	if err := u.ensureOwner(ctx, current.RestaurantID, ownerID); err != nil {
		return err
	}
	if err := u.menu.DeleteItem(ctx, id); err != nil {
		return err
	}
	u.invalidateMenu(ctx, current.RestaurantID)
	_ = u.events.Publish(ctx, EventMenuItemDeleted, map[string]string{"id": id.String(), "restaurant_id": current.RestaurantID.String()})
	return nil
}

func (u *RestaurantUsecase) GetMenuByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]entities.MenuCategory, error) {
	key := menuKey(restaurantID)
	var cached []entities.MenuCategory
	if ok, err := u.cache.Get(ctx, key, &cached); err == nil && ok {
		return cached, nil
	}
	menu, err := u.menu.GetByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, err
	}
	_ = u.cache.Set(ctx, key, menu, u.cacheTTL)
	return menu, nil
}

func (u *RestaurantUsecase) SetMenuItemAvailability(ctx context.Context, id, ownerID uuid.UUID, available bool) (*entities.MenuItem, error) {
	current, err := u.menu.GetItemByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := u.ensureOwner(ctx, current.RestaurantID, ownerID); err != nil {
		return nil, err
	}
	item, err := u.menu.SetAvailability(ctx, id, available)
	if err != nil {
		return nil, err
	}
	u.invalidateMenu(ctx, item.RestaurantID)
	_ = u.events.Publish(ctx, EventMenuAvailability, item)
	return item, nil
}

func (u *RestaurantUsecase) ensureOwner(ctx context.Context, restaurantID, ownerID uuid.UUID) error {
	r, err := u.restaurants.GetByID(ctx, restaurantID)
	if err != nil {
		return err
	}
	if r.OwnerID != ownerID {
		return services.ErrForbidden
	}
	return nil
}

func (u *RestaurantUsecase) invalidateRestaurant(ctx context.Context, id uuid.UUID) {
	_ = u.cache.Delete(ctx, restaurantKey(id))
	u.invalidateRestaurantLists(ctx)
	u.invalidateMenu(ctx, id)
}

func (u *RestaurantUsecase) invalidateRestaurantLists(ctx context.Context) {
	_ = u.cache.DeletePattern(ctx, "restaurants:list:*")
	_ = u.cache.DeletePattern(ctx, "restaurants:top:*")
}

func (u *RestaurantUsecase) invalidateMenu(ctx context.Context, restaurantID uuid.UUID) {
	_ = u.cache.Delete(ctx, menuKey(restaurantID))
}

func restaurantKey(id uuid.UUID) string { return "restaurants:detail:" + id.String() }
func menuKey(id uuid.UUID) string       { return "restaurants:menu:" + id.String() }

func normalizePage(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return page, size
}

func validateRestaurant(r *entities.Restaurant) error {
	if r == nil || r.OwnerID == uuid.Nil || strings.TrimSpace(r.Name) == "" || strings.TrimSpace(r.City) == "" || strings.TrimSpace(r.CuisineType) == "" {
		return services.ErrInvalidInput
	}
	return nil
}

func validateMenuItem(item *entities.MenuItem) error {
	if item == nil || item.RestaurantID == uuid.Nil || item.CategoryID == uuid.Nil || strings.TrimSpace(item.Name) == "" || item.Price < 0 {
		return services.ErrInvalidInput
	}
	return nil
}
