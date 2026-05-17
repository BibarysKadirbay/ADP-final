package main

import (
	"log"

	"api-gateway/internal/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
	orderpb "api-gateway/internal/grpc/orderpb"
	restaurantpb "api-gateway/internal/grpc/restaurantpb"
	userpb "api-gateway/internal/grpc/userpb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()

	userConn := dial(cfg.UserGRPCAddr)
	defer userConn.Close()
	orderConn := dial(cfg.OrderGRPCAddr)
	defer orderConn.Close()
	restaurantConn := dial(cfg.RestaurantGRPCAddr)
	defer restaurantConn.Close()

	authHandler := handlers.NewAuthHandler(userpb.NewUserServiceClient(userConn))
	restaurantHandler := handlers.NewRestaurantHandler(restaurantpb.NewRestaurantServiceClient(restaurantConn))
	orderHandler := handlers.NewOrderHandler(orderpb.NewOrderServiceClient(orderConn))

	r := gin.Default()
	r.Use(corsMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "api-gateway running"})
	})

	// Team Member 1 — Auth & Users
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	// Team Member 2 — Restaurants
	r.POST("/restaurants", middleware.JWTAuth(cfg.JWTSecret), restaurantHandler.CreateRestaurant)
	r.GET("/restaurants", restaurantHandler.ListRestaurants)
	r.GET("/restaurants/:id", restaurantHandler.GetRestaurant)
	r.POST("/restaurants/:id/menu", middleware.JWTAuth(cfg.JWTSecret), restaurantHandler.AddMenuItem)
	r.GET("/restaurants/:id/menu", restaurantHandler.GetMenu)
	r.PATCH("/restaurants/:id/menu/:menuItemId/availability", middleware.JWTAuth(cfg.JWTSecret), restaurantHandler.UpdateAvailability)

	// Team Member 3 — Orders
	r.POST("/orders", middleware.JWTAuth(cfg.JWTSecret), orderHandler.CreateOrder)
	r.GET("/orders/:id", orderHandler.GetOrder)
	r.GET("/users/:id/orders", orderHandler.GetUserOrders)
	r.GET("/users/:id", authHandler.GetUser)
	r.PATCH("/orders/:id/status", middleware.JWTAuth(cfg.JWTSecret), orderHandler.UpdateStatus)
	r.PATCH("/orders/:id/cancel", middleware.JWTAuth(cfg.JWTSecret), orderHandler.CancelOrder)

	log.Println("api-gateway listening on :" + cfg.HTTPPort)
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal(err)
	}
}

func dial(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc dial %s: %v", addr, err)
	}
	return conn
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
