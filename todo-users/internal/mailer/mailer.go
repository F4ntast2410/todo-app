package mailer

import (
	"fmt"
	"mime"
	"net/smtp"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type SMTPMailer struct {
	cfg Config
}

func New(cfg Config) *SMTPMailer {
	return &SMTPMailer{cfg: cfg}
}

func (m *SMTPMailer) Send2FACode(to string, code string) error {
	subject := mime.BEncoding.Encode("UTF-8", "Код подтверждения")
	body := fmt.Sprintf(
		"Ваш код подтверждения: %s\r\n\r\nКод действителен 5 минут. Если вы не запрашивали вход — просто игнорируйте это письмо.",
		code,
	)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s\r\n",
		m.cfg.From, to, subject, body,
	))

	addr := fmt.Sprintf("%s:%s", m.cfg.Host, m.cfg.Port)

	// Для порта 587 (STARTTLS)
	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	return smtp.SendMail(addr, auth, m.cfg.From, []string{to}, msg)
}
