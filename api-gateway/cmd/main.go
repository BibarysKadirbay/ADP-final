package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userpb "food-delivery-system/proto/userpb"
)

func main() {

	conn, err := grpc.NewClient(
		"user-service:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		panic(err)
	}

	userClient := userpb.NewUserServiceClient(conn)
	orderConn, err := grpc.NewClient(
		"order-service:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	orderClient := orderpb.NewOrderServiceClient(orderConn)
	r := gin.Default()

	r.POST("/register", func(c *gin.Context) {

		var body struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		res, err := userClient.RegisterUser(
			context.Background(),
			&userpb.RegisterRequest{
				Name:     body.Name,
				Email:    body.Email,
				Password: body.Password,
			},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	r.POST("/login", func(c *gin.Context) {

		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		res, err := userClient.LoginUser(
			context.Background(),
			&userpb.LoginRequest{
				Email:    body.Email,
				Password: body.Password,
			},
		)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res)
	})
	r.POST("/orders", func(c *gin.Context) {

		var body struct {
			UserID       string `json:"user_id"`
			RestaurantID string `json:"restaurant_id"`
			TotalPrice   int64  `json:"total_price"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		res, err := orderClient.CreateOrder(
			context.Background(),
			&orderpb.CreateOrderRequest{
				UserId:       body.UserID,
				RestaurantId: body.RestaurantID,
				TotalPrice:   body.TotalPrice,
			},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res)
	})
	r.Run(":8080")
}
