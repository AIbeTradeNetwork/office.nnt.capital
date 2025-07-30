package payment

import (
	"context"

	"github.com/google/uuid"

	"server/internal/domain"
	"server/internal/provider/payment/generator"
)

func (c *Client) OrderCreate(ctx context.Context, order domain.OrderCreateReq) (*domain.CreateOrderResp, error) {
	currency := generator.Currency(order.CurrencyCode)
	if !currency.IsValid() {
		return nil, Error("CurrencyIsValid").SetCode(domain.ErrCurrencyCodeInvalid)
	}

	res, err := generator.OrderCreate(ctx, c.graphqlClient, generator.OrderCreateIn{
		Amount: order.Amount,
		Callbacks: generator.Callbacks{
			FailureURL: order.Callbacks.FailureURL,
			ReturnURL:  order.Callbacks.ReturnURL,
			SuccessURL: order.Callbacks.SuccessURL,
		},
		Comment:      order.Comment,
		Currency:     currency,
		CustomFields: convertToCustomField(order.CustomFields),
		Customer: generator.CustomerIn{
			Account: order.Customer.Account,
			Email:   order.Customer.Email,
			Phone:   order.Customer.Phone,
		},
		ExpirationDateTime: order.ExpirationDateTime,
		PaymentSystemID:    order.PaymentSystemID,
	})
	if err != nil {
		return nil, Error("order_create").SetCode(domain.ErrServer).Add(err)
	}

	return &domain.CreateOrderResp{
		BillId: res.GetOrder_create().BillId,
		URL:    res.GetOrder_create().Url,
	}, nil
}

func (c *Client) OrderCancel(ctx context.Context, id uuid.UUID) (bool, error) {
	res, err := generator.OrderCancel(ctx, c.graphqlClient, id)
	if err != nil {
		return false, Error("order_cancel").SetCode(domain.ErrServer).Add(err)
	}

	return res.GetOrder_cancel(), nil
}

func (c *Client) OrderInfo(ctx context.Context, id uuid.UUID) (*domain.OrderInfo, error) {
	res, err := generator.OrderInfo(ctx, c.graphqlClient, id)
	if err != nil {
		return nil, Error("order_create").SetCode(domain.ErrServer).Add(err)
	}

	oi := res.GetOrder_info()

	currencyCode := domain.CurrencyCode(oi.Currency)
	if !currencyCode.IsValid() {
		return nil, Error("CurrencyCodeIsValid").SetCode(domain.ErrCurrencyCodeInvalid)
	}

	status := domain.OrderStatus(oi.Status)
	if !status.IsValid() {
		return nil, Error("OrderStatusIsValid").SetCode(domain.ErrServer)
	}

	return &domain.OrderInfo{
		Amount:           oi.Amount,
		AmountPaid:       oi.AmountPaid,
		BillID:           oi.BillId,
		ChangedDateTime:  oi.ChangedDateTime,
		Comment:          oi.Comment,
		CreationDateTime: oi.CreationDateTime,
		CurrencyCode:     currencyCode,
		CustomFields:     convertToDomainCustomField(oi.CustomFields),
		Customer: domain.Customer{
			Account: oi.Customer.Account,
			Email:   oi.Customer.Email,
			Phone:   oi.Customer.Phone,
		},
		ExpirationDateTime: oi.ExpirationDateTime,
		MerchantID:         oi.MerchantId,
		PayURL:             oi.PayUrl,
		PaymentSystemID:    oi.PaymentSystemID,
		Status:             status,
	}, nil
}
