package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ebitezion/backend-framework/internal/accounts"
	"github.com/ebitezion/backend-framework/internal/notifications"
	"github.com/ebitezion/backend-framework/internal/validator"
	"github.com/julienschmidt/httprouter"
	"github.com/jung-kurt/gofpdf"
)

// Retrieve the "id" URL parameter from the current request context, then convert
// it to an integer and return it. If the operation isn't successful, return 0
// and an error.
func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// Define an envelope type.
type envelope map[string]interface{}

// Define a writeJSON() helper for sending responses. This takes the destination
// http.ResponseWriter, the HTTP status code to send, the data to encode to JSON,
// and a header map containing any additional HTTP headers we want to include in
// the response.
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// Encode the data to JSON, returning the error if there was one.
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications
	js = append(js, '\n')

	// Add any headers that we want to include. We loop through the header map and
	// add each header to the http.ResponseWriter header map. Note that it's OK if
	// the provided header map is nil. Go doesn't throw an error if you try to range
	// over (or generally, read from) a nil map.
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Add the "Content-Type: application/json" header, then write the status code
	// and JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Use the http.MaxBytesReaeder() to limit the size of the request body to 1MB.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Initialize the json.Decoder, and call the DisallowUnknownFields() method on it
	// before decoding. This means that if the JSON from the client now includes any
	// field which cannot be mapped to the target destination, the decoder will return
	// an error instead of just ignoring the field.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Use the Decode() method to decode the body contents into the dst interface.
	// Notice that when we call Decode() we pass a *pointer* to the dst interface as
	// the target decode destination. If there was an error during decoding, we also
	// use the generic errorResponse() helper to send the client a 400 Bad Request
	// response containing the error message.
	err := dec.Decode(dst)
	if err != nil {
		// If there is an error during decoding, start the triage...
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError. If it doesn, then return a plain-english error message
		// which includes the location of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
		// for syntax errors in the JSON. So we check for this using errors.Is() and
		// return a generic error message.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// Likewise, catch any *json.UnmarshalTypeError errors. These occur when the
		// JSON value is the wrong type for the target destination. If the error relates
		// to a specific field, then we include that in our error message to make it
		// easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", unmarshalTypeError.Offset)

		// An io.EOF error will be returned by Decode() if the request body is empty. We
		// check for this with errors.Is() and return a plain-english error message instead.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// If the JSON contains a field which cannot be mapped to the target destination
		// then Decode() will now return an error message in the format "json: unknown
		// field "<name>"". We check for this, extract the field name from the error,
		// and interpolate it into our custom error message.
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// If the request body exceeds 1MB in size the decode will now fail with the
		// error "http: request body too large".
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		// A json.InvalidUnmarshalError error will be returned if we pass a non-nil
		// pointer to Decode(). We catch this and panic, rather than returning an error
		// to our handler.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		// For anything else, return the error message as-is.
		default:
			return err
		}
	}

	// Call Decode() again, using a pointer to an empty anonymous struct as the
	// destination. If the request body only contained a single JSON value this will
	// return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body and we return our own custom error message.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// The readString() helper returns a string value from the query string, or the provided
// default value if no matching key could be found.
func (app *application) readString(qs url.Values, key, defaultValue string) string {
	// Extract the value for a given key from the query string. If no key exists this
	// will return the empty string "".
	s := qs.Get(key)

	// If no key exists (or the value is empty) then return the default value
	if s == "" {
		return defaultValue
	}

	// Otherwise return the string
	return s
}

// The readCSV() helper reads a string value from the query string and then splits it
// into a slice on the comma character. If no matching key can be found, it returns
// the provided default value.
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	// Extract the value from the query string.
	csv := qs.Get(key)

	// If no key exists (or the value is empty) then return the default value
	if csv == "" {
		return defaultValue
	}

	// Otherwise parse the value into a []string slice and return it.
	return strings.Split(csv, ",")
}

// The readInt() helper reads a string value from the query string and converts it to an
// integer before returning. If no matching key can be found it returns the provided
// default value. If the value couldn't be converted to an integer, then we record an
// error message in the provided Validator instance.
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	// Extract the value from the query string.
	s := qs.Get(key)

	// If no key exists (or the value is empty) then return the default value
	if s == "" {
		return defaultValue
	}

	// Try to conver the value to an int. If this fails, add an error message to the
	// validator instance and return the default value.
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	// Otherwise return the converted integer value.
	return i
}

// generateRandomNumber gives a random number of a given length
func (app *application) generateRandomNumber(length int) (int, error) {
	if length < 1 {
		return 0, fmt.Errorf("length should be at least 1")
	}

	// Calculate the minimum and maximum values for the specified length
	min := int(math.Pow10(length - 1))
	max := int(math.Pow10(length)) - 1

	if min >= max {
		return 0, fmt.Errorf("invalid length")
	}

	// Initialize the random number generator with a seed based on the current time
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between min and max (inclusive)
	return rand.Intn(max-min+1) + min, nil
}

