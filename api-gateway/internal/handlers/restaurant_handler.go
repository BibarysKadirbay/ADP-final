package handlers

import (
	"context"
	"net/http"
	"time"

	restaurantpb "api-gateway/internal/grpc/restaurantpb"

	"github.com/gin-gonic/gin"
)

// RestaurantHandler — Team Member 2: restaurant & menu routes.
type RestaurantHandler struct {
	restaurants restaurantpb.RestaurantServiceClient
}

func NewRestaurantHandler(restaurants restaurantpb.RestaurantServiceClient) *RestaurantHandler {
	return &RestaurantHandler{restaurants: restaurants}
}

func (h *RestaurantHandler) CreateRestaurant(c *gin.Context) {
	var body struct {
		OwnerID     string `json:"owner_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CuisineType string `json:"cuisine_type"`
		Address     string `json:"address"`
		City        string `json:"city"`
		ImageURL    string `json:"image_url"`
		IsOpen      bool   `json:"is_open"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.OwnerID == "" {
		if uid, ok := c.Get("user_id"); ok {
			body.OwnerID = uid.(string)
		}
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	res, err := h.restaurants.CreateRestaurant(ctx, &restaurantpb.CreateRestaurantRequest{
		OwnerId: body.OwnerID, Name: body.Name, Description: body.Description,
		CuisineType: body.CuisineType, Address: body.Address, City: body.City,
		ImageUrl: body.ImageURL, IsOpen: body.IsOpen,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *RestaurantHandler) ListRestaurants(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	res, err := h.restaurants.ListRestaurants(ctx, &restaurantpb.ListRestaurantsRequest{
		Pagination: &restaurantpb.Pagination{Page: 1, PageSize: 50},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *RestaurantHandler) GetRestaurant(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	res, err := h.restaurants.GetRestaurantById(ctx, &restaurantpb.GetRestaurantByIdRequest{Id: c.Param("id")})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *RestaurantHandler) AddMenuItem(c *gin.Context) {
	var body struct {
		CategoryID  string  `json:"category_id"`
		OwnerID     string  `json:"owner_id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		ImageURL    string  `json:"image_url"`
		IsAvailable bool    `json:"is_available"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.OwnerID == "" {
		if uid, ok := c.Get("user_id"); ok {
			body.OwnerID = uid.(string)
		}
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	res, err := h.restaurants.CreateMenuItem(ctx, &restaurantpb.CreateMenuItemRequest{
		RestaurantId: c.Param("id"), CategoryId: body.CategoryID, OwnerId: body.OwnerID,
		Name: body.Name, Description: body.Description, Price: body.Price,
		ImageUrl: body.ImageURL, IsAvailable: body.IsAvailable,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *RestaurantHandler) GetMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	res, err := h.restaurants.GetMenuByRestaurant(ctx, &restaurantpb.GetMenuByRestaurantRequest{
		RestaurantId: c.Param("id"),
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *RestaurantHandler) UpdateAvailability(c *gin.Context) {
	var body struct {
		OwnerID     string `json:"owner_id"`
		IsAvailable bool   `json:"is_available"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	res, err := h.restaurants.SetMenuItemAvailability(ctx, &restaurantpb.SetMenuItemAvailabilityRequest{
		Id: c.Param("menuItemId"), OwnerId: body.OwnerID, IsAvailable: body.IsAvailable,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
