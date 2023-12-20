package configs

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"healthcare/models/schema"
	"log"
)

var DB *gorm.DB

func Init() {
	ConnectDB()
	InitialMigration()
}

func ConnectDB() {
	var err error

	configuration := GetConfig()

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		configuration.DB_USERNAME,
		configuration.DB_PASSWORD,
		configuration.DB_HOST,
		configuration.DB_PORT,
		configuration.DB_NAME,
	)

	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Println("Failed to Connect Database")
	}
	log.Println("Connected to Database")
}

func InitialMigration() {
	DB.AutoMigrate(
		&schema.User{},
		&schema.Admin{},
		&schema.Doctor{},
		&schema.Medicine{},
		&schema.Article{},
		&schema.DoctorTransaction{},
		&schema.MedicineTransaction{},
		&schema.MedicineDetails{},
		&schema.Checkout{},
		&schema.Roomchat{},
		&schema.Message{},
	)
}

func ConnectDBTest() *gorm.DB {

	TDB_Username := "root"
	TDB_Password := ""
	TDB_Port := "3306"
	TDB_Host := "localhost"
	TDB_Name := "finaldb"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		TDB_Username, TDB_Password, TDB_Host, TDB_Port, TDB_Name)

	log.Printf("Connection String: %s\n", dsn)

	var errDB error
	DB, errDB = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if errDB != nil {
		panic("Failed to Connect Database")
	}
	return DB
}
