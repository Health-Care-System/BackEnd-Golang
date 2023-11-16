package routes

import (
	"healthcare/controllers"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	avilableDoctor := e.Group("/patients")
	avilableDoctor.GET("/doctors/available", controllers.GetAvailableDoctor)
	avilableDoctor.GET("/doctors", controllers.GetSpecializeDoctor)
}
