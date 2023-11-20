package middlewares

import (
	"healthcare/utils/helper"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// authentication and authorization USER
const role = "user"

// check "role" only
func UserRoleAuth(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		token, err := ExtractToken(c)
		if err != nil {
			return err
		}

		if token.Valid {
			userClaims := token.Claims.(jwt.MapClaims)
			userRole := userClaims["role"].(string)

			if userRole == role {
				return next(c)
			} else {
				return c.JSON(http.StatusForbidden, helper.ErrorResponse("You are not Authorized to Access this Resource"))
			}
		} else {
			return c.JSON(http.StatusUnauthorized, helper.ErrorResponse("Invalid or Expired Token"))
		}
	}
}

// check "role" and "id"
func UserIDRoleAuth(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		token, err := ExtractToken(c)
		if err != nil {
			return err
		}

		if token.Valid {
			userClaims := token.Claims.(jwt.MapClaims)
			userRole := userClaims["role"].(string)
			userID := int(userClaims["id"].(float64))

			if userRole == role {
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
