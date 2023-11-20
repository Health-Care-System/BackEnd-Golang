package configs

import (
	"fmt"
	"healthcare/models/schema"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
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
	)
}
