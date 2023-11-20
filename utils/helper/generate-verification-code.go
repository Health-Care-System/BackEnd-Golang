package helper

import (
	"crypto/rand"
	"math/big"
)

// GenerateVerificationCode generates a random alphanumeric verification code.
func GenerateVerificationCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	codeLength := 8
	randomBytes := make([]byte, codeLength)
	for i := range randomBytes {
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		randomBytes[i] = charset[randIndex.Int64()]
	}

	return string(randomBytes)
}

// // GetDoctorOTPByEmail mengambil OTP dari database berdasarkan email dokter
// func GetDoctorOTPByEmail(email string) (string, error) {
// 	var doctor schema.Doctor
// 	if err := configs.DB.Where("email = ?", email).First(&doctor).Error; err != nil {
// 		return "", err
// 	}

// 	return doctor.OTP, nil
// }

// // VerifyDoctorOTP memverifikasi OTP untuk email yang diberikan dan dokter
// func VerifyDoctorOTP(email, enteredOTP string) bool {
// 	storedOTP, err := GetDoctorOTPByEmail(email)
// 	if err != nil {
// 		return false
// 	}
// 	return enteredOTP == storedOTP
// }
