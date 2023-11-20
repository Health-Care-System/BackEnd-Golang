package helper

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, content string) error {
	// Konfigurasi server email
	smtpServer := os.Getenv("SMTPSERVER")
	smtpPortStr := os.Getenv("SMTPPORT")
	smtpUsername := os.Getenv("SMTPUSERNAME")
	smtpPassword := os.Getenv("SMTPPASSWORD")

	// Periksa apakah semua variabel lingkungan
	if smtpServer == "" || smtpPortStr == "" || smtpUsername == "" || smtpPassword == "" {
		return errors.New("Konfigurasi SMTP belum lengkap. Mohon atur semua variabel lingkungan yang diperlukan.")
	}

	// Konversi smtpPortStr ke int
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Println("Gagal mengonversi port SMTP ke integer:", err)
		return err
	}

	// Buat pesan email baru
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUsername)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	// Tambahkan bagian HTML ke body email
	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				/* CSS */
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
					padding: 20px;
				}
				.container {
					background-color: #ffffff;
					padding: 20px;
					border-radius: 10px;
				}
				h1 {
					color: #007bff;
				}
				p {
					color: #333;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>%s</h1>
				<p>%s</p>
			</div>
		</body>
		</html>
	`, subject, content)

	m.SetBody("text/html", htmlBody)

	// Konfigurasi dialer untuk server email
	d := gomail.NewDialer(smtpServer, smtpPort, smtpUsername, smtpPassword)

	// Kirim email
	if err := d.DialAndSend(m); err != nil {
		log.Println("Gagal mengirim email:", err)
		return err
	}

	log.Printf("Email berhasil dikirim ke %s\n", to)

	return nil
}

// SendNotificationEmail mengirimkan email notifikasi
func SendNotificationEmail(to, fullname, notificationType, userType string) error {
	go func() {
		var subject, body string

		switch notificationType {
		case "login":
			subject = "Login Notification"
			body = "Hello, "

			// Customize the greeting based on the user type
			if userType == "drg" {
				body += "Drg. "
			}

			body += fullname + "! Anda telah berhasil masuk."
		case "register":
			subject = "Registration Notification"
			body = "Hello, "

			// Customize the greeting based on the user type
			if userType == "drg" {
				body += "Drg. "
			}

			body += fullname + "! Anda telah berhasil mendaftar."
		default:
			err := errors.New("Jenis notifikasi tidak valid")
			log.Println(err)
			return
		}

		err := SendEmail(to, subject, body)
		if err != nil {
			log.Println("Gagal mengirim email:", err)
		}
	}()

	return nil
}
