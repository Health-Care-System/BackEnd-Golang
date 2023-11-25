package middlewares

import (
	"healthcare/utils/helper"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AdminRoleAuth(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		token, err := ExtractToken(c)
		if err != nil {
			return err
		}

		if token.Valid {
			userClaims := token.Claims.(jwt.MapClaims)
			userRole := userClaims["role"].(string)
			userID := int(userClaims["id"].(float64))

			if userRole == "admin" {
				c.Set("userID", userID)
				return next(c)
			} else {
				return c.JSON(http.StatusForbidden, helper.ErrorResponse("You are not Authorized to Access this Resource"))
			}
		} else {
			return c.JSON(http.StatusUnauthorized, helper.ErrorResponse("Invalid or Expired Token"))
		}
	}
}
