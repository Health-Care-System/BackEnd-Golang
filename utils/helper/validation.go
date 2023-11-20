package helper

import "time"

func GenderIsValid(gender string) bool {
	if gender == "female" || gender == "male" {
		return true
	}
	return false
}

func BloodTypeIsValid(bloodType string) bool {
	if bloodType == "A" || bloodType == "B" || bloodType == "AB" || bloodType == "O" {
		return true
	}
	return false
}

func BirthdateIsValid(birthdate string) bool {
	if birthdate != "" {
		_, err := time.Parse("2006-01-02", birthdate)
		return err == nil
	}
	return false
}