package tonapi

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
)

type RateResp struct {
	Rates Rates `json:"rates"`
}
type Prices struct {
	Usdt float64 `json:"USDT"`
}
type Ton struct {
	Prices Prices `json:"prices"`
}
type Rates struct {
	Ton Ton `json:"TON"`
}

func GetTonRate() (decimal.Decimal, error) {
	resp, err := http.Get("https://tonapi.io/v2/rates?tokens=ton&currencies=usdt")
	if err != nil {
		return decimal.Zero, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return decimal.Zero, err
	}
	var rates RateResp
	err = json.Unmarshal(body, &rates)
	if err != nil {
		return decimal.Zero, err
	}

	return decimal.NewFromFloat(rates.Rates.Ton.Prices.Usdt), nil
}
