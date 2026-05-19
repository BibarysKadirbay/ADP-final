package handlers

import (
	"context"
	"net/http"
	"time"

	deliverypb "api-gateway/internal/grpc/deliverypb"

	"github.com/gin-gonic/gin"
)

type DeliveryHandler struct {
	deliveries deliverypb.DeliveryServiceClient
}

func NewDeliveryHandler(deliveries deliverypb.DeliveryServiceClient) *DeliveryHandler {
	return &DeliveryHandler{deliveries: deliveries}
}

func (h *DeliveryHandler) GetDelivery(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.deliveries.GetDeliveryById(ctx, &deliverypb.GetDeliveryByIdRequest{Id: c.Param("id")})
	if err != nil {
		writeGRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *DeliveryHandler) GetOrderDeliveries(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.deliveries.GetDeliveriesByOrder(ctx, &deliverypb.GetDeliveriesByOrderRequest{OrderId: c.Param("id")})
	if err != nil {
		writeGRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *DeliveryHandler) UpdateStatus(c *gin.Context) {
	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.deliveries.UpdateDeliveryStatus(ctx, &deliverypb.UpdateDeliveryStatusRequest{
		Id:     c.Param("id"),
		Status: body.Status,
	})
	if err != nil {
		writeGRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}
