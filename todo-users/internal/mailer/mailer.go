package mailer

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"time"
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
	subject := "Код подтверждения"
	body := fmt.Sprintf(
		"Ваш код подтверждения: %s\r\n\r\nКод действителен 5 минут. Если вы не запрашивали вход — просто игнорируйте это письмо.",
		code,
	)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s\r\n",
		m.cfg.From, to, subject, body,
	))

	addr := fmt.Sprintf("%s:%s", m.cfg.Host, m.cfg.Port)

	// Если используем порт 465 (Implicit TLS)
	if m.cfg.Port == "465" {
		tlsconfig := &tls.Config{
			ServerName: m.cfg.Host,
		}

		// 🛠️ Используем DialWithDialer, принудительно задаем "tcp4" и таймаут 10 секунд
		dialer := &net.Dialer{
			Timeout: 10 * time.Second,
		}

		conn, err := tls.DialWithDialer(dialer, "tcp4", addr, tlsconfig)
		if err != nil {
			return fmt.Errorf("tls dial failed: %w", err)
		}
		defer conn.Close()

		// ... остальной код отправки без изменений
		client, err := smtp.NewClient(conn, m.cfg.Host)
		if err != nil {
			return fmt.Errorf("smtp client creation failed: %w", err)
		}
		defer client.Quit()

		auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("auth failed: %w", err)
		}

		if err = client.Mail(m.cfg.From); err != nil {
			return err
		}
		if err = client.Rcpt(to); err != nil {
			return err
		}

		w, err := client.Data()
		if err != nil {
			return err
		}
		if _, err = w.Write(msg); err != nil {
			return err
		}
		return w.Close()
	}

	// Для порта 587 (STARTTLS)
	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	return smtp.SendMail(addr, auth, m.cfg.From, []string{to}, msg)
}
