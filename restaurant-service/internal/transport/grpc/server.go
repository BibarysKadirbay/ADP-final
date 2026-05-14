package grpc

import (
	"context"

	"github.com/aitu/food-delivery/restaurant-service/internal/domain/entities"
	"github.com/aitu/food-delivery/restaurant-service/internal/domain/services"
	"github.com/aitu/food-delivery/restaurant-service/internal/infrastructure/grpc/restaurantpb"
	"github.com/aitu/food-delivery/restaurant-service/internal/usecase"
	"github.com/google/uuid"
)

type Server struct {
	restaurantpb.UnimplementedRestaurantServiceServer
	uc          *usecase.RestaurantUsecase
	serviceName string
}

func NewServer(uc *usecase.RestaurantUsecase, serviceName string) *Server {
	return &Server{uc: uc, serviceName: serviceName}
}

func (s *Server) CreateRestaurant(ctx context.Context, req *restaurantpb.CreateRestaurantRequest) (*restaurantpb.RestaurantResponse, error) {
	ownerID, err := uuid.Parse(req.OwnerId)
	if err != nil {
		return nil, toStatus(services.ErrInvalidInput)
	}
	r, err := s.uc.CreateRestaurant(ctx, &entities.Restaurant{
		OwnerID: ownerID, Name: req.Name, Description: req.Description, CuisineType: req.CuisineType,
		Address: req.Address, City: req.City, ImageURL: req.ImageUrl, IsOpen: req.IsOpen,
	})
	if err != nil {
		return nil, toStatus(err)
	}
	return &restaurantpb.RestaurantResponse{Restaurant: toPBRestaurant(r)}, nil
}

func (s *Server) GetRestaurantById(ctx context.Context, req *restaurantpb.GetRestaurantByIdRequest) (*restaurantpb.RestaurantResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, toStatus(services.ErrInvalidInput)
	}
	r, err := s.uc.GetRestaurantByID(ctx, id)
	if err != nil {
		return nil, toStatus(err)
	}
	return &restaurantpb.RestaurantResponse{Restaurant: toPBRestaurant(r)}, nil
}

func (s *Server) UpdateRestaurant(ctx context.Context, req *restaurantpb.UpdateRestaurantRequest) (*restaurantpb.RestaurantResponse, error) {
	id, ownerID, err := parsePair(req.Id, req.OwnerId)
	if err != nil {
		return nil, toStatus(err)
	}
	r, err := s.uc.UpdateRestaurant(ctx, ownerID, &entities.Restaurant{
		ID: id, Name: req.Name, Description: req.Description, CuisineType: req.CuisineType,
		Address: req.Address, City: req.City, ImageURL: req.ImageUrl, IsOpen: req.IsOpen, OwnerID: ownerID,
	})
	if err != nil {
		return nil, toStatus(err)
	}
	return &restaurantpb.RestaurantResponse{Restaurant: toPBRestaurant(r)}, nil
}

func (s *Server) DeleteRestaurant(ctx context.Context, req *restaurantpb.DeleteRestaurantRequest) (*restaurantpb.DeleteResponse, error) {
	id, ownerID, err := parsePair(req.Id, req.OwnerId)
	if err != nil {
		return nil, toStatus(err)
	}
	if err := s.uc.DeleteRestaurant(ctx, id, ownerID); err != nil {
		return nil, toStatus(err)
	}
	return &restaurantpb.DeleteResponse{Success: true, Message: "restaurant deleted"}, nil
}

func (s *Server) ListRestaurants(ctx context.Context, req *restaurantpb.ListRestaurantsRequest) (*restaurantpb.ListRestaurantsResponse, error) {
	items, meta, err := s.uc.ListRestaurants(ctx, filterFromPB(req.Filter), page(req.Pagination), size(req.Pagination))
	if err != nil {
		return nil, toStatus(err)
	}
	return &restaurantpb.ListRestaurantsResponse{Restaurants: toPBRestaurants(items), Meta: &restaurantpb.PageMeta{Page: int32(meta.Page), PageSize: int32(meta.PageSize), Total: meta.Total}}, nil
}

