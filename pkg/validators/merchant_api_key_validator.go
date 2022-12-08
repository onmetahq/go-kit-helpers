package validators

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	ctxLogger "github.com/onmetahq/go-kit-helpers/pkg/logger"
	"github.com/onmetahq/go-kit-helpers/pkg/models"
	metahttp "github.com/onmetahq/meta-http/pkg/meta_http"
	ctxKeys "github.com/onmetahq/meta-http/pkg/models"
)

type Merchant struct {
	ID         string    `json:"id,omitempty"`
	Email      string    `json:"email"`
	OTP        string    `json:"otp"`
	IsVerfied  bool      `json:"isVerified"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt,omitempty"`
	WebhookUrl string    `json:"webhookUrl"`
	APIKey     string    `json:"apiKey"`
	APISecret  string    `json:"apiSecret"`
}

type MerchantAPIResponse struct {
	Success bool     `json:"success"`
	Data    Merchant `json:"data"`
}

type KeyValidator interface {
	ValidateKey(ctx context.Context, apikey string) (Merchant, error)
}

type KeyStore interface {
	Put(key string, mer Merchant) error
	Get(key string) (Merchant, error)
}

type DefaultValidator struct {
	client *metahttp.Client
	store  KeyStore
}

type DefaultStore struct {
	data map[string]Merchant
}

func (store DefaultStore) Put(key string, mer Merchant) error {
	store.data[key] = mer
	return nil
}

func (store DefaultStore) Get(key string) (Merchant, error) {
	if mer, ok := store.data[key]; ok {
		return mer, nil
	}

	return Merchant{}, models.ErrNotFound
}

func NewValidator(client *metahttp.Client) KeyValidator {
	return &DefaultValidator{
		client: client,
		store: DefaultStore{
			data: map[string]Merchant{},
		},
	}
}

func NewValidatorWithStore(client *metahttp.Client, store KeyStore) KeyValidator {
	return &DefaultValidator{
		client: client,
		store:  store,
	}
}

func (svc DefaultValidator) fetchMerchantDetails(ctx context.Context, apikey string) (Merchant, error) {
	var res MerchantAPIResponse
	req := map[string]string{"apikey": apikey}
	_, err := svc.client.Post(ctx, "", map[string]string{}, req, &res)

	if err != nil {
		return Merchant{}, err
	}

	return res.Data, nil
}

func (svc DefaultValidator) ValidateKey(ctx context.Context, apikey string) (Merchant, error) {
	mer, err := svc.store.Get(apikey)

	if err != nil && err == models.ErrNotFound {
		mer, err = svc.fetchMerchantDetails(ctx, apikey)

		if err != nil {
			return Merchant{}, err
		}

		svc.store.Put(apikey, mer)
		return mer, nil
	}

	if err != nil {
		return Merchant{}, err
	}

	return mer, nil
}

func MerchantAPIKeyValidator(svc KeyValidator, logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			lg := ctxLogger.NewCtxLogger(logger)

			apikey, ok := ctx.Value(ctxKeys.MerchantAPIKey).(string)
			if !ok {
				lg.Context(ctx).Error().Log("msg", "Invalid Merchant API key")
				return nil, models.ErrUnauthorized
			}

			mer, err := svc.ValidateKey(ctx, apikey)
			if err != nil || mer.ID == "" {
				lg.Context(ctx).Error().Log("msg", "Merchant API key does not exist")
				return nil, models.ErrUnauthorized
			}

			ctx = context.WithValue(ctx, ctxKeys.TenantID, mer.ID)
			return next(ctx, request)
		}
	}
}
