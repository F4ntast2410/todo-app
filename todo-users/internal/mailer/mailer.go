package mailer

import (
	"crypto/tls"
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
		"Ваш код подтверждения: %s\r\n\r\nКод действителен 5 минут.",
		code,
	)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s\r\n",
		m.cfg.From, to, subject, body,
	))

	addr := fmt.Sprintf("%s:%s", m.cfg.Host, m.cfg.Port)

	// Для порта 465 используем tls.Dial вместо smtp.SendMail
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName: m.cfg.Host,
	})
	if err != nil {
		return fmt.Errorf("tls dial error: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, m.cfg.Host)
	if err != nil {
		return fmt.Errorf("smtp client error: %w", err)
	}
	defer client.Quit()

	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("auth error: %w", err)
	}

	if err := client.Mail(m.cfg.From); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	return w.Close()
}
