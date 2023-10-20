package converter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/valyala/fasthttp"
)

var (
	baseURL = "https://openexchangerates.org/api"
	appID   = os.Getenv("OPEN_EXCHANGE_RATES_APP_ID") // Use environment variable for App ID
)

type ExchangeRates struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
}

type ErrorResponse struct {
	Message string `json:"error"`
}

func createRequest(url string) (*fasthttp.Request, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	return req, nil
}

func doRequest(req *fasthttp.Request) (*fasthttp.Response, error) {
	resp := fasthttp.AcquireResponse()
	err := fasthttp.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %v", err)
	}
	return resp, nil
}

func handleResponse(resp *fasthttp.Response, target interface{}) error {
	if resp.StatusCode() != fasthttp.StatusOK {
		var errorResponse ErrorResponse
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			return fmt.Errorf("API error: %s", errorResponse.Message)
		}
		return fmt.Errorf("API error: status code %d", resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), target); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}
	return nil
}

func GetExchangeRates(baseCurrency string) (ExchangeRates, error) {
	url := fmt.Sprintf("%s/latest.json?app_id=%s&base=%s", baseURL, appID, baseCurrency)
	req, err := createRequest(url)
	if err != nil {
		return ExchangeRates{}, err
	}

	resp, err := doRequest(req)
	if err != nil {
		return ExchangeRates{}, err
	}
	defer fasthttp.ReleaseResponse(resp)

	var rates ExchangeRates
	if err := handleResponse(resp, &rates); err != nil {
		return ExchangeRates{}, err
	}

	return rates, nil
}

func GetHistoricalExchangeRates(date, baseCurrency string) (ExchangeRates, error) {
	url := fmt.Sprintf("%s/historical/%s.json?app_id=%s&base=%s", baseURL, date, appID, baseCurrency)
	req, err := createRequest(url)
	if err != nil {
		return ExchangeRates{}, err
	}

	resp, err := doRequest(req)
	if err != nil {
		return ExchangeRates{}, err
	}
	defer fasthttp.ReleaseResponse(resp)

	var rates ExchangeRates
	if err := handleResponse(resp, &rates); err != nil {
		return ExchangeRates{}, err
	}

	return rates, nil
}

func ConvertCurrency(from, to string, amount float64) (float64, error) {
	url := fmt.Sprintf("%s/convert?app_id=%s&from=%s&to=%s&amount=%f", baseURL, appID, from, to, amount)
	req, err := createRequest(url)
	if err != nil {
		return 0, err
	}

	resp, err := doRequest(req)
	if err != nil {
		return 0, err
	}
	defer fasthttp.ReleaseResponse(resp)

	var result map[string]float64
	if err := handleResponse(resp, &result); err != nil {
		return 0, err
	}

	return result["result"], nil
}

func GetUsage() (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/usage.json?app_id=%s", baseURL, appID)
	req, err := createRequest(url)
	if err != nil {
		return nil, err
	}

	resp, err := doRequest(req)
	if err != nil {
		return nil, err
	}
	defer fasthttp.ReleaseResponse(resp)

	var usage map[string]interface{}
	if err := handleResponse(resp, &usage); err != nil {
		return nil, err
	}

	return usage, nil
}

func GetAvailableCurrencies() (map[string]string, error) {
	url := fmt.Sprintf("%s/currencies.json", baseURL)
	req, err := createRequest(url)
	if err != nil {
		return nil, err
	}

	resp, err := doRequest(req)
	if err != nil {
		return nil, err
	}
	defer fasthttp.ReleaseResponse(resp)

	var currencies map[string]string
	if err := handleResponse(resp, &currencies); err != nil {
		return nil, err
	}

	return currencies, nil
}
