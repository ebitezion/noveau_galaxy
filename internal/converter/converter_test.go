package converter

import (
	"testing"
)

func TestGetExchangeRates(t *testing.T) {

	appID = "YOUR_APP_ID"

	rates, err := GetExchangeRates("USD")

	if err != nil {
		t.Errorf("Error fetching exchange rates: %v", err)
	}

	if rates.Base != "USD" {
		t.Errorf("Expected base currency to be USD, got %s", rates.Base)
	}
}

func TestGetHistoricalExchangeRates(t *testing.T) {

	appID = "YOUR_APP_ID"

	rates, err := GetHistoricalExchangeRates("2023-10-18", "USD")

	if err != nil {
		t.Errorf("Error fetching historical exchange rates: %v", err)
	}

	if rates.Base != "USD" {
		t.Errorf("Expected base currency to be USD, got %s", rates.Base)
	}
}

func TestConvertCurrency(t *testing.T) {

	appID = "YOUR_APP_ID"

	amount := 100.0
	result, err := ConvertCurrency("USD", "EUR", amount)

	if err != nil {
		t.Errorf("Error converting currency: %v", err)
	}

	if result <= 0 {
		t.Errorf("Invalid conversion result: %f", result)
	}
}

func TestGetUsage(t *testing.T) {

	appID = "YOUR_APP_ID"

	usage, err := GetUsage()

	if err != nil {
		t.Errorf("Error fetching usage information: %v", err)
	}

	if len(usage) == 0 {
		t.Errorf("Empty usage information")
	}
}

func TestGetAvailableCurrencies(t *testing.T) {

	appID = "YOUR_APP_ID"

	currencies, err := GetAvailableCurrencies()

	if err != nil {
		t.Errorf("Error fetching available currencies: %v", err)
	}

	if len(currencies) == 0 {
		t.Errorf("Empty available currencies")
	}
}
