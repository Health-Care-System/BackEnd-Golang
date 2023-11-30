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
	gAdmins.PUT("/profile", controllers.UpdateAdminController, AdminJWT)
	gAdmins.GET("/profile", controllers.GetAdminProfileController, AdminJWT)
	gAdmins.POST("/register/doctor", controllers.RegisterDoctorByAdminController, AdminJWT)
	gAdmins.GET("/list/doctors", controllers.GetAllDoctorByAdminController, AdminJWT)
	gAdmins.GET("/users", controllers.GetAllUserByAdminController, AdminJWT)
	gAdmins.GET("/user/:id", controllers.GetUserIDbyAdminController, AdminJWT)
	gAdmins.DELETE("/user/:id", controllers.DeleteUserByAdminController, AdminJWT)
	gAdmins.PUT("/update/doctor/:id", controllers.UpdateDoctorByAdminController, AdminJWT)
	gAdmins.DELETE("/delete/doctor/:id", controllers.DeleteDoctorByAdminController, AdminJWT)
	gAdmins.PUT("/update/payment/:id", controllers.UpdatePaymentStatusByAdminController, AdminJWT)
	gAdmins.POST("/medicines", controllers.CreateMedicineController, AdminJWT)
	gAdmins.GET("/medicines", controllers.GetMedicineAdminController)
	gAdmins.GET("/medicines/:id", controllers.GetMedicineAdminByIDController)
	// gAdmins.GET("/medicine", controllers.GetMedicineByNameAdminController)
	gAdmins.PUT("/medicines/:id", controllers.UpdateMedicineController, AdminJWT)
	gAdmins.DELETE("/medicines/:id", controllers.DeleteMedicineController, AdminJWT)
	gAdmins.GET("/medicines/:id/image", controllers.GetImageMedicineController)
	gAdmins.PUT("/medicines/:id/image", controllers.UpdateImageMedicineController, AdminJWT)
	gAdmins.DELETE("/medicines/:id/image", controllers.DeleteImageMedicineController, AdminJWT)

	gUsers := e.Group("/users")
	gUsers.POST("/register", controllers.RegisterUserController)
	gUsers.POST("/login", controllers.LoginUserController)
	gUsers.GET("/profile", controllers.GetUserController, UserJWT)
	gUsers.PUT("/profile", controllers.UpdateUserController, UserJWT)
	gUsers.DELETE("", controllers.DeleteUserController, UserJWT)
	gUsers.GET("/medicines", controllers.GetMedicineUserController)
	gUsers.GET("/medicines/:id", controllers.GetMedicineUserByIDController)
	gUsers.GET("/doctors/available", controllers.GetAvailableDoctor)
	gUsers.GET("/doctors", controllers.GetSpecializeDoctor)
	gUsers.GET("/articles", controllers.GetAllArticles)
	gUsers.GET("/articles/:id", controllers.GetArticleByID)
	gUsers.GET("/article", controllers.GetAllArticlesByTitle)
	gUsers.POST("/doctor-payments/:doctor_id", controllers.CreateDoctorTransactionController, UserJWT)
	gUsers.GET("/doctor-payments", controllers.GetAllDoctorTransactionsController, UserJWT)
	gUsers.GET("/doctor-payments/:transaction_id", controllers.GetDoctorTransactionController, UserJWT)
	gUsers.POST("/chats/:transaction_id", controllers.CreateRoomchatController, UserJWT)
	gUsers.GET("/chats/:roomchat_id", controllers.GetUserRoomchatController, UserJWT)
	gUsers.POST("/chats/:roomchat_id/message", controllers.CreateComplaintMessageController, UserJWT)

	gDoctors := e.Group("/doctors")
	gDoctors.POST("/login", controllers.LoginDoctorController)
	gDoctors.GET("/profile", controllers.GetDoctorProfileController, DoctorJWT)
	gDoctors.GET("/:idDokter", controllers.GetDoctorByIDController)
	gDoctors.PUT("/profile", controllers.UpdateDoctorController, DoctorJWT)
	gDoctors.DELETE("", controllers.DeleteDoctorController, DoctorJWT)
	gDoctors.GET("", controllers.GetAllDoctorController)
	gDoctors.GET("/articles", controllers.DoctorGetAllArticles, DoctorJWT)
	gDoctors.GET("/articles/:id", controllers.DoctorGetArticleByID, DoctorJWT)
	gDoctors.POST("/articles", controllers.CreateArticle, DoctorJWT)
	gDoctors.PUT("/articles/:id", controllers.UpdateArticleById, DoctorJWT)
	gDoctors.DELETE("/articles/:id", controllers.DeleteArticleById, DoctorJWT)
	gDoctors.GET("/doctor-consultation", controllers.GetAllDoctorConsutationController, DoctorJWT)
	gDoctors.GET("/chats", controllers.GetAllDoctorRoomchatController, DoctorJWT)
	gDoctors.GET("/chats/:roomchat_id", controllers.GetDoctorRoomchatController, DoctorJWT)
	gDoctors.POST("/chats/:roomchat_id/message", controllers.CreateAdviceMessageController, DoctorJWT)
	gDoctors.GET("/manage-user", controllers.GetManagePatientController, DoctorJWT)
	gDoctors.PUT("/manage-user", controllers.UpdateManagePatientController, DoctorJWT)

	e.POST("/chatbot", controllers.Chatbot)
	e.POST("/customerservice", controllers.CustomerService)

}
