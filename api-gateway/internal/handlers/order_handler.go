package handlers

import (
	"context"
	"net/http"
	"time"

	orderpb "api-gateway/internal/grpc/orderpb"

	"github.com/gin-gonic/gin"
)

// OrderHandler — Team Member 3: order routes.
type OrderHandler struct {
	orders orderpb.OrderServiceClient
}

func NewOrderHandler(orders orderpb.OrderServiceClient) *OrderHandler {
	return &OrderHandler{orders: orders}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var body struct {
		UserID       string `json:"user_id"`
		RestaurantID string `json:"restaurant_id"`
		TotalPrice   int64  `json:"total_price"`
		Address      string `json:"address"`
		Comment      string `json:"comment"`
		UserEmail    string `json:"user_email"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.UserID == "" {
		if uid, ok := c.Get("user_id"); ok {
			body.UserID = uid.(string)
		}
	}
	if body.UserEmail == "" {
		if email, ok := c.Get("email"); ok {
			body.UserEmail = email.(string)
		}
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	res, err := h.orders.CreateOrder(ctx, &orderpb.CreateOrderRequest{
		UserId: body.UserID, RestaurantId: body.RestaurantID, TotalPrice: body.TotalPrice,
		Address: body.Address, Comment: body.Comment, UserEmail: body.UserEmail,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	res, err := h.orders.GetOrder(ctx, &orderpb.GetOrderRequest{OrderId: c.Param("id")})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	res, err := h.orders.GetOrdersByUser(ctx, &orderpb.GetOrdersByUserRequest{UserId: c.Param("id")})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	var body struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	res, err := h.orders.UpdateOrderStatus(ctx, &orderpb.UpdateOrderStatusRequest{
		OrderId: c.Param("id"), Status: body.Status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	res, err := h.orders.CancelOrder(ctx, &orderpb.CancelOrderRequest{OrderId: c.Param("id")})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