// generateRandomString makes a random character of a given length
func (app *application) generateRandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%^&*()")

	rand.Seed(time.Now().UnixNano())

	randomString := make([]rune, length)
	for i := range randomString {
		randomString[i] = letters[rand.Intn(len(letters))]
	}
	return string(randomString)
}

// // RenderTemplate renders an HTML template with the provided data and layout.
// // It allows you to specify a template function map for custom functions.
// func (app *application) RenderTemplate(w http.ResponseWriter, templateFiles []string, data interface{}, layout string, funcMap template.FuncMap) {
// 	if funcMap != nil {
// 		tmpl, err := app.TemplateFunction(templateFiles, funcMap, layout)
// 		if err != nil {
// 			return
// 		}
// 		app.ExecuteTemplate(w, tmpl, data, layout)
// 		return
// 	}
// 	tmpl := app.ParseTemplate(templateFiles)
// 	app.ExecuteTemplate(w, tmpl, data, layout)
// }

// // TemplateFunction creates and returns an HTML template with custom functions.
// func (app *application) TemplateFunction(templateFiles []string, funcMap template.FuncMap, layout string) (*template.Template, error) {
// 	tmpl, err := template.New(layout).Funcs(funcMap).ParseFiles(templateFiles...)
// 	if err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}
// 	return tmpl, nil
// }

// // ParseTemplate parses the HTML template files and returns a template.
// func (app *application) ParseTemplate(templateFiles []string) *template.Template {
// 	tmpl, err := template.ParseFiles(templateFiles...)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return tmpl
// }

// // ExecuteTemplate executes the provided HTML template with the given data.
// func (app *application) ExecuteTemplate(w http.ResponseWriter, tmpl *template.Template, data interface{}, layout string) {
// 	layout = app.RemoveSlashFromPath(layout)
// 	err := tmpl.ExecuteTemplate(w, layout, data)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// // RemoveSlashFromPath removes the slashes from a path and returns the last part.
// func (app *application) RemoveSlashFromPath(path string) string {
// 	NewPath := strings.Split(path, "/")
// 	return NewPath[len(NewPath)-1]
// }

func (app *application) RenderTemplate(w http.ResponseWriter, templateFiles []string, data interface{}, layout string, funcMap template.FuncMap) {
	var tmpl *template.Template

	app.mu.Lock()
	defer app.mu.Unlock()

	if layout != "" {
		tmpl = app.getTemplate(layout)
	} else {
		tmpl = app.ParseTemplate(templateFiles)
	}

	if funcMap != nil {
		tmpl.Funcs(funcMap)
	}

	app.ExecuteTemplate(w, tmpl, data, layout)
}

func (app *application) getTemplate(name string) *template.Template {
	if tmpl, ok := app.templates[name]; ok {
		return tmpl
	}

	tmpl, err := template.ParseFiles(name)
	if err != nil {
		log.Println(err)
		return nil
	}

	app.templates[name] = tmpl
	return tmpl
}

func (app *application) ParseTemplate(templateFiles []string) *template.Template {
	app.mu.Lock()
	defer app.mu.Unlock()

	tmpl, err := template.ParseFiles(templateFiles...)
	if err != nil {
		log.Println(err)
	}
	return tmpl
}

