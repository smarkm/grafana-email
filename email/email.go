package email

import (
	"log"

	"gopkg.in/gomail.v2"
	"smark.freecoop.net/grafana-email/model"
)

func SendEmail(s *model.Schedule) {
	m := gomail.NewMessage()

	// Set the sender address
	m.SetHeader("From", "sender@example.com")

	// Set the recipient address
	m.SetHeader("To", "recipient@example.com")

	// Set the email subject
	m.SetHeader("Subject", "Hello, Golang Email!")

	// Set the email body (plain text)
	m.SetBody("text/plain", "This is the email body.")

	// Create a new SMTP client
	d := gomail.NewDialer("smtp.example.com", 587, "your_username", "your_password")

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Error sending email:", err)
	}

	log.Println("Email sent successfully.")
}
