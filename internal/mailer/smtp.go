package mailer

import "gopkg.in/gomail.v2"

type smtpMailService struct {
	from   string
	dialer *gomail.Dialer
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewSMTPMailService(cfg SMTPConfig) MailService {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	return smtpMailService{from: cfg.From, dialer: dialer}
}

func (ss smtpMailService) Send(msg MailMessage) error {
	m := gomail.NewMessage()
	m.SetHeader("From", ss.from)
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Subject)
	// m.SetAddressHeader("Bcc", "dan@example.com", "Dan")
	if msg.IsHTML {
		m.SetBody("text/plain", string(msg.Body))
	}
	// m.Attach("txt.md")
	return ss.dialer.DialAndSend(m)
}
