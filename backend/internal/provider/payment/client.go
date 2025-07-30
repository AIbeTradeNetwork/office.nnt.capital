package payment

import (
	"net/http"
	"server/internal/domain"
	"server/internal/provider/payment/generator"

	"github.com/Khan/genqlient/graphql"
)

const (
	errorSource = "[provider.payment]"
)

type Config struct {
	URL   string
	Login string
	Key   string
}

type Client struct {
	graphqlClient graphql.Client
	config        Config
}

func New(c Config) *Client {
	httpClient := http.Client{
		Transport: &authedTransport{
			app:     c.Login,
			key:     c.Key,
			wrapped: http.DefaultTransport,
		},
	}

	graphqlClient := graphql.NewClient(c.URL, &httpClient)

	return &Client{
		graphqlClient: graphqlClient,
		config:        c,
	}
}

type authedTransport struct {
	app string
	key string

	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Auth-App", t.app)
	req.Header.Set("X-Auth-Key", t.key)

	return t.wrapped.RoundTrip(req)
}

func Error(source string) *domain.Error {
	return domain.NewError(errorSource + "[" + source + "]")
}

func convertToCustomField(req []domain.CustomField) []generator.CustomFieldIn {
	var res []generator.CustomFieldIn

	for _, v := range req {
		res = append(res, generator.CustomFieldIn{
			Name:  v.Name,
			Value: v.Value,
		})
	}

	return res
}

func convertToDomainCustomField(req []generator.OrderInfoOrder_infoOrderCustomFieldsCustomField) []domain.CustomField {
	var res []domain.CustomField

	for _, v := range req {
		res = append(res, domain.CustomField{
			Name:  v.Name,
			Value: v.Value,
		})
	}

	return res
}
