package middlewares

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"time"
)

func CreateToken(userId int, SecretKey string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["role"] = "user"
	claims["userId"] = int(userId)
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}

func ExtractTokenUser(e echo.Context) (int, string) {
	user := e.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		role := claims["role"].(string)
		userId := int(claims["userId"].(float64))
		return userId, role
	}
	return 0, "Token tidak valid"
}
