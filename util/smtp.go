package util

import (
	"bytes"
	"github.com/BATUCHKA/real-estate-back/database"
	"github.com/BATUCHKA/real-estate-back/database/models"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strings"
)

func SMTPSend(to string, subject string, body string) {
	from := os.Getenv("SMTP_EMAIL")
	pass := os.Getenv("SMTP_PASSWORD")
	msg := `From: OdinTech <no-reply@odintech.mn>` + "\r\n" +
		`To: ` + strings.TrimSpace(to) + "\r\n" +
		`Subject: ` + strings.TrimSpace(subject) + "\r\n" +
		`MIME-version: 1.0;` + "\r\n" +
		`Content-Type: text/html; charset="UTF-8";` + "\r\n\r\n" +
		body

		// log.Println(auth, to, msg)
	host := "email-smtp.ap-southeast-1.amazonaws.com"
	auth := smtp.PlainAuth("", from, pass, host)

	err := smtp.SendMail(host+":587", auth, "no-reply@odintech.mn", []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}

func SendBaseSMTP(toEmail string, subject string, smtpText string) {
	if t, err := template.ParseFiles("template/base.html", "template/organization.html"); err != nil {
		log.Println("hi it's here ", err.Error())
	} else {
		ConfigKeyVal.Load()
		templateData := struct {
			SmtpText string
			Email    string
		}{
			SmtpText: smtpText,
			Email:    toEmail,
		}
		buf := new(bytes.Buffer)
		if err = t.ExecuteTemplate(buf, "base", templateData); err != nil {
			log.Println("hi it's here 2 ", err.Error())
		}
		body := buf.String()
		SMTPSend(toEmail, subject, body)

		// db := database.Database
		// newNotf := models.Notifications{
		// 	ToEmail:  toEmail,
		// 	SmtpText: smtpText,
		// 	// IsSeen:   false,
		// }
		// if result := db.GormDB.Create(&newNotf); result.Error != nil {
		// 	log.Println("notf save error: ", err)
		// 	return
		// }
	}
}

// func SendOrgReqResSmtp(toEmail string, otpCode string) {
// 	if t, err := template.ParseFiles("template/base.html", "template/forgot_password_email.html"); err != nil {
// 		log.Println(err.Error())
// 	} else {
// 		ConfigKeyVal.Load()
// 		templateData := struct {
// 			TitleText string
// 			OtpCode   string
// 		}{
// 			TitleText: "Нууц үг сэргээх!",
// 			OtpCode:   otpCode,
// 		}
// 		buf := new(bytes.Buffer)
// 		if err = t.ExecuteTemplate(buf, "base", templateData); err != nil {
// 			log.Println(err.Error())
// 		}
// 		body := buf.String()
// 		SMTPSend(toEmail, "ESAN нууц үг сэргээх", body)
// 	}
// }
