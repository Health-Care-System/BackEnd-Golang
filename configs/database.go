package configs

import (
	"fmt"
	"healthcare/models/schema"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	TDB_Username := os.Getenv("DB_USERNAME")
	TDB_Password := os.Getenv("DB_PASSWORD")
	TDB_Host := os.Getenv("DB_HOST")
	TDB_Port := os.Getenv("DB_PORT")
	TDB_Name := os.Getenv("DB_NAME")

	if TDB_Username == "" || TDB_Password == "" || TDB_Host == "" || TDB_Port == "" || TDB_Name == "" {
		fmt.Println("one or more database environment variables are not set")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		TDB_Username, TDB_Password, TDB_Host, TDB_Port, TDB_Name)

	var errDB error
	DB, errDB = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if errDB != nil {
		log.Fatalf("errornya adalah: %v", errDB)
	}
	return DB
}
