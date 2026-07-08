package controller

import (
	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/service"
)

type AuthController interface {
	Register(email, password, name string) (*dto.UserResponse, error)
	Login(email, password string) (*dto.LoginResponse, error)
	Profile(userID uint) (*dto.UserResponse, error)
}

type authController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return &authController{authService: authService}
}

func (c *authController) Register(email, password, name string) (*dto.UserResponse, error) {
	user, err := c.authService.RegisterUser(email, password, name)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (c *authController) Login(email, password string) (*dto.LoginResponse, error) {
	token, user, err := c.authService.LoginUser(email, password)
	if err != nil {
		return nil, err
	}

	userDTO := &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}

	return &dto.LoginResponse{
		Token: token,
		User:  *userDTO,
	}, nil
}

func (c *authController) Profile(userID uint) (*dto.UserResponse, error) {
	user, err := c.authService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}