func (s *Server) SearchRestaurants(ctx context.Context, req *restaurantpb.SearchRestaurantsRequest) (*restaurantpb.ListRestaurantsResponse, error) {
	items, meta, err := s.uc.SearchRestaurants(ctx, req.Query, filterFromPB(req.Filter), page(req.Pagination), size(req.Pagination))
	if err != nil {
		return nil, toStatus(err)
	}
	return &restaurantpb.ListRestaurantsResponse{Restaurants: toPBRestaurants(items), Meta: &restaurantpb.PageMeta{Page: int32(meta.Page), PageSize: int32(meta.PageSize), Total: meta.Total}}, nil
}

func (s *Server) GetTopRatedRestaurants(ctx context.Context, req *restaurantpb.TopRatedRestaurantsRequest) (*restaurantpb.ListRestaurantsResponse, error) {
	items, err := s.uc.TopRated(ctx, req.City, int(req.Limit))
	if err != nil {
		return nil, toStatus(err)
	}
	return &restaurantpb.ListRestaurantsResponse{Restaurants: toPBRestaurants(items), Meta: &restaurantpb.PageMeta{Page: 1, PageSize: int32(len(items)), Total: int64(len(items))}}, nil
}

func (s *Server) CreateCategory(ctx context.Context, req *restaurantpb.CreateCategoryRequest) (*restaurantpb.CategoryResponse, error) {
	restaurantID, ownerID, err := parsePair(req.RestaurantId, req.OwnerId)
	if err != nil {
		return nil, toStatus(err)
	}
	c, err := s.uc.CreateCategory(ctx, ownerID, &entities.Category{RestaurantID: restaurantID, Name: req.Name})
	if err != nil {
		return nil, toStatus(err)
	}
	return &restaurantpb.CategoryResponse{Category: toPBCategory(c)}, nil
}

func (s *Server) UpdateCategory(ctx context.Context, req *restaurantpb.UpdateCategoryRequest) (*restaurantpb.CategoryResponse, error) {
	id, ownerID, err := parsePair(req.Id, req.OwnerId)
	if err != nil {
		return nil, toStatus(err)
	}
	c, err := s.uc.UpdateCategory(ctx, ownerID, &entities.Category{ID: id, Name: req.Name})
	if err != nil {
		return nil, toStatus(err)
	}
	return &restaurantpb.CategoryResponse{Category: toPBCategory(c)}, nil
}

func (s *Server) DeleteCategory(ctx context.Context, req *restaurantpb.DeleteCategoryRequest) (*restaurantpb.DeleteResponse, error) {
	id, ownerID, err := parsePair(req.Id, req.OwnerId)
	if err != nil {
		return nil, toStatus(err)
	}
	if err := s.uc.DeleteCategory(ctx, id, ownerID); err != nil {
		return nil, toStatus(err)
	}
	return &restaurantpb.DeleteResponse{Success: true, Message: "category deleted"}, nil
}

func (s *Server) ListCategories(ctx context.Context, req *restaurantpb.ListCategoriesRequest) (*restaurantpb.ListCategoriesResponse, error) {
	id, err := uuid.Parse(req.RestaurantId)
	if err != nil {
		return nil, toStatus(services.ErrInvalidInput)
	}
	items, err := s.uc.ListCategories(ctx, id)
	if err != nil {
		return nil, toStatus(err)
	}
	return &restaurantpb.ListCategoriesResponse{Categories: toPBCategories(items)}, nil
}

func (s *Server) CreateMenuItem(ctx context.Context, req *restaurantpb.CreateMenuItemRequest) (*restaurantpb.MenuItemResponse, error) {
	restaurantID, ownerID, err := parsePair(req.RestaurantId, req.OwnerId)
	if err != nil {
		return nil, toStatus(err)
	}
	categoryID, err := uuid.Parse(req.CategoryId)
	if err != nil {
		return nil, toStatus(services.ErrInvalidInput)
	}
	item, err := s.uc.CreateMenuItem(ctx, ownerID, &entities.MenuItem{RestaurantID: restaurantID, CategoryID: categoryID, Name: req.Name, Description: req.Description, Price: req.Price, ImageURL: req.ImageUrl, IsAvailable: req.IsAvailable})
	if err != nil {
		return nil, toStatus(err)
	}
	return &restaurantpb.MenuItemResponse{Item: toPBMenuItem(item)}, nil
}

