package restaurantpb

import (
	"context"

	"google.golang.org/grpc"
)

type RestaurantServiceServer interface {
	CreateRestaurant(context.Context, *CreateRestaurantRequest) (*RestaurantResponse, error)
	GetRestaurantById(context.Context, *GetRestaurantByIdRequest) (*RestaurantResponse, error)
	UpdateRestaurant(context.Context, *UpdateRestaurantRequest) (*RestaurantResponse, error)
	DeleteRestaurant(context.Context, *DeleteRestaurantRequest) (*DeleteResponse, error)
	ListRestaurants(context.Context, *ListRestaurantsRequest) (*ListRestaurantsResponse, error)
	SearchRestaurants(context.Context, *SearchRestaurantsRequest) (*ListRestaurantsResponse, error)
	GetTopRatedRestaurants(context.Context, *TopRatedRestaurantsRequest) (*ListRestaurantsResponse, error)
	CreateCategory(context.Context, *CreateCategoryRequest) (*CategoryResponse, error)
	UpdateCategory(context.Context, *UpdateCategoryRequest) (*CategoryResponse, error)
	DeleteCategory(context.Context, *DeleteCategoryRequest) (*DeleteResponse, error)
	ListCategories(context.Context, *ListCategoriesRequest) (*ListCategoriesResponse, error)
	CreateMenuItem(context.Context, *CreateMenuItemRequest) (*MenuItemResponse, error)
	UpdateMenuItem(context.Context, *UpdateMenuItemRequest) (*MenuItemResponse, error)
	DeleteMenuItem(context.Context, *DeleteMenuItemRequest) (*DeleteResponse, error)
	GetMenuByRestaurant(context.Context, *GetMenuByRestaurantRequest) (*MenuResponse, error)
	SetMenuItemAvailability(context.Context, *SetMenuItemAvailabilityRequest) (*MenuItemResponse, error)
	HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error)
}

type UnimplementedRestaurantServiceServer struct{}

func RegisterRestaurantServiceServer(s grpc.ServiceRegistrar, srv RestaurantServiceServer) {
	s.RegisterService(&RestaurantService_ServiceDesc, srv)
}

