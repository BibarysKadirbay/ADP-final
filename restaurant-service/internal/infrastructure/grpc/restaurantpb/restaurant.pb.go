package restaurantpb

import "google.golang.org/protobuf/types/known/timestamppb"

type Pagination struct {
	Page     int32
	PageSize int32
}

type PageMeta struct {
	Page     int32
	PageSize int32
	Total    int64
}

type Restaurant struct {
	Id           string
	OwnerId      string
	Name         string
	Description  string
	CuisineType  string
	Address      string
	City         string
	Rating       float64
	TotalReviews int32
	ImageUrl     string
	IsOpen       bool
	CreatedAt    *timestamppb.Timestamp
	UpdatedAt    *timestamppb.Timestamp
}

type Category struct {
	Id           string
	RestaurantId string
	Name         string
	CreatedAt    *timestamppb.Timestamp
}

type MenuItem struct {
	Id           string
	CategoryId   string
	RestaurantId string
	Name         string
	Description  string
	Price        float64
	ImageUrl     string
	IsAvailable  bool
	CreatedAt    *timestamppb.Timestamp
	UpdatedAt    *timestamppb.Timestamp
}

type MenuCategory struct {
	Category *Category
	Items    []*MenuItem
}

type RestaurantFilter struct {
	Query         string
	CuisineType   string
	City          string
	OpenOnly      *bool
	SortBy        string
	SortDirection string
}

type CreateRestaurantRequest struct {
	OwnerId     string
	Name        string
	Description string
	CuisineType string
	Address     string
	City        string
	ImageUrl    string
	IsOpen      bool
}
type GetRestaurantByIdRequest struct{ Id string }
type UpdateRestaurantRequest struct {
	Id          string
	OwnerId     string
	Name        string
	Description string
	CuisineType string
	Address     string
	City        string
	ImageUrl    string
	IsOpen      bool
}
type DeleteRestaurantRequest struct {
	Id      string
	OwnerId string
}
type ListRestaurantsRequest struct {
	Pagination *Pagination
	Filter     *RestaurantFilter
}
type SearchRestaurantsRequest struct {
	Query      string
	Pagination *Pagination
	Filter     *RestaurantFilter
}
type TopRatedRestaurantsRequest struct {
	City  string
	Limit int32
}
type RestaurantResponse struct{ Restaurant *Restaurant }
type ListRestaurantsResponse struct {
	Restaurants []*Restaurant
	Meta        *PageMeta
}

type CreateCategoryRequest struct {
	RestaurantId string
	OwnerId      string
	Name         string
}
type UpdateCategoryRequest struct {
	Id      string
	OwnerId string
	Name    string
}
type DeleteCategoryRequest struct {
	Id      string
	OwnerId string
}
type ListCategoriesRequest struct{ RestaurantId string }
type CategoryResponse struct{ Category *Category }
type ListCategoriesResponse struct{ Categories []*Category }

type CreateMenuItemRequest struct {
	RestaurantId string
	CategoryId   string
	OwnerId      string
	Name         string
	Description  string
	Price        float64
	ImageUrl     string
	IsAvailable  bool
}
type UpdateMenuItemRequest struct {
	Id          string
	OwnerId     string
	CategoryId  string
	Name        string
	Description string
	Price       float64
	ImageUrl    string
	IsAvailable bool
}
type DeleteMenuItemRequest struct {
	Id      string
	OwnerId string
}
type GetMenuByRestaurantRequest struct{ RestaurantId string }
type SetMenuItemAvailabilityRequest struct {
	Id          string
	OwnerId     string
	IsAvailable bool
}
type MenuItemResponse struct{ Item *MenuItem }
type MenuResponse struct{ Categories []*MenuCategory }
type DeleteResponse struct {
	Success bool
	Message string
}
type HealthCheckRequest struct{}
type HealthCheckResponse struct {
	Status  string
	Service string
}
