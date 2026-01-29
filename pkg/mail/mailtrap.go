package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"shopify/internal/utils"
	"shopify/pkg/logger"
	"time"

	"github.com/rs/zerolog"
)

type MailtrapConfig struct {
	MailSender string
	NameSender string
	MailUrl    string
	MailApiKey string
}

type MailtrapProvider struct {
	client *http.Client
	config *MailtrapConfig
	logger *zerolog.Logger
}

func NewMailtrapProvider(config *MailConfig) (EmailProviderService, error) {
	mailtrapCfg, ok := config.ProviderConfig["mailConfig"].(map[string]any)
	if !ok {
		return nil, utils.NewError("Invalid or missing Mail configuaration", utils.ErrCodeInternalServer)
	}

	return &MailtrapProvider{
		client: &http.Client{Timeout: config.Timeout},
		config: &MailtrapConfig{
			MailSender: mailtrapCfg["mail_sender"].(string),
			NameSender: mailtrapCfg["name_sender"].(string),
			MailUrl:    mailtrapCfg["mail_url"].(string),
			MailApiKey: mailtrapCfg["mail_api_key"].(string),
		},
		logger: config.Logger,
	}, nil
}
func (p *MailtrapProvider) SendEmail(ctx context.Context, email *Email) error {
	traceID := logger.GetTraceID(ctx)
	start := time.Now()

	time.Sleep(5 * time.Second)

	email.From = Address{
		Email: p.config.MailSender,
		Name:  p.config.NameSender,
	}

	payload, err := json.Marshal(email)
	if err != nil {
		return utils.WrapError(err, "Failed to marshal email", utils.ErrCodeInternalServer)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.config.MailUrl, bytes.NewReader(payload))
	if err != nil {
		return utils.WrapError(err, "Failed to create request", utils.ErrCodeInternalServer)
	}

	req.Header.Add("Authorization", "Bearer "+p.config.MailApiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		p.logger.Error().Str("trace_id", traceID).
			Dur("duration", time.Since(start)).
			Str("operation", "send_mail").
			Err(err).
			Msg("Failed to send request")
		return utils.WrapError(err, "Failed to send request", utils.ErrCodeInternalServer)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		p.logger.Error().Str("trace_id", traceID).
			Dur("duration", time.Since(start)).
			Str("operation", "send_mail").
			Int("status_code", resp.StatusCode).
			Str("response_body", string(body)).
			Msg("Unexpected response from mailtrap")

		return utils.NewError(fmt.Sprintf("Unexpected response from mailtrap with code %d: %s", resp.StatusCode, string(body)), utils.ErrCodeInternalServer)
	}
	return nil
}