var RestaurantService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "restaurant.v1.RestaurantService",
	HandlerType: (*RestaurantServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateRestaurant", Handler: createRestaurantHandler},
		{MethodName: "GetRestaurantById", Handler: getRestaurantByIDHandler},
		{MethodName: "UpdateRestaurant", Handler: updateRestaurantHandler},
		{MethodName: "DeleteRestaurant", Handler: deleteRestaurantHandler},
		{MethodName: "ListRestaurants", Handler: listRestaurantsHandler},
		{MethodName: "SearchRestaurants", Handler: searchRestaurantsHandler},
		{MethodName: "GetTopRatedRestaurants", Handler: topRatedRestaurantsHandler},
		{MethodName: "CreateCategory", Handler: createCategoryHandler},
		{MethodName: "UpdateCategory", Handler: updateCategoryHandler},
		{MethodName: "DeleteCategory", Handler: deleteCategoryHandler},
		{MethodName: "ListCategories", Handler: listCategoriesHandler},
		{MethodName: "CreateMenuItem", Handler: createMenuItemHandler},
		{MethodName: "UpdateMenuItem", Handler: updateMenuItemHandler},
		{MethodName: "DeleteMenuItem", Handler: deleteMenuItemHandler},
		{MethodName: "GetMenuByRestaurant", Handler: getMenuByRestaurantHandler},
		{MethodName: "SetMenuItemAvailability", Handler: setMenuItemAvailabilityHandler},
		{MethodName: "HealthCheck", Handler: healthCheckHandler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/restaurant.proto",
}

func unary[Req any, Resp any](srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor, info *grpc.UnaryServerInfo, call func(RestaurantServiceServer, context.Context, *Req) (*Resp, error)) (any, error) {
	in := new(Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return call(srv.(RestaurantServiceServer), ctx, in)
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return call(srv.(RestaurantServiceServer), ctx, req.(*Req))
	}
	return interceptor(ctx, in, info, handler)
}

func createRestaurantHandler(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
	return unary[CreateRestaurantRequest, RestaurantResponse](s, c, d, i, &grpc.UnaryServerInfo{Server: s, FullMethod: "/restaurant.v1.RestaurantService/CreateRestaurant"}, func(s RestaurantServiceServer, c context.Context, r *CreateRestaurantRequest) (*RestaurantResponse, error) {
		return s.CreateRestaurant(c, r)
	})
}
func getRestaurantByIDHandler(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
	return unary[GetRestaurantByIdRequest, RestaurantResponse](s, c, d, i, &grpc.UnaryServerInfo{Server: s, FullMethod: "/restaurant.v1.RestaurantService/GetRestaurantById"}, func(s RestaurantServiceServer, c context.Context, r *GetRestaurantByIdRequest) (*RestaurantResponse, error) {
		return s.GetRestaurantById(c, r)
	})
}
func updateRestaurantHandler(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
	return unary[UpdateRestaurantRequest, RestaurantResponse](s, c, d, i, &grpc.UnaryServerInfo{Server: s, FullMethod: "/restaurant.v1.RestaurantService/UpdateRestaurant"}, func(s RestaurantServiceServer, c context.Context, r *UpdateRestaurantRequest) (*RestaurantResponse, error) {
		return s.UpdateRestaurant(c, r)
	})
}
func deleteRestaurantHandler(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
	return unary[DeleteRestaurantRequest, DeleteResponse](s, c, d, i, &grpc.UnaryServerInfo{Server: s, FullMethod: "/restaurant.v1.RestaurantService/DeleteRestaurant"}, func(s RestaurantServiceServer, c context.Context, r *DeleteRestaurantRequest) (*DeleteResponse, error) {
		return s.DeleteRestaurant(c, r)
	})
}
func listRestaurantsHandler(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
	return unary[ListRestaurantsRequest, ListRestaurantsResponse](s, c, d, i, &grpc.UnaryServerInfo{Server: s, FullMethod: "/restaurant.v1.RestaurantService/ListRestaurants"}, func(s RestaurantServiceServer, c context.Context, r *ListRestaurantsRequest) (*ListRestaurantsResponse, error) {
		return s.ListRestaurants(c, r)
	})
}
func searchRestaurantsHandler(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
	return unary[SearchRestaurantsRequest, ListRestaurantsResponse](s, c, d, i, &grpc.UnaryServerInfo{Server: s, FullMethod: "/restaurant.v1.RestaurantService/SearchRestaurants"}, func(s RestaurantServiceServer, c context.Context, r *SearchRestaurantsRequest) (*ListRestaurantsResponse, error) {
		return s.SearchRestaurants(c, r)
	})
}
func topRatedRestaurantsHandler(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
	return unary[TopRatedRestaurantsRequest, ListRestaurantsResponse](s, c, d, i, &grpc.UnaryServerInfo{Server: s, FullMethod: "/restaurant.v1.RestaurantService/GetTopRatedRestaurants"}, func(s RestaurantServiceServer, c context.Context, r *TopRatedRestaurantsRequest) (*ListRestaurantsResponse, error) {
		return s.GetTopRatedRestaurants(c, r)
	})
}
func createCategoryHandler(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
	return unary[CreateCategoryRequest, CategoryResponse](s, c, d, i, &grpc.UnaryServerInfo{Server: s, FullMethod: "/restaurant.v1.RestaurantService/CreateCategory"}, func(s RestaurantServiceServer, c context.Context, r *CreateCategoryRequest) (*CategoryResponse, error) {
		return s.CreateCategory(c, r)
	})
}
func updateCategoryHandler(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
	return unary[UpdateCategoryRequest, CategoryResponse](s, c, d, i, &grpc.UnaryServerInfo{Server: s, FullMethod: "/restaurant.v1.RestaurantService/UpdateCategory"}, func(s RestaurantServiceServer, c context.Context, r *UpdateCategoryRequest) (*CategoryResponse, error) {
		return s.UpdateCategory(c, r)
	})
}
func deleteCategoryHandler(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
	return unary[DeleteCategoryRequest, DeleteResponse](s, c, d, i, &grpc.UnaryServerInfo{Server: s, FullMethod: "/restaurant.v1.RestaurantService/DeleteCategory"}, func(s RestaurantServiceServer, c context.Context, r *DeleteCategoryRequest) (*DeleteResponse, error) {
		return s.DeleteCategory(c, r)
	})
}
func listCategoriesHandler(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
	return unary[ListCategoriesRequest, ListCategoriesResponse](s, c, d, i, &grpc.UnaryServerInfo{Server: s, FullMethod: "/restaurant.v1.RestaurantService/ListCategories"}, func(s RestaurantServiceServer, c context.Context, r *ListCategoriesRequest) (*ListCategoriesResponse, error) {
		return s.ListCategories(c, r)
	})
}
func createMenuItemHandler(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
	return unary[CreateMenuItemRequest, MenuItemResponse](s, c, d, i, &grpc.UnaryServerInfo{Server: s, FullMethod: "/restaurant.v1.RestaurantService/CreateMenuItem"}, func(s RestaurantServiceServer, c context.Context, r *CreateMenuItemRequest) (*MenuItemResponse, error) {
		return s.CreateMenuItem(c, r)
	})
}
func updateMenuItemHandler(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
	return unary[UpdateMenuItemRequest, MenuItemResponse](s, c, d, i, &grpc.UnaryServerInfo{Server: s, FullMethod: "/restaurant.v1.RestaurantService/UpdateMenuItem"}, func(s RestaurantServiceServer, c context.Context, r *UpdateMenuItemRequest) (*MenuItemResponse, error) {
		return s.UpdateMenuItem(c, r)
	})
}
func deleteMenuItemHandler(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
	return unary[DeleteMenuItemRequest, DeleteResponse](s, c, d, i, &grpc.UnaryServerInfo{Server: s, FullMethod: "/restaurant.v1.RestaurantService/DeleteMenuItem"}, func(s RestaurantServiceServer, c context.Context, r *DeleteMenuItemRequest) (*DeleteResponse, error) {
		return s.DeleteMenuItem(c, r)
	})
}
func getMenuByRestaurantHandler(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
	return unary[GetMenuByRestaurantRequest, MenuResponse](s, c, d, i, &grpc.UnaryServerInfo{Server: s, FullMethod: "/restaurant.v1.RestaurantService/GetMenuByRestaurant"}, func(s RestaurantServiceServer, c context.Context, r *GetMenuByRestaurantRequest) (*MenuResponse, error) {
		return s.GetMenuByRestaurant(c, r)
	})
}
func setMenuItemAvailabilityHandler(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
	return unary[SetMenuItemAvailabilityRequest, MenuItemResponse](s, c, d, i, &grpc.UnaryServerInfo{Server: s, FullMethod: "/restaurant.v1.RestaurantService/SetMenuItemAvailability"}, func(s RestaurantServiceServer, c context.Context, r *SetMenuItemAvailabilityRequest) (*MenuItemResponse, error) {
		return s.SetMenuItemAvailability(c, r)
	})
}
func healthCheckHandler(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
	return unary[HealthCheckRequest, HealthCheckResponse](s, c, d, i, &grpc.UnaryServerInfo{Server: s, FullMethod: "/restaurant.v1.RestaurantService/HealthCheck"}, func(s RestaurantServiceServer, c context.Context, r *HealthCheckRequest) (*HealthCheckResponse, error) {
		return s.HealthCheck(c, r)
	})
}
