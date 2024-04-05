package model

import "errors"

type CurrencyType int32

const (
	NONE_CURRENCY CurrencyType = iota
	TWD
	USD
	YEN
)

var (
	currencyMapper = make(map[CurrencyType]*Currency)
)

func init() {
	currencyMapper = map[CurrencyType]*Currency{
		NONE_CURRENCY: {
			CurrencyType: NONE_CURRENCY,
			CurrencyName: "無",
		},
		TWD: {
			CurrencyType: TWD,
			CurrencyName: "台幣",
		},
		USD: {
			CurrencyType: USD,
			CurrencyName: "美元",
		},
	}
}

type Currency struct {
	CurrencyType CurrencyType `json:"currencyType"`
	CurrencyName string       `json:"currencyName"`
}

func GetAllCurrencyTypes() []*Currency {

	currencies := []*Currency{}

	for _, v := range currencyMapper {
		currencies = append(currencies, v)
	}
	return currencies
}

func GetCurrencyType(currencyType int32) (*Currency, error) {

	currency, ok := currencyMapper[CurrencyType(currencyType)]

	if !ok {
		return nil, errors.New("Cannot find currency type")
	}
	return currency, nil
}
