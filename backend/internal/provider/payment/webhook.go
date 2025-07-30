package payment

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"

	"server/internal/domain"
)

func (c *Client) Handle(w http.ResponseWriter, r *http.Request, process func(ctx context.Context, order *domain.OrderIn) error) error {
	orderIn, err := getOrderFromRequest(r)
	if err != nil {
		http.Error(w, "getOrderFromRequest", http.StatusBadRequest)

		return Error("getOrderFromRequest").SetCode(domain.ErrBadRequest).Add(err)
	}

	fmt.Println(r.Header.Get("X-Api-Signature-Sha256"), c.config, orderIn)

	err = testOrderSecurityCode(
		r.Header.Get("X-Api-Signature-Sha256"),
		c.config.Key,
		orderIn,
	)
	if err != nil {
		http.Error(w, "testOrderSecurityCode", http.StatusBadRequest)

		return Error("testOrderSecurityCode").SetCode(domain.ErrBadRequest).Add(err)
	}

	err = process(r.Context(), orderIn)
	if err != nil {
		http.Error(w, "process", http.StatusBadRequest)

		return Error("process").SetCode(domain.ErrBadRequest).Add(err)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{ "status": "ok" }`))

	return nil
}

func getOrderFromRequest(r *http.Request) (*domain.OrderIn, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, Error("Content-Type").SetCode(domain.ErrBadRequest)
	}

	var res domain.OrderIn

	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return nil, Error("NewDecoder").SetCode(domain.ErrBadRequest).Add(err)
	}

	return &res, nil
}

func testOrderSecurityCode(
	securityTokenFromUser string,
	merchantPassword string,
	order *domain.OrderIn,
) error {
	mac := hmac.New(sha256.New, []byte(merchantPassword))

	_, err := mac.Write([]byte(fmt.Sprintf(
		"%s|%s|%s",
		uuid.MustParse(order.BillID),
		order.CurrencyCode,
		order.Status,
	)))
	if err != nil {
		return Error("mac.Write").SetCode(domain.ErrBadRequest)
	}

	securityTokenFromUS := hex.EncodeToString(mac.Sum(nil))

	if securityTokenFromUser != securityTokenFromUS {
		return Error("Check securityTokenFromUser").SetCode(domain.ErrAccessDenied)
	}

	return nil
}