func (app *application) ExecuteTemplate(w http.ResponseWriter, tmpl *template.Template, data interface{}, layout string) {
	if tmpl == nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	layout = app.RemoveSlashFromPath(layout)
	err := tmpl.ExecuteTemplate(w, layout, data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (app *application) RemoveSlashFromPath(path string) string {
	NewPath := strings.Split(path, "/")
	return NewPath[len(NewPath)-1]
}

func newApplication() *application {
	return &application{
		templates: make(map[string]*template.Template),
	}
}

func createPdf(data interface{}) (string, error) {

	// Unmarshal JSON data
	jsonData, err := json.Marshal(data)
	if err != nil {

		return "", fmt.Errorf("error marshalling JSON:", err)

	}

	var transactions []accounts.Transaction
	err = json.Unmarshal(jsonData, &transactions)
	if err != nil {

		return "", fmt.Errorf("error marshalling JSON:", err)
	}
	// Create PDF
	pdf := gofpdf.New(gofpdf.OrientationPortrait, gofpdf.UnitPoint, gofpdf.PageSizeA4, "")
	pdf.SetHeaderFunc(func() {
		pdf.SetFont("Arial", "B", 16)
		pdf.Cell(0, 10, "Nouveau Mobile Account Statement")
		pdf.Ln(40)
	})

	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "B", 12)

	// Add content to PDF with spacing between key and value
	// Add content to PDF with spacing between key and value
	for _, t := range transactions {
		pdf.Cell(0, 10, fmt.Sprintf("Transaction ID: %d", t.ID))
		pdf.Ln(16)

		pdf.Cell(0, 10, fmt.Sprintf("Transaction Type: %s", t.Transaction))
		pdf.Ln(16)

		pdf.Cell(0, 10, fmt.Sprintf("Transaction Amount: $%.2f", t.TransactionAmount))
		pdf.Ln(16)

		pdf.Cell(0, 10, fmt.Sprintf("Sender Account: %s, Bank: %s", t.SenderAccountNumber, t.SenderBankNumber))
		pdf.Ln(16)

		pdf.Cell(0, 10, fmt.Sprintf("Receiver Account: %s, Bank: %s", t.ReceiverAccountNumber, t.ReceiverBankNumber))
		pdf.Ln(16)

		pdf.Cell(0, 10, fmt.Sprintf("Fee Amount: $%.2f", t.FeeAmount))
		pdf.Ln(16)

		pdf.Cell(0, 10, fmt.Sprintf("Timestamp: %s", t.Timestamp))
		pdf.Ln(34) // Increased spacing between transactions
	}

	pdfPath := "cmd/web/static/files/Transactions.pdf"
	// Save the PDF to a file or do something else with it
	err = pdf.OutputFileAndClose(pdfPath)
	if err != nil {

		return "", fmt.Errorf("error saving PDF:", err)
	}

	fmt.Println("PDF created successfully!")
	return pdfPath, nil
}
func createExcelSheet(data interface{}) (string, error) {

	//Unmarshal JSON data
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON:", err)
	}

	var transactions []accounts.Transaction
	err = json.Unmarshal(jsonData, &transactions)
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON:", err)
	}
	f := excelize.NewFile()

	// Create a new sheet.
	index := f.NewSheet("TransactionSheet")

	// Set column headers
	headers := []string{"ID", "Transaction", "Type", "SenderAccountNumber", "SenderBankNumber", "ReceiverAccountNumber", "ReceiverBankNumber", "TransactionAmount", "FeeAmount", "Timestamp"}
	for col, header := range headers {
		cell := fmt.Sprintf("%c%d", 'A'+col, 1)
		f.SetCellValue("TransactionSheet", cell, header)
	}

	// Populate data from the Transaction struct array to the spreadsheet
	for row, transaction := range transactions {
		// Use reflection to get the field values dynamically
		for col := 0; col < len(headers); col++ {
			cell := fmt.Sprintf("%c%d", 'A'+col, row+2)
			f.SetCellValue("TransactionSheet", cell, getFieldValue(transaction, headers[col]))
		}
	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	excelPath := "cmd/web/static/files/Transactions.xlsx"
	// Save spreadsheet by the given path.
	if err := f.SaveAs(excelPath); err != nil {
		return "", fmt.Errorf("error saving excel:", err)
	}
	fmt.Println("Excel Sheet created successfully")
	return excelPath, nil
}

// Helper function to get field value using reflection
func getFieldValue(transaction accounts.Transaction, field string) interface{} {
	r := reflect.ValueOf(transaction)
	f := reflect.Indirect(r).FieldByName(field)
	return f.Interface()
}

// this function is how to use the notification package
func (app *application) Notification(token string, sendersAccountNumber string, receiversAccountNumber string, amount string) error {
	sender, err := accounts.FetchAccountMeta(sendersAccountNumber)
	if err != nil {
		return err
	}

	receiver, err := accounts.FetchAccountMeta(receiversAccountNumber)
	if err != nil {
		return err
	}

	users := [2]*accounts.AccountHolderDetails{sender, receiver}

	for i := range users {

		ns := notifications.NotificationService{}
		User := notifications.User{
			ID:       1,
			Username: "adeoluwa",
			Email:    "akanbiadenugba699@gmail.com",
			Phone:    users[i].ContactNumber1,
		}

		notification := notifications.Notification{
			User:    User,
			Message: fmt.Sprintf("Amount of %s was transfered from %s to %s", amount, sender.AccountNumber, receiver.AccountNumber),
		}
		notifications.SendNotification(ns, notification)
	}

	return nil
}

// this function is how to use the notification package
func (app *application) VerificationNotification(token int, email string) error {

	ns := notifications.NotificationService{}
	User := notifications.User{
		ID:       1,
		Username: "adeoluwa",
		Email:    "akanbiadenugba699@gmail.com",
		Phone:    "08088974888",
	}
	fmt.Println("hit")
	notification := notifications.Notification{
		User:    User,
		Message: fmt.Sprintf("Your verication token is %d", token),
	}
	notifications.SendNotification(ns, notification)

	return nil
}
