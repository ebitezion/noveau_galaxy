package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ebitezion/backend-framework/internal/accounts"
	cashpickup "github.com/ebitezion/backend-framework/internal/cash_pickup"
)

type CashPickupPageData struct {
	Data      []cashpickup.CashPickup
	AdminName string
}
type BalanceEnquiryPageData struct {
	Data *accounts.BalanceEnquiry
	name string
}
type AllAccountPageData struct {
	Accounts  []accounts.AccountDetails
	AdminName string
}
type AllTransactionsPageData struct {
	Transactions []accounts.Transaction
	AdminName    string
}
type PageData struct {
	AdminName string
}

// RenderIndexPage renders a HTML page
func (app *application) RenderIndexPage(w http.ResponseWriter, r *http.Request) {

	FullName, err := GetAdminName(r)
	if err != nil {
		log.Println(err)
	}
	pageData := AllTransactionsPageData{
		AdminName: FullName,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/index.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/index.html", nil)
}

// RenderLoginPage renders a HTML page
func (app *application) RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/login.html"}, nil, "cmd/web/views/login.html", nil)
}

// RenderSignUpPage renders a HTML page
func (app *application) RenderSignUpPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/signup.html"}, nil, "cmd/web/views/signup.html", nil)
}

// RenderCreateAccountPage renders a HTML page
func (app *application) RenderCreateAccountPage(w http.ResponseWriter, r *http.Request) {

	FullName, err := GetAdminName(r)
	if err != nil {
		log.Println(err)
	}
	pageData := AllTransactionsPageData{
		AdminName: FullName,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/createAccount.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/createAccount.html", nil)
}

// RenderDepositInitiationPage renders a HTML page
func (app *application) RenderDepositInitiationPage(w http.ResponseWriter, r *http.Request) {

	FullName, err := GetAdminName(r)
	if err != nil {
		log.Println(err)
	}
	pageData := AllTransactionsPageData{
		AdminName: FullName,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/deposit.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/deposit.html", nil)
}

// RenderCreditInitiationPage renders a HTML page
func (app *application) RenderCreditInitiationPage(w http.ResponseWriter, r *http.Request) {

	FullName, err := GetAdminName(r)
	if err != nil {
		log.Println(err)
	}
	pageData := AllTransactionsPageData{
		AdminName: FullName,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/credit.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/credit.html", nil)
}

// RenderBatchTransactionPage renders a HTML page
func (app *application) RenderBatchTransactionPage(w http.ResponseWriter, r *http.Request) {

	FullName, err := GetAdminName(r)
	if err != nil {
		log.Println(err)
	}
	pageData := AllTransactionsPageData{
		AdminName: FullName,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/batch_transaction.html", "cmd/web/views/layout.html"}, pageData, "cmd/web/views/layout.html", nil)
}

func (app *application) RenderBalanceEnquiry(w http.ResponseWriter, r *http.Request) {

	FullName, err := GetAdminName(r)
	if err != nil {
		log.Println(err)
	}
	pageData := AllTransactionsPageData{
		AdminName: FullName,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/balanceEnquiry.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/balanceEnquiry.html", nil)
}

func (app *application) RenderAccountHistory(w http.ResponseWriter, r *http.Request) {

	FullName, err := GetAdminName(r)
	if err != nil {
		log.Println(err)
	}
	pageData := AllTransactionsPageData{
		AdminName: FullName,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/accountHistory.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/accountHistory.html", nil)
}

func (app *application) RenderBusinessesPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/businesses.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/businesses.html", nil)
}
func (app *application) RenderPartnersPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/partners.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/partners.html", nil)
}
func (app *application) RenderKycPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/kyc.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/kyc.html", nil)
}
func (app *application) RenderCurrencyConversionPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/currencyConverter.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/currencyConverter.html", nil)
}
func (app *application) RenderTeamsPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/teams.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/teams.html", nil)
}
func (app *application) RenderRolesPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/roles.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/roles.html", nil)
}
func (app *application) RenderSystemLogsPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/systemLogs.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/systemLogs.html", nil)
}
func (app *application) RenderBlockAccountPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/blockAccountPage.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/blockAccountPage.html", nil)
}
func (app *application) RenderUnblockAccountPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/unblockAccountPage.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/unblockAccountPage.html", nil)
}
func (app *application) RenderUkAccountPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/ukAccountPage.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/ukAccountPage.html", nil)
}
func (app *application) RenderUsAccountPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/usAccountPage.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/usAccountPage.html", nil)
}
func (app *application) RenderCashPickupPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/cashPickup.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/cashPickup.html", nil)
}