func (s *Server) UpdateMenuItem(ctx context.Context, req *restaurantpb.UpdateMenuItemRequest) (*restaurantpb.MenuItemResponse, error) {
	id, ownerID, err := parsePair(req.Id, req.OwnerId)
	if err != nil {
		return nil, toStatus(err)
	}
	categoryID, err := uuid.Parse(req.CategoryId)
	if err != nil {
		return nil, toStatus(services.ErrInvalidInput)
	}
	item, err := s.uc.UpdateMenuItem(ctx, ownerID, &entities.MenuItem{ID: id, CategoryID: categoryID, Name: req.Name, Description: req.Description, Price: req.Price, ImageURL: req.ImageUrl, IsAvailable: req.IsAvailable})
	if err != nil {
		return nil, toStatus(err)
	}
	return &restaurantpb.MenuItemResponse{Item: toPBMenuItem(item)}, nil
}

func (s *Server) DeleteMenuItem(ctx context.Context, req *restaurantpb.DeleteMenuItemRequest) (*restaurantpb.DeleteResponse, error) {
	id, ownerID, err := parsePair(req.Id, req.OwnerId)
	if err != nil {
		return nil, toStatus(err)
	}
	if err := s.uc.DeleteMenuItem(ctx, id, ownerID); err != nil {
		return nil, toStatus(err)
	}
	return &restaurantpb.DeleteResponse{Success: true, Message: "menu item deleted"}, nil
}

func (s *Server) GetMenuByRestaurant(ctx context.Context, req *restaurantpb.GetMenuByRestaurantRequest) (*restaurantpb.MenuResponse, error) {
	id, err := uuid.Parse(req.RestaurantId)
	if err != nil {
		return nil, toStatus(services.ErrInvalidInput)
	}
	menu, err := s.uc.GetMenuByRestaurant(ctx, id)
	if err != nil {
		return nil, toStatus(err)
	}
	return &restaurantpb.MenuResponse{Categories: toPBMenu(menu)}, nil
}

func (s *Server) SetMenuItemAvailability(ctx context.Context, req *restaurantpb.SetMenuItemAvailabilityRequest) (*restaurantpb.MenuItemResponse, error) {
	id, ownerID, err := parsePair(req.Id, req.OwnerId)
	if err != nil {
		return nil, toStatus(err)
	}
	item, err := s.uc.SetMenuItemAvailability(ctx, id, ownerID, req.IsAvailable)
	if err != nil {
		return nil, toStatus(err)
	}
	return &restaurantpb.MenuItemResponse{Item: toPBMenuItem(item)}, nil
}

func (s *Server) HealthCheck(context.Context, *restaurantpb.HealthCheckRequest) (*restaurantpb.HealthCheckResponse, error) {
	return &restaurantpb.HealthCheckResponse{Status: "SERVING", Service: s.serviceName}, nil
}

func parsePair(a, b string) (uuid.UUID, uuid.UUID, error) {
	left, err := uuid.Parse(a)
	if err != nil {
		return uuid.Nil, uuid.Nil, services.ErrInvalidInput
	}
	right, err := uuid.Parse(b)
	if err != nil {
		return uuid.Nil, uuid.Nil, services.ErrInvalidInput
	}
	return left, right, nil
}

func page(p *restaurantpb.Pagination) int {
	if p == nil {
		return 1
	}
	return int(p.Page)
}

func size(p *restaurantpb.Pagination) int {
	if p == nil {
		return 20
	}
	return int(p.PageSize)
}

func filterFromPB(f *restaurantpb.RestaurantFilter) entities.RestaurantFilter {
	if f == nil {
		return entities.RestaurantFilter{}
	}
	return entities.RestaurantFilter{Query: f.Query, CuisineType: f.CuisineType, City: f.City, OpenOnly: f.OpenOnly, SortBy: f.SortBy, SortDirection: f.SortDirection}
}
