package web

type DoctorRegisterResponse struct {
	Fullname       string `json:"fullname" form:"fullname" `
	Email          string `json:"email" form:"email"`
	Status         bool   `json:"status" form:"status"`
	Price          int    `json:"price" form:"price"`
	Gender         string `json:"gender" form:"gender"`
	Specialist     string `json:"specialist" form:"specialist"`
	ProfilePicture string `json:"profile_picture" form:"profile_picture"`
	NoSTR          int    `json:"no_str" form:"no_str"`
	Experience     string `json:"experience" form:"experience"`
	Alumnus        string `json:"alumnus" form:"alumnus"`
}
type DoctorLoginResponse struct {
	Fullname string `json:"fullname" form:"fullname"`
	Email    string `json:"email" form:"email"`
	Token    string `json:"token" form:"token"`
}

type DoctorUpdateResponse struct {
	Fullname         string `json:"fullname" form:"fullname"`
	Email            string `json:"email" form:"email"`
	Gender           string `json:"gender" form:"gender"`
	Specialist       string `json:"specialist" form:"specialist"`
	ProfilePicture   string `json:"profile_picture" form:"profile_picture"`
	NoSTR            int    `json:"no_str" form:"no_str"`
	Experience       string `json:"experience" form:"experience"`
	Alumnus          string `json:"alumnus" form:"alumnus"`
	Status           bool   `json:"status" form:"status"`
	AboutDoctor      string `json:"about_doctor" form:"about_doctor"`
	LocationPractice string `json:"location_practice" form:"location_practice" `
}

type DoctorAllResponse struct {
	ID             uint   `json:"id" form:"id"`
	ProfilePicture string `json:"profile_picture" form:"profile_picture"`
	Fullname       string `json:"fullname" form:"fullname"`
	Specialist     string `json:"specialist" form:"specialist"`
	Price          int    `json:"price" form:"price"`
	Status         bool   `json:"status" form:"status"`
}

type DoctorAllResponseByAdmin struct {
	ID             uint   `json:"id" form:"id"`
	ProfilePicture string `json:"profile_picture" form:"profile_picture"`
	Fullname       string `json:"fullname" form:"fullname"`
	Gender         string `json:"gender" form:"gender"`
	Email          string `json:"email" form:"email"`
	Status         bool   `json:"status" form:"status"`
	Price          int    `json:"price" form:"price"`
	Specialist     string `json:"specialist" form:"specialist"`
	Experience     string `json:"experience" form:"experience"`
	NoSTR          int    `json:"no_str" form:"no_str"`
	Role           string `json:"role" form:"role"`
	Alumnus        string `json:"alumnus" form:"alumnus"`
	// DoctorTransaction []DoctorTransaction `gorm:"ForeignKey:DoctorID;references:ID"`
}

type DoctorIDResponse struct {
  
	ID               uint   `json:"id" form:"id"`
	ProfilePicture   string `json:"profile_picture" form:"profile_picture"`
	Status           bool   `json:"status" form:"status"`
	Fullname         string `json:"fullname" form:"fullname"`
	Specialist       string `json:"specialist" form:"specialist"`
	Price            int    `json:"price" form:"price"`
	Experience       string `json:"experience" form:"experience"`
	AboutDoctor      string `json:"about_doctor" form:"about_doctor"`
	NoSTR            int    `json:"no_str" form:"no_str"`
	LocationPractice string `json:"location_practice" form:"location_practice"`
	Alumnus          string `json:"alumnus" form:"alumnus"`
}


type DoctorProfile struct {
	Fullname       string `json:"fullname" form:"fullname"`
	Email          string `json:"email" form:"email"`
	Gender         string `json:"gender" form:"gender"`
	Specialist     string `json:"specialist" form:"specialist"`
	ProfilePicture string `json:"profile_picture" form:"profile_picture"`
	NoSTR          int    `json:"no_str" form:"no_str"`
	Experience     string `json:"experience" form:"experience"`
	Alumnus        string `json:"alumnus" form:"alumnus"`
	AboutDoctor    string `json:"about_doctor" form:"about_doctor"`
}

type DoctorRoomchat struct {
	ID			   uint   `json:"id" form:"id"`
	Fullname       string `json:"fullname" form:"fullname"`
	Status         bool   `json:"status" form:"status"`
	ProfilePicture string `json:"profile_picture" form:"profile_picture"`
}

// type DoctorPatientResponse struct {
// 	UserID              uint      `json:"user_id" form:"user_id"`
// 	Fullname            string    `json:"fullname" form:"fullname"`
// 	DoctorTransactionID uint      `json:"doctor_transaction_id" form:"doctor_transaction"`
// 	CreatedAt           time.Time `json:"created_at" form:"created_at"`
// 	HealthDetails       string    `json:"health_details" form:"health_details"`
// 	PatientStatus       string    `json:"status" form:"status"`
// }

