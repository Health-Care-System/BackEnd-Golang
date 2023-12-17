package helper

import (
	"fmt"
	"healthcare/configs"
	"healthcare/models/schema"
	"time"

	"gorm.io/gorm"
)

// UpdatePasswordInDatabase
func UpdatePasswordInDatabase(db *gorm.DB, tableName, email, hashedPassword, otp string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Update the password based on email and OTP
	err := db.Table(tableName).Where("email = ? AND otp = ?", email, otp).Update("password", hashedPassword).Error
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserVerificationStatus(email string, isVerified bool) error {
	var user schema.User
	if err := configs.DB.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error; err != nil {
		return err
	}

	// Update user verification status without modifying 'gender'
	if err := configs.DB.Model(&user).Select("is_verified").Updates(map[string]interface{}{"is_verified": isVerified}).Error; err != nil {
		return err
	}

	return nil
}

func UpdatePasswordAndMarkVerified(db *gorm.DB, tableName, email, hashedPassword, otp string) error {
	return db.Table(tableName).
		Where("email = ? AND otp = ?", email, otp).
		Updates(map[string]interface{}{
			"password":    hashedPassword,
			"is_verified": true, 
			"updated_at":  time.Now(),
		}).Error
}
