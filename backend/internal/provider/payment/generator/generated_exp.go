package generator

import "slices"

var AllCurrency = []Currency{
	CurrencyUsdt,
}

func (e Currency) IsValid() bool {
	return slices.Contains(AllCurrency, e)
}
