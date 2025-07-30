package payment

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"server/internal/domain"
)

func TestNew(t *testing.T) {
	t.SkipNow()

	ctx := context.Background()

	cl := New(Config{
		URL:   "https://pay.aibetrade.com/graphql",
		Login: "",
		Key:   "",
	})

	paymentSystemID, _ := uuid.Parse("eabf1c97-3c27-47a9-b56c-55edfbe84a22")

	r, err := cl.OrderCreate(ctx, domain.OrderCreateReq{
		Amount: 1,
		Callbacks: domain.Callbacks{
			FailureURL: "2",
			ReturnURL:  "3",
			SuccessURL: "4",
		},
		Comment:      "5",
		CurrencyCode: "USDT",
		CustomFields: nil,
		Customer: domain.Customer{
			Account: "7",
			Email:   "8",
			Phone:   "9",
		},
		ExpirationDateTime: time.Now(),
		PaymentSystemID:    paymentSystemID,
	})
	//r, err := cl.OrderInfo(ctx, uuid.New())
	if err != nil {
		t.Error("OrderInfo", err)
		return
	}

	t.Log(r)
}
