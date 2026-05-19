package handlers

import (
	"context"
	"net/http"
	"time"

	paymentpb "api-gateway/internal/grpc/paymentpb"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	payments paymentpb.PaymentServiceClient
}

func NewPaymentHandler(payments paymentpb.PaymentServiceClient) *PaymentHandler {
	return &PaymentHandler{payments: payments}
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.payments.GetPayment(ctx, &paymentpb.GetPaymentRequest{PaymentId: c.Param("id")})
	if err != nil {
		writeGRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *PaymentHandler) GetOrderPayments(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.payments.ListPaymentsByOrder(ctx, &paymentpb.ListPaymentsByOrderRequest{OrderId: c.Param("id")})
	if err != nil {
		writeGRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}
