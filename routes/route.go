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
	gAdmins.GET("/profile", controllers.GetAdminProfileController, AdminJWT)
	gAdmins.PUT("/profile", controllers.UpdateAdminController, AdminJWT)
	gAdmins.POST("/doctors/register", controllers.RegisterDoctorByAdminController, AdminJWT)
	gAdmins.GET("/doctors", controllers.GetAllDoctorByAdminController, AdminJWT)
	gAdmins.GET("/users", controllers.GetAllUserByAdminController, AdminJWT)
	gAdmins.GET("/user/:user_id", controllers.GetUserIDbyAdminController, AdminJWT)
	gAdmins.DELETE("/user/:user_id", controllers.DeleteUserByAdminController, AdminJWT)
	gAdmins.GET("/doctor/:doctor_id", controllers.GetDoctorIDbyAdminController, AdminJWT)
	gAdmins.PUT("/doctor/:doctor_id", controllers.UpdateDoctorByAdminController, AdminJWT)
	gAdmins.DELETE("/doctor/:doctor_id", controllers.DeleteDoctorByAdminController, AdminJWT)
	gAdmins.PUT("/doctor-payments/:transaction_id", controllers.UpdatePaymentStatusByAdminController, AdminJWT)
	gAdmins.GET("/doctor-payment/:user_id", controllers.GetUserPaymentsByAdminsController, AdminJWT)
	gAdmins.GET("/doctor-payments", controllers.GetAllDoctorsPaymentsByAdminsController, AdminJWT)
	gAdmins.GET("/doctor-payment", controllers.GetDoctorTransactionByIDController, AdminJWT)
	gAdmins.POST("/medicines", controllers.CreateMedicineController, AdminJWT)
	gAdmins.GET("/medicines", controllers.GetMedicineAdminController, AdminJWT)
	gAdmins.GET("/medicines/:medicine_id", controllers.GetMedicineAdminByIDController, AdminJWT)
	gAdmins.PUT("/medicines/:medicine_id", controllers.UpdateMedicineController, AdminJWT)
	gAdmins.DELETE("/medicines/:medicine_id", controllers.DeleteMedicineController, AdminJWT)
	gAdmins.GET("/medicines/:medicine_id/image", controllers.GetImageMedicineController, AdminJWT)
	gAdmins.PUT("/medicines/:medicine_id/image", controllers.UpdateImageMedicineController, AdminJWT)
	gAdmins.DELETE("/medicines/:medicine_id/image", controllers.DeleteImageMedicineController, AdminJWT)
	gAdmins.PUT("/medicines-payments/checkout/:checkout_id", controllers.UpdateCheckoutController, AdminJWT)
	gAdmins.GET("/medicines-payments/checkout", controllers.GetAdminCheckoutController, AdminJWT)
	gAdmins.GET("/medicines-payments/checkout/:checkout_id", controllers.GetAdminCheckoutByIDController, AdminJWT)

	gUsers := e.Group("/users")
	gUsers.POST("/register", controllers.RegisterUserController)
	gUsers.POST("/login", controllers.LoginUserController)
	gUsers.GET("/profile", controllers.GetUserController, UserJWT)
	gUsers.PUT("/profile", controllers.UpdateUserController, UserJWT)
	gUsers.DELETE("", controllers.DeleteUserController, UserJWT)
	gUsers.GET("/medicines", controllers.GetMedicineUserController)
	gUsers.GET("/medicines/:medicine_id", controllers.GetMedicineUserByIDController)
	gUsers.GET("/doctors/available", controllers.GetAvailableDoctor)
	gUsers.GET("/doctors", controllers.GetSpecializeDoctor)
	gUsers.GET("/articles", controllers.GetAllArticles)
	gUsers.GET("/articles/:article_id", controllers.GetArticleByID)
	gUsers.GET("/article", controllers.GetAllArticlesByTitle)
	gUsers.POST("/doctor-payments/:doctor_id", controllers.CreateDoctorTransactionController, UserJWT)
	gUsers.GET("/doctor-payments", controllers.GetAllDoctorTransactionsController, UserJWT)
	gUsers.GET("/doctor-payments/:transaction_id", controllers.GetDoctorTransactionController, UserJWT)
	gUsers.POST("/chats/:transaction_id", controllers.CreateRoomchatController, UserJWT)
	gUsers.GET("/chats/:roomchat_id", controllers.GetUserRoomchatController, UserJWT)
	gUsers.POST("/chats/:roomchat_id/message", controllers.CreateComplaintMessageController, UserJWT)
	gUsers.POST("/medicines-payments", controllers.CreateMedicineTransaction, UserJWT)
	gUsers.GET("/medicines-payments", controllers.GetMedicineTransactionController, UserJWT)
	gUsers.GET("/medicines-payments/:medtrans_id", controllers.GetMedicineTransactionByIDController, UserJWT)
	gUsers.DELETE("/medicines-payments/:medtrans_id", controllers.DeleteMedicineTransactionController, UserJWT)
	gUsers.POST("/medicines-payments/checkout", controllers.CreateCheckoutController, UserJWT)
	gUsers.GET("/medicines-payments/checkout", controllers.GetUserCheckoutController, UserJWT)
	gUsers.GET("/medicines-payments/checkout/:checkout_id", controllers.GetUserCheckoutByIDController, UserJWT)
	gUsers.POST("/get-otp", controllers.GetOTPForPasswordReset)
	gUsers.POST("/verify-otp", controllers.VerifyOTP)
	gUsers.POST("/change-password", controllers.ResetPassword)

	gDoctors := e.Group("/doctors")
	gDoctors.POST("/login", controllers.LoginDoctorController)
	gDoctors.GET("/profile", controllers.GetDoctorProfileController, DoctorJWT)
	gDoctors.GET("/:doctor_id", controllers.GetDoctorByIDController)
	gDoctors.PUT("/profile", controllers.UpdateDoctorController, DoctorJWT)
	gDoctors.PUT("/status", controllers.ChangeDoctorStatusController, DoctorJWT)
	gDoctors.DELETE("", controllers.DeleteDoctorController, DoctorJWT)
	gDoctors.GET("/articles", controllers.DoctorGetAllArticles, DoctorJWT)
	gDoctors.GET("/articles/:article_id", controllers.DoctorGetArticleByID, DoctorJWT)
	gDoctors.POST("/articles", controllers.CreateArticle, DoctorJWT)
	gDoctors.PUT("/articles/:article_id", controllers.UpdateArticleById, DoctorJWT)
	gDoctors.DELETE("/articles/:article_id", controllers.DeleteArticleById, DoctorJWT)
	gDoctors.GET("/doctor-consultation", controllers.GetAllDoctorConsultationController, DoctorJWT)
	gDoctors.GET("/chats", controllers.GetAllDoctorRoomchatController, DoctorJWT)
	gDoctors.GET("/chats/:roomchat_id", controllers.GetDoctorRoomchatController, DoctorJWT)
	gDoctors.POST("/chats/:roomchat_id/message", controllers.CreateAdviceMessageController, DoctorJWT)
	gDoctors.GET("/manage-user", controllers.GetManageUserController, DoctorJWT)
	gDoctors.PUT("/manage-user/:transaction_id", controllers.UpdateManageUserController, DoctorJWT)
	gDoctors.POST("/get-otp", controllers.GetOTPForPasswordReset)
	gDoctors.POST("/verify-otp", controllers.VerifyOTP)
	gDoctors.POST("/change-password", controllers.ResetDoctorPassword)

	e.POST("/chatbot", controllers.Chatbot)
	e.POST("/customerservice", controllers.CustomerService)

}
