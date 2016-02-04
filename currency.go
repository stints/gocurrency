package currency

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/shopspring/decimal"
)

var (
	CurrencyTypes = map[string]Currency{
		"AUD": Currency{code: "AUD", symbol: "$"},
		"BGN": Currency{code: "BGN", symbol: "лв"},
		"BRL": Currency{code: "BRL", symbol: "R$"},
		"CAD": Currency{code: "CAD", symbol: "$"},
		"CHF": Currency{code: "CHF", symbol: "CHF"},
		"CNY": Currency{code: "CNY", symbol: "¥"},
		"CZK": Currency{code: "CZK", symbol: "Kč"},
		"DKK": Currency{code: "DKK", symbol: "kr"},
		"GBP": Currency{code: "GBP", symbol: "£"},
		"HKD": Currency{code: "HKD", symbol: "$"},
		"HRK": Currency{code: "HRK", symbol: "kn"},
		"HUF": Currency{code: "HUF", symbol: "Ft"},
		"IDR": Currency{code: "IDR", symbol: "Rp"},
		"ILS": Currency{code: "ILS", symbol: "₪"},
		"INR": Currency{code: "INR", symbol: "₹"},
		"JPY": Currency{code: "JPY", symbol: "¥"},
		"KRW": Currency{code: "KRW", symbol: "₩"},
		"MXN": Currency{code: "MXN", symbol: "$"},
		"MYR": Currency{code: "MYR", symbol: "RM"},
		"NOK": Currency{code: "NOK", symbol: "kr"},
		"NZD": Currency{code: "NZD", symbol: "$"},
		"PHP": Currency{code: "PHP", symbol: "₱"},
		"PLN": Currency{code: "PLN", symbol: "zł"},
		"RON": Currency{code: "RON", symbol: "lei"},
		"RUB": Currency{code: "RUB", symbol: "руб"},
		"SEK": Currency{code: "SEK", symbol: "kr"},
		"SGD": Currency{code: "SGD", symbol: "$"},
		"THB": Currency{code: "THB", symbol: "฿"},
		"TRY": Currency{code: "TRY", symbol: "₺"},
		"USD": Currency{code: "USD", symbol: "$"},
		"ZAR": Currency{code: "ZAR", symbol: "R"},
	}
	url = "https://api.fixer.io/latest?base=%s&symbols=%s"
)

type Currency struct {
	symbol string
	code   string
}

type MoneyObject struct {
	money    decimal.Decimal
	currency Currency
}

func Money(value interface{}, code string) (MoneyObject, error) {
	currency, found := CurrencyTypes[code]
	if !found {
		return MoneyObject{}, errors.New("Code not found.")
	}

	var money decimal.Decimal
	var moneyObject MoneyObject

	switch v := value.(type) {
	case string:
		m, err := decimal.NewFromString(v)
		if err != nil {
			return MoneyObject{}, err
		}
		money = m
	case float32:
		money = decimal.NewFromFloat(float64(v))
	case float64:
		money = decimal.NewFromFloat(v)
	case int:
		money = decimal.NewFromFloat(float64(v))
	default:
		return MoneyObject{}, errors.New("Value could not be translated.")
	}

	moneyObject.money = money
	moneyObject.currency = currency

	return moneyObject, nil
}

func (m MoneyObject) Convert(code string) (MoneyObject, error) {
	currency, found := CurrencyTypes[code]
	if !found {
		return MoneyObject{}, errors.New("Code not found.")
	}

	resp, err := http.Get(fmt.Sprintf(url, m.currency.code, code))
	if err != nil {
		return MoneyObject{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return MoneyObject{}, err
	}

	var apiObj map[string]*json.RawMessage
	var rates map[string]float64

	err = json.Unmarshal(body, &apiObj)
	if err != nil {
		return MoneyObject{}, err
	}

	err = json.Unmarshal(*apiObj["rates"], &rates)
	if err != nil {
		return MoneyObject{}, err
	}

	conversionRate := rates[code]
	convertedMoney := MoneyObject{
		money:    m.money.Mul(decimal.NewFromFloat(conversionRate)),
		currency: currency,
	}

	return convertedMoney, nil
}

func (m MoneyObject) String() string {
	return fmt.Sprintf("%s%s", m.currency.symbol, m.money)
}
