package handler

import (
	"net/http"
	"strconv"

	"github.com/entertrans/backend-bogor.git/internal/controller"
	"github.com/entertrans/backend-bogor.git/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Profile(c *gin.Context)
}

type authHandler struct {
	authController controller.AuthController
}

func NewAuthHandler(authController controller.AuthController) AuthHandler {
	return &authHandler{authController: authController}
}

func (h *authHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userDTO, err := h.authController.Register(req.Email, req.Password, req.Name)
	if err != nil {
		response.SendErrorResponse(c, http.StatusConflict, err.Error())
		return
	}

	response.SendResponse(c, http.StatusCreated, "User registered successfully", userDTO)
}

func (h *authHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	loginDTO, err := h.authController.Login(req.Email, req.Password)
	if err != nil {
		response.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.SendResponse(c, http.StatusOK, "Login successful", loginDTO)
}

func (h *authHandler) Profile(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.SendErrorResponse(c, http.StatusUnauthorized, "User not found in context")
		return
	}

	userID, err := strconv.ParseUint(userIDVal.(string), 10, 32)
	if err != nil {
		response.SendErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	userDTO, err := h.authController.Profile(uint(userID))
	if err != nil {
		response.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	response.SendResponse(c, http.StatusOK, "Profile retrieved successfully", userDTO)
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
