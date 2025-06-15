package model

import "github.com/golang-jwt/jwt/v5"

type AdminResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type AdminCreateRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AdminUpdateRequest struct {
	ID uint `json:"id" validate:"required"`
	AdminCreateRequest
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type TokenPyload struct {
	AdminID  uint   `json:"admin_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}
