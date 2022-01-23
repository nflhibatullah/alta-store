package middlewares

import (
	"altastore/constant"
	"altastore/delivery/common"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func CreateToken(userId int, role string, email string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["role"] = role
	claims["email"] = email
	claims["userId"] = int(userId)
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(constant.JWT_SECRET_KEY))
}

func ExtractTokenUser(e echo.Context) (common.JWTPayload, error) {
	user := e.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		role := claims["role"].(string)
		email := claims["email"].(string)
		userId := int(claims["userId"].(float64))
		return common.JWTPayload{
			ID: userId,
			Email: email,
			Role: role,
		}, nil
	}
	return common.JWTPayload{}, errors.New("invalid token")
}
