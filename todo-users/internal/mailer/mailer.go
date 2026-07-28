package mailer

import (
	"context"
	"encoding/base64"
	"fmt"
	"mime"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
	From         string
}

type GmailAPIMailer struct {
	cfg Config
}

func New(cfg Config) *GmailAPIMailer {
	return &GmailAPIMailer{cfg: cfg}
}

func (m *GmailAPIMailer) Send2FACode(to string, code string) error {
	ctx := context.Background()

	// Настройка OAuth2 конфига
	oauthConfig := &oauth2.Config{
		ClientID:     m.cfg.ClientID,
		ClientSecret: m.cfg.ClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{gmail.GmailSendScope},
	}

	token := &oauth2.Token{
		RefreshToken: m.cfg.RefreshToken,
	}

	client := oauthConfig.Client(ctx, token)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to create gmail client: %w", err)
	}

	// Формирование письма в формате RFC 822
	subject := mime.BEncoding.Encode("UTF-8", "Код подтверждения")
	body := fmt.Sprintf("Ваш код подтверждения: %s\n\nКод действителен 5 минут.", code)

	rawMessage := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		m.cfg.From, to, subject, body,
	)

	// Gmail API требует URL-safe Base64 кодирование сообщения
	encodedMessage := base64.URLEncoding.EncodeToString([]byte(rawMessage))

	message := &gmail.Message{
		Raw: encodedMessage,
	}

	_, err = srv.Users.Messages.Send("me", message).Do()
	if err != nil {
		return fmt.Errorf("failed to send message via Gmail API: %w", err)
	}

	return nil
}
