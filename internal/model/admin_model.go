package model

import "github.com/golang-jwt/jwt/v5"

type AdminResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type AdminCreateRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AdminUpdateRequest struct {
	ID uint `json:"id" validate:"required"`
	AdminCreateRequest
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type TokenPyload struct {
	AdminID uint   `json:"admin_id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	jwt.RegisteredClaims
}
