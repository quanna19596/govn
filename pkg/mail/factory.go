package mail

import (
	"fmt"
	"shopify/internal/utils"
)

type ProviderType string

const (
	ProviderMailtrap ProviderType = "mailtrap"
)

type MailtrapProviderFactory struct{}

func (f *MailtrapProviderFactory) CreateProvider(config *MailConfig) (EmailProviderService, error) {
	return NewMailtrapProvider(config)
}

func NewProviderFactory(providerType ProviderType) (ProviderFactory, error) {
	switch providerType {
	case ProviderMailtrap:
		return &MailtrapProviderFactory{}, nil
	default:
		return nil, utils.NewError(fmt.Sprintf("Unsupported provider type: %s", utils.ErrorCode(providerType)), utils.ErrCodeInternalServer)
	}

}
