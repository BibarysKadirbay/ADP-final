package grpc

import (
	"github.com/aitu/food-delivery/restaurant-service/internal/domain/entities"
	"github.com/aitu/food-delivery/restaurant-service/internal/infrastructure/grpc/restaurantpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toPBRestaurant(r *entities.Restaurant) *restaurantpb.Restaurant {
	if r == nil {
		return nil
	}
	return &restaurantpb.Restaurant{
		Id: r.ID.String(), OwnerId: r.OwnerID.String(), Name: r.Name, Description: r.Description,
		CuisineType: r.CuisineType, Address: r.Address, City: r.City, Rating: r.Rating,
		TotalReviews: r.TotalReviews, ImageUrl: r.ImageURL, IsOpen: r.IsOpen,
		CreatedAt: timestamppb.New(r.CreatedAt), UpdatedAt: timestamppb.New(r.UpdatedAt),
	}
}

func toPBRestaurants(items []entities.Restaurant) []*restaurantpb.Restaurant {
	out := make([]*restaurantpb.Restaurant, 0, len(items))
	for i := range items {
		out = append(out, toPBRestaurant(&items[i]))
	}
	return out
}

func toPBCategory(c *entities.Category) *restaurantpb.Category {
	if c == nil {
		return nil
	}
	return &restaurantpb.Category{Id: c.ID.String(), RestaurantId: c.RestaurantID.String(), Name: c.Name, CreatedAt: timestamppb.New(c.CreatedAt)}
}

func toPBCategories(items []entities.Category) []*restaurantpb.Category {
	out := make([]*restaurantpb.Category, 0, len(items))
	for i := range items {
		out = append(out, toPBCategory(&items[i]))
	}
	return out
}

func toPBMenuItem(item *entities.MenuItem) *restaurantpb.MenuItem {
	if item == nil {
		return nil
	}
	return &restaurantpb.MenuItem{
		Id: item.ID.String(), CategoryId: item.CategoryID.String(), RestaurantId: item.RestaurantID.String(),
		Name: item.Name, Description: item.Description, Price: item.Price, ImageUrl: item.ImageURL,
		IsAvailable: item.IsAvailable, CreatedAt: timestamppb.New(item.CreatedAt), UpdatedAt: timestamppb.New(item.UpdatedAt),
	}
}

func toPBMenu(menu []entities.MenuCategory) []*restaurantpb.MenuCategory {
	out := make([]*restaurantpb.MenuCategory, 0, len(menu))
	for _, cat := range menu {
		items := make([]*restaurantpb.MenuItem, 0, len(cat.Items))
		for i := range cat.Items {
			items = append(items, toPBMenuItem(&cat.Items[i]))
		}
		out = append(out, &restaurantpb.MenuCategory{Category: toPBCategory(&cat.Category), Items: items})
	}
	return out
}