func (app *application) RenderAllCashPickupPage(w http.ResponseWriter, r *http.Request) {

	response, err := cashpickup.GetAllCashPickups()
	if err != nil {
		fmt.Println(err)
		return
	}
	FullName, err := GetAdminName(r)
	if err != nil {
		log.Println(err)
	}
	pageData := CashPickupPageData{
		Data:      response,
		AdminName: FullName,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/allCashPickup.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/allCashPickup.html", nil)
}

func (app *application) RenderApproveCashPickupPage(w http.ResponseWriter, r *http.Request) {

	response, err := cashpickup.GetAllCashPickups()
	if err != nil {
		fmt.Println(err)
		return
	}
	FullName, err := GetAdminName(r)
	if err != nil {
		log.Println(err)
	}
	pageData := CashPickupPageData{
		Data:      response,
		AdminName: FullName,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/cashPickupApproval.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/cashPickupApproval.html", nil)
}
func (app *application) RenderApproveWithdrawalPage(w http.ResponseWriter, r *http.Request) {

	//get all transactions
	accountNumber := os.Getenv("WITHDRAWAL_ACCOUNT_NUMBER")
	fmt.Println(accountNumber)
	data, err := accounts.ProcessAccount([]string{"", "acmt", "1012", accountNumber})

	if err != nil {
		fmt.Println(err)
	}
	Transactions, ok := data.([]accounts.Transaction)
	if !ok {
		fmt.Println("Failed to convert to []AccountDetails")
		return
	}

	FullName, err := GetAdminName(r)
	if err != nil {
		log.Println(err)
	}

	pageData := AllTransactionsPageData{
		Transactions: Transactions,
		AdminName:    FullName,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/withdrawalApprovals.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/withdrawalApprovals.html", nil)
}
func (app *application) RenderApproveTransferPage(w http.ResponseWriter, r *http.Request) {

	response, err := cashpickup.GetAllCashPickups()
	if err != nil {
		fmt.Println(err)
		return
	}
	FullName, err := GetAdminName(r)
	if err != nil {
		log.Println(err)
	}
	pageData := CashPickupPageData{
		Data:      response,
		AdminName: FullName,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/cashPickupApproval.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/cashPickupApproval.html", nil)
}
func (app *application) RenderApproveDepositPage(w http.ResponseWriter, r *http.Request) {

	accountNumber := os.Getenv("DEPOSIT_ACCOUNT_NUMBER")
	fmt.Println(accountNumber)
	data, err := accounts.ProcessAccount([]string{"", "acmt", "1004", accountNumber})

	if err != nil {
		fmt.Println(err)
	}

	Transactions, ok := data.([]accounts.Transaction)
	if !ok {
		fmt.Println("Failed to convert to []AccountDetails")
		return
	}

	FullName, err := GetAdminName(r)
	if err != nil {
		log.Println(err)
	}
	pageData := AllTransactionsPageData{
		Transactions: Transactions,
		AdminName:    FullName,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/depositApprovals.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/depositApprovals.html", nil)
}
func (app *application) RenderUserCashPickupPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/cashPickup.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/cashPickup.html", nil)
}
func (app *application) RenderWithdrawalPage(w http.ResponseWriter, r *http.Request) {

	FullName, err := GetAdminName(r)
	if err != nil {
		log.Println(err)
	}
	pageData := AllTransactionsPageData{
		AdminName: FullName,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/withdrawalPage.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/withdrawalPage.html", nil)
}

func (app *application) RenderTransactionsPage(w http.ResponseWriter, r *http.Request) {
	//get all transactions
	data, err := accounts.ProcessAccount([]string{"", "acmt", "1008"})
	if err != nil {
		fmt.Println(err)
	}
	Transactions, ok := data.([]accounts.Transaction)
	if !ok {
		fmt.Println("Failed to convert to []AccountDetails")
		return
	}
	//update excel sheet
	_, err = createExcelSheet(Transactions)
	if err != nil {
		fmt.Println(err, "error creating excel sheet")
		return
	}
	//update pdf file
	_, err = createPdf(Transactions)
	if err != nil {
		fmt.Println(err, "error creating pdf sheet")
		return
	}

	pageData := AllTransactionsPageData{
		Transactions: Transactions,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/allTransactions.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/allTransactions.html", nil)
}

func (app *application) RenderInflowPage(w http.ResponseWriter, r *http.Request) {
	//get all transactions

	accountNumber := os.Getenv("DEPOSIT_ACCOUNT_NUMBER")
	fmt.Println(accountNumber)
	data, err := accounts.ProcessAccount([]string{"", "acmt", "1004", accountNumber})

	if err != nil {
		fmt.Println(err)
	}

	Transactions, ok := data.([]accounts.Transaction)
	if !ok {
		fmt.Println("Failed to convert to []AccountDetails")
		return
	}
	//update excel sheet
	_, err = createExcelSheet(Transactions)
	if err != nil {
		fmt.Println(err, "error creating excel sheet")
		return
	}
	//update pdf file
	_, err = createPdf(Transactions)
	if err != nil {
		fmt.Println(err, "error creating pdf sheet")
		return
	}
	FullName, err := GetAdminName(r)
	if err != nil {
		log.Println(err)
	}
	pageData := AllTransactionsPageData{
		Transactions: Transactions,
		AdminName:    FullName,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/inflow.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/inflow.html", nil)
}

func (app *application) RenderOutflowPage(w http.ResponseWriter, r *http.Request) {
	//get all transactions
	accountNumber := os.Getenv("WITHDRAWAL_ACCOUNT_NUMBER")
	fmt.Println(accountNumber)
	data, err := accounts.ProcessAccount([]string{"", "acmt", "1012", accountNumber})

	if err != nil {
		fmt.Println(err)
	}
	Transactions, ok := data.([]accounts.Transaction)
	if !ok {
		fmt.Println("Failed to convert to []AccountDetails")
		return
	}

	FullName, err := GetAdminName(r)
	if err != nil {
		log.Println(err)
	}
	//update excel sheet
	_, err = createExcelSheet(Transactions)
	if err != nil {
		fmt.Println(err, "error creating excel sheet")
		return
	}
	//update pdf file
	_, err = createPdf(Transactions)
	if err != nil {
		fmt.Println(err, "error creating pdf sheet")
		return
	}

	pageData := AllTransactionsPageData{
		Transactions: Transactions,
		AdminName:    FullName,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/outflow.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/outflow.html", nil)
}

// RenderBatchTransactionPage renders a HTML page
func (app *application) RenderAllAccountsPage(w http.ResponseWriter, r *http.Request) {
	data, err := accounts.ProcessAccount([]string{"", "acmt", "1000"})
	if err != nil {
		fmt.Println(err)
		return
	}
	accountDetails, ok := data.([]accounts.AccountDetails)
	if !ok {
		fmt.Println("Failed to convert to []AccountDetails")
		return
	}
	FullName, err := GetAdminName(r)
	if err != nil {
		log.Println(err)
	}
	pageData := AllAccountPageData{
		Accounts:  accountDetails,
		AdminName: FullName,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/allAccounts.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/allAccounts.html", nil)
}

func GetAdminName(r *http.Request) (string, error) {
	session, err := store.Get(r, "JwtToken")
	if err != nil {
		return "", err
	}
	Fullname := session.Values["fullname"].(string)

	return Fullname, nil
}
func (app *application) RenderBeneficiariesPage(w http.ResponseWriter, r *http.Request) {

	FullName, err := GetAdminName(r)
	if err != nil {
		log.Println(err)
	}
	pageData := AllTransactionsPageData{
		AdminName: FullName,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/viewBeneficiaries.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/viewBeneficiaries.html", nil)
}
