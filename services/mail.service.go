package services

import (
	"log"

	"gopkg.in/gomail.v2"
)

type MailService struct{}

func NewMailService() *MailService {
	return &MailService{}
}

const (
	CONFIG_SMTP_HOST     = "smtp.gmail.com"
	CONFIG_SMTP_PORT     = 587
	CONFIG_SENDER_NAME   = "mikhael <mikhaeljhon22@gmail.com>"
	CONFIG_AUTH_EMAIL    = "mikhaeljhon22@gmail.com"
	CONFIG_AUTH_PASSWORD = "silv cfrk ebjz nquy"
)

// Mengirim satu email
func (s *MailService) SendEmail(to, subject, body string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Println("Mail sent to", to)
	return nil
}

// Struct job email
type EmailJob struct {
	To      string
	Subject string
	Body    string
}

func StartEmailWorkerPool(mailService *MailService, numWorkers int) chan EmailJob {
	jobs := make(chan EmailJob, 100)
	for w := 1; w <= numWorkers; w++ {
		go func(id int) {
			for job := range jobs {
				log.Printf("Worker %d sending email to %s\n", id, job.To)
				if err := mailService.SendEmail(job.To, job.Subject, job.Body); err != nil {
					log.Printf("Failed to send email to %s: %v\n", job.To, err)
				}
			}
		}(w)
	}
	return jobs
}
