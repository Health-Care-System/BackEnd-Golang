package middlewares

import (
	"healthcare/utils/helper"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// create token
func GenerateToken(userID uint, email string, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // 3 hari
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ExtractToken(c echo.Context) (*jwt.Token, error) {
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return nil, c.JSON(http.StatusUnauthorized, helper.ErrorResponse("Missing Token"))
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, c.JSON(http.StatusUnauthorized, helper.ErrorResponse("Invalid or Expired Token"))
	}

	return token, nil
}
