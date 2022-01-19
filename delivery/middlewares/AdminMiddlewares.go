package middlewares

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

//func CheckRole(next echo.Group) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		_, role := ExtractTokenUser(c)
//
//		if role == "user" {
//			fmt.Println(role)
//			return c.JSON(
//				http.StatusUnauthorized, map[string]interface{}{
//					"Message": "Unauthorized",
//				},
//			)
//		}
//	}
//}

func CheckRole(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, role := ExtractTokenUser(c)

		if role == "user" {
			fmt.Println(role)
			return c.JSON(
				http.StatusUnauthorized, map[string]interface{}{
					"Message": "Unauthorized",
				},
			)
		}
		return next(c)
	}
}
