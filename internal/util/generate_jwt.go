package util

import (
	"os"
	"strconv"
	"time"

	"github.com/Bangdams/web-profile-API/internal/entity"
	"github.com/Bangdams/web-profile-API/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(request *entity.Admin) (string, error) {
	var token model.TokenPyload
	duration := os.Getenv("DURATION_JWT_ACCESS_TOKEN")
	lifeTime, _ := strconv.Atoi(duration)

	now := time.Now()
	token.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    "QuizKu",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * time.Duration(lifeTime))),
	}

	token.AdminID = request.ID
	token.Username = request.Username
	token.Name = request.Name

	_token := jwt.NewWithClaims(jwt.SigningMethodHS256, token)
	return _token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func GenerateRefreshToken(request *entity.Admin) (string, error) {
	var token model.TokenPyload
	duration := os.Getenv("DURATION_JWT_REFRESH_TOKEN")
	lifeTime, _ := strconv.Atoi(duration)

	now := time.Now()
	token.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:   "WebProfile",
		IssuedAt: jwt.NewNumericDate(now),
		// ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour * time.Duration(lifeTime))),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * time.Duration(lifeTime))),
	}

	token.AdminID = request.ID
	token.Username = request.Username
	token.Name = request.Name

	_token := jwt.NewWithClaims(jwt.SigningMethodHS256, token)
	return _token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}
