package handlers

import (
	"context"
	"net/http"
	"time"

	userpb "api-gateway/internal/grpc/userpb"

	"github.com/gin-gonic/gin"
)

// AuthHandler — Team Member 1: user auth & profile routes.
type AuthHandler struct {
	users userpb.UserServiceClient
}

func NewAuthHandler(users userpb.UserServiceClient) *AuthHandler {
	return &AuthHandler{users: users}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
		Address  string `json:"address"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	res, err := h.users.RegisterUser(ctx, &userpb.RegisterUserRequest{
		Name: body.Name, Email: body.Email, Password: body.Password,
		Phone: body.Phone, Address: body.Address,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	res, err := h.users.LoginUser(ctx, &userpb.LoginUserRequest{
		Email: body.Email, Password: body.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) GetUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	res, err := h.users.GetUser(ctx, &userpb.GetUserRequest{UserId: c.Param("id")})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
