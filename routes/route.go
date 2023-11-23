package routes

import (
	"healthcare/controllers"
	"healthcare/middlewares"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {

	AdminJWT := middlewares.AdminRoleAuth
	UserJWT := middlewares.UserIDRoleAuth
	DoctorJWT := middlewares.DoctorIDRoleAuth

	gAdmins := e.Group("/admins")
	gAdmins.POST("/login", controllers.LoginAdminController)
	gAdmins.POST("/register/doctor", controllers.RegisterDoctorByAdminController, AdminJWT)
	gAdmins.GET("/list/doctors", controllers.GetAllDoctorByAdminController, AdminJWT)
	gAdmins.PUT("/update/doctor/:id", controllers.UpdateDoctorByAdminController, AdminJWT)
	gAdmins.PUT("/update/payment/:id", controllers.UpdatePaymentStatusByAdminController, AdminJWT)
	gAdmins.DELETE("/delete/doctor/:id", controllers.DeleteDoctorByAdminController, AdminJWT)
	gAdmins.PUT("/:id", controllers.UpdateAdminController, AdminJWT)
	gAdmins.POST("/medicines", controllers.CreateMedicineController)
	gAdmins.GET("/medicines", controllers.GetAllMedicinesAdminController)
	gAdmins.GET("/medicines/:id", controllers.GetMedicineController)
	gAdmins.GET("/medicines/?name=", controllers.GetMedicineByNameAdminController)
	gAdmins.PUT("/medicines/:id", controllers.UpdateMedicineController)
	gAdmins.DELETE("/medicines/:id", controllers.DeleteMedicineController)
	gAdmins.GET("/medicines/:id/image", controllers.GetImageMedicineController)
	gAdmins.PUT("/medicines/:id/image", controllers.UpdateImageMedicineController)
	gAdmins.DELETE("/medicines/:id/image", controllers.DeleteImageMedicineController)

	gUsers := e.Group("/users")
	gUsers.POST("/register", controllers.RegisterUserController)
	gUsers.POST("/login", controllers.LoginUserController)
	gUsers.GET("/profile", controllers.GetUserController, UserJWT)
	gUsers.PUT("/profile", controllers.UpdateUserController, UserJWT)
	gUsers.DELETE("", controllers.DeleteUserController, UserJWT)
	gUsers.GET("/medicines", controllers.GetAllMedicinesUserController)
	gUsers.GET("/medicines/?name=", controllers.GetMedicineByNameUserController)
	gUsers.GET("/doctors/available", controllers.GetAvailableDoctor)
	gUsers.GET("/doctors", controllers.GetSpecializeDoctor)

	gDoctors := e.Group("/doctors")
	gDoctors.POST("/login", controllers.LoginDoctorController)
	gDoctors.GET("/profile", controllers.GetDoctorProfileController, DoctorJWT)
	gDoctors.GET("/:id", controllers.GetDoctorByIDController)
	gDoctors.PUT("/profile", controllers.UpdateDoctorController, DoctorJWT)
	gDoctors.DELETE("", controllers.DeleteDoctorController, DoctorJWT)
	gDoctors.GET("", controllers.GetAllDoctorController)
	// gDoctors.GET("manage/patitient", controllers.GetAllPatientsController, DoctorJWT)
	// gDoctors.GET("manage/patitient/:status", controllers.GetPatientsByStatusController, DoctorJWT)
	// gDoctors.PUT("manage/patitient/:idTransaksi", controllers.UpdatePatientController, DoctorJWT)
}
