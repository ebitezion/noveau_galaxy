package main

import (
	"net/http"
	"strconv"

	"github.com/ebitezion/backend-framework/internal/converter"
)

func (app *application) ExchangeRatesHandler(w http.ResponseWriter, r *http.Request) {
	baseCurrency := r.URL.Query().Get("base")

	if baseCurrency == "" {
		http.Error(w, "Base currency is required", http.StatusBadRequest)
		return
	}

	rates, err := converter.GetExchangeRates(baseCurrency)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	data := envelope{
		"responseCode": "00",
		"status":       "success",
		"rates":        rates,
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) HistoricalExchangeRatesHandler(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	baseCurrency := r.URL.Query().Get("base")

	if date == "" || baseCurrency == "" {
		http.Error(w, "Date and base currency are required", http.StatusBadRequest)
		return
	}

	rates, err := converter.GetHistoricalExchangeRates(date, baseCurrency)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"responseCode":              "00",
		"status":                    "success",
		"historical_exchange_rates": rates,
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) ConvertCurrencyHandler(w http.ResponseWriter, r *http.Request) {

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	amountStr := r.URL.Query().Get("amount")

	if from == "" || to == "" || amountStr == "" {
		http.Error(w, "From currency, to currency, and amount are required", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	result, err := converter.ConvertCurrency(from, to, amount)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"responseCode": "00",
		"status":       "success",
		"result":       result,
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) UsageHandler(w http.ResponseWriter, r *http.Request) {
	usage, err := converter.GetUsage()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	data := envelope{
		"responseCode": "00",
		"status":       "success",
		"usage":        usage,
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) AvailableCurrenciesHandler(w http.ResponseWriter, r *http.Request) {
	currencies, err := converter.GetAvailableCurrencies()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"responseCode": "00",
		"status":       "success",
		"currencies":   currencies,
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
