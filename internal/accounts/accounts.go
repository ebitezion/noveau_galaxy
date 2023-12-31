package accounts

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ebitezion/backend-framework/internal/appauth"
	"github.com/ebitezion/backend-framework/internal/data"
	"github.com/ebitezion/backend-framework/internal/nuban"
	"github.com/ebitezion/backend-framework/internal/ukaccountgen"
	"github.com/shopspring/decimal"
)

/*
Accounts package to deal with all account related queries.

@TODO Implement the ISO20022 standard
http://www.iso20022.org/full_catalogue.page - acmt

@TODO Consider moving checkBalances, updateBalance to here

Accounts (acmt) transactions are as follows:
1  - AccountOpeningInstructionV05
2  - AccountDetailsConfirmationV05
3  - AccountModificationInstructionV05
5  - RequestForAccountManagementStatusReportV03
6  - AccountManagementStatusReportV04
7  - AccountOpeningRequestV02
8  - AccountOpeningAmendmentRequestV02
9  - AccountOpeningAdditionalInformationRequestV02
10 - AccountRequestAcknowledgementV02
11 - AccountRequestRejectionV02
12 - AccountAdditionalInformationRequestV02
13 - AccountReportRequestV02
14 - AccountReportV02
15 - AccountExcludedMandateMaintenanceRequestV02
16 - AccountExcludedMandateMaintenanceAmendmentRequestV02
17 - AccountMandateMaintenanceRequestV02
18 - AccountMandateMaintenanceAmendmentRequestV02
19 - AccountClosingRequestV02
20 - AccountClosingAmendmentRequestV02
21 - AccountClosingAdditionalInformationRequestV02
22 - IdentificationModificationAdviceV02
23 - IdentificationVerificationRequestV02
24 - IdentificationVerificationReportV02

### Custom functionality
1000 - ListAllAccounts (@FIXME Used for now by anyone, close down later)
1001 - ListSingleAccount
1002 - CheckAccountByID
1003 - BalanceEnquiry
1004 - AccountHistory
1005 - AccountMetaData
1006 - AllAccounts
1007 - UpdateAccountInformation
1008 - AllTransactions
1009 - CreateSpecialAccount
1010 - BlockAccount
1011 - UnblockAccount
1012 - OutflowHistory

*/

/*
acmt~1~

	AccountHolderGivenName~
	AccountHolderFamilyName~
	AccountHolderDateOfBirth~
	AccountHolderIdentificationNumber~
	AccountHolderContactNumber1~
	AccountHolderContactNumber2~
	AccountHolderEmailAddress~
	AccountHolderAddressLine1~
	AccountHolderAddressLine2~
	AccountHolderAddressLine3~
	AccountHolderPostalCode
*/
type AccountHolder struct {
	AccountNumber string
	BankNumber    string
}

type AccountHolderDetails struct {
	AccountNumber        string
	BankNumber           string
	GivenName            string
	FamilyName           string
	DateOfBirth          string
	IdentificationNumber string
	IdentificationType   string
	ContactNumber1       string
	ContactNumber2       string
	EmailAddress         string
	AddressLine1         string
	AddressLine2         string
	AddressLine3         string
	PostalCode           string
	Image                string
	Country              string
}
type SpecialAccountDetails struct {
	AccountNumber     string
	BankNumber        string
	AccountHolderName string
	Status            string
	AccountBalance    decimal.Decimal
	Overdraft         decimal.Decimal
	AvailableBalance  decimal.Decimal
	AccountProfile    AccountProfile
}
type AccountProfile struct {
	Creator       string
	AccountName   string
	Purpose       string
	AccountNumber string
	BankNumber    string
}

type AccountDetails struct {
	AccountNumber        string
	BankNumber           string
	AccountHolderName    string
	Status               string
	AccountBalance       decimal.Decimal
	Overdraft            decimal.Decimal
	AvailableBalance     decimal.Decimal
	AccountHolderDetails AccountHolderDetails
}

type BalanceEnquiry struct {
	AccountHolderName string `json:"accountHolderName"`
	AccountNumber     string `json:"accountNumber"`
	LedgerBalance     string `json:"ledgerBalance"`
	Status            string `json:"status"`
}

type Transaction struct {
	ID                    int     `json:"id"`
	Transaction           string  `json:"transaction"`
	Type                  int     `json:"type"`
	SenderAccountNumber   string  `json:"senderAccountNumber"`
	SenderBankNumber      string  `json:"senderBankNumber"`
	ReceiverAccountNumber string  `json:"receiverAccountNumber"`
	ReceiverBankNumber    string  `json:"receiverBankNumber"`
	TransactionAmount     float64 `json:"transactionAmount"`
	Narration             string  `json:"narration"`
	FeeAmount             float64 `json:"feeAmount"`
	Timestamp             string  `json:"timestamp"`
	Initiator             string  `json:"initiator"`
}

// Set up some defaults
const (
	BANK_NUMBER       = "a0299975-b8e2-4358-8f1a-911ee12dbaac"
	OPENING_BALANCE   = 100.
	OPENING_OVERDRAFT = 0.
)

var accountMutex sync.Mutex

func ProcessAccount(data []string) (result interface{}, err error) {

	if len(data) < 3 {
		return "", errors.New("accounts.ProcessAccount: Not enough fields, minimum 3")
	}

	acmtType, err := strconv.ParseInt(data[2], 10, 64)
	if err != nil {
		return "", errors.New("accounts.ProcessAccount: Could not get ACMT type")
	}

	// Switch on the acmt type
	switch acmtType {
	case 1, 7:
		/*
		   @TODO
		   The differences between AccountOpeningInstructionV05 and AccountOpeningRequestV02 will be explored in detail, for now we treat the same - open an account
		*/

		result, err = openAccount(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
		break
	case 1000:
		//TODO: check permissions for this
		result, err = fetchAccounts(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
		break
	case 1001:
		result, err = fetchSingleAccount(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
		break
	case 1002:
		if len(data) < 4 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		result, err = fetchSingleAccountByID(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
		break
	case 1003:
		if len(data) < 3 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		result, err = fetchAccountDetails(data)

		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
		break
	case 1004:
		if len(data) < 3 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		result, err = fetchAccountHistory(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
		break
	case 1007:
		if len(data) < 3 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		result, err = updateAccountInformation(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
		break
	case 1008:
		if len(data) < 3 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		result, err = allTransactions(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
		break
	case 1009:
		if len(data) < 2 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		result, err = openSpecialAccount(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
		break
	case 1010:
		if len(data) < 2 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		result, err = blockAccount(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
		break
	case 1011:
		if len(data) < 3 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		result, err = unblockAccount(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
		break
	case 1012:
		if len(data) < 3 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		result, err = fetchOutflowHistory(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
		break

	default:
		err = errors.New("accounts.ProcessAccount: ACMT transaction code invalid")
		break

	}

	return
}
func unblockAccount(data []string) (result string, err error) {
	// Validate string against required info/length
	if len(data) < 2 {
		err := errors.New("accounts.closeAccount: Not all fields present")
		return "", err
	}

	// Check if account already exists, check on ID number
	accountHolder, _ := getAccountMeta(data[3])
	if accountHolder.AccountNumber == "" {
		return "", errors.New("accounts.closeAccount: Account does not exist. " + accountHolder.AccountNumber)
	}

	err = doUnblockAccount(accountHolder.AccountNumber)
	if err != nil {
		return "", errors.New("accounts.closeAccount: " + err.Error())
	}

	return "Account Unblocked Successfully", nil
}
func blockAccount(data []string) (result string, err error) {
	// Validate string against required info/length
	if len(data) < 2 {
		err := errors.New("accounts.closeAccount: Not all fields present")
		return "", err
	}
	fmt.Println(data[3])
	// Check if account already exists, check on ID number
	accountHolder, _ := getAccountMeta(data[3])
	if accountHolder.AccountNumber == "" {
		return "", errors.New("accounts.closeAccount: Account does not exist. " + accountHolder.AccountNumber)
	}

	err = doBlockAccount(accountHolder.AccountNumber)
	if err != nil {
		return "", errors.New("accounts.closeAccount: " + err.Error())
	}

	return "Account Blocked Successfully", nil
}

// GetBenefciaries gets all bebeficiaries related to a particular user
func GetBenefciaries(accountNumber string) (beneficiaries []data.Beneficiary, err error) {

	//get user_id
	userId, err := GetUserId(accountNumber)
	if err != nil {
		if err == data.ErrRecordNotFound {
			return nil, errors.New("accounts.CreateBeneficiary: No records found")

		} else {
			return nil, err
		}
	}

	//get beneficiaries
	beneficiaries, err = FetchBenefciaries(userId)
	if err != nil {
		return nil, err

	}

	return beneficiaries, nil
}

// CreateBeneficiary creates a new beneficiary tied to an account
func CreateBeneficiary(beneficiary *data.Beneficiary) error {

	//get user_id
	userId, err := GetUserId(beneficiary.UserAccountNumber)
	if err != nil {
		if err == data.ErrRecordNotFound {
			return errors.New("accounts.CreateBeneficiary: No records found")

		} else {
			return err
		}
	}

	//create new beneficiary
	err = StoreBeneficiary(beneficiary, userId)
	if err != nil {
		return err

	}
	return nil

}
func FetchAccountNumber(username string) (AccountNumber string, Fullname string, err error) {

	if username == "" {
		return "", "", errors.New("accounts.fetchAccountMeta: Username not present")
	}

	accountNumber, fullname, err := getSingleAccountNumberByUsername(username)
	if err != nil {
		return "", "", errors.New("accounts.fetchAccountMeta: " + err.Error())
	}

	return accountNumber, fullname, nil
}
func FetchAccountMeta(accountNumber string) (AccountHolderDetails *AccountHolderDetails, err error) {

	if accountNumber == "" {
		return nil, errors.New("accounts.fetchAccountMeta: Account number not present")
	}

	accountMeta, err := getAccountMeta(accountNumber)
	if err != nil {
		return nil, errors.New("accounts.fetchAccountMeta: " + err.Error())
	}

	return &accountMeta, nil
}
func allTransactions(data []string) (result []Transaction, err error) {
	// Fetch all accounts. This fetches non-sensitive information (no balances)
	transactions, err := getAllTransactions()

	if err != nil {
		return nil, errors.New("accounts.allTransactions: " + err.Error())
	}

	return transactions, nil
}

func fetchAccounts(data []string) (result []AccountDetails, err error) {
	// Fetch all accounts. This fetches non-sensitive information (no balances)
	accounts, err := getAllAccountDetails()
	if err != nil {
		return nil, errors.New("accounts.fetchAccounts: " + err.Error())
	}

	return accounts, nil
}
func openSpecialAccount(data []string) (result string, err error) {
	// Validate string against required info/length
	if len(data) < 3 {
		err = errors.New("accounts.openAccount: Not all fields present")
		//@TODO Add to documentation rather than returning here
		//result = "ERROR: acmt transactions must be as follows:acmt~AcmtType~AccountHolderGivenName~AccountHolderFamilyName~AccountHolderDateOfBirth~AccountHolderIdentificationNumber~AccountHolderContactNumber1~AccountHolderContactNumber2~AccountHolderEmailAddress~AccountHolderAddressLine1~AccountHolderAddressLine2~AccountHolderAddressLine3~AccountHolderPostalCode"
		return
	}

	//@TODO add check to know if an account already exist

	// @FIXME: Remove new line from data
	data[len(data)-1] = strings.Replace(data[len(data)-1], "\n", "", -1)

	// Create account
	accountDetails, err := setSpecialAccountDetails(data)
	if err != nil {
		return "", errors.New("accounts.openAccount: " + err.Error())
	}

	err = createSpecialAccount(&accountDetails)
	if err != nil {
		return "", errors.New("accounts.openAccount: " + err.Error())
	}

	result = accountDetails.AccountNumber
	return
}
func openAccount(data []string) (result string, err error) {
	// Validate string against required info/length
	if len(data) < 14 {
		err = errors.New("accounts.openAccount: Not all fields present")
		//@TODO Add to documentation rather than returning here
		//result = "ERROR: acmt transactions must be as follows:acmt~AcmtType~AccountHolderGivenName~AccountHolderFamilyName~AccountHolderDateOfBirth~AccountHolderIdentificationNumber~AccountHolderContactNumber1~AccountHolderContactNumber2~AccountHolderEmailAddress~AccountHolderAddressLine1~AccountHolderAddressLine2~AccountHolderAddressLine3~AccountHolderPostalCode"
		return
	}

	// Test: acmt~1~Kyle~Redelinghuys~19000101~190001011234098~1112223456~~email@domain.com~Physical Address 1~~~1000
	// Check if account already exists, check on ID number
	accountHolder, _ := getAccountMeta(data[6])

	if accountHolder.AccountNumber != "" {
		return "", errors.New("accounts.openAccount: Account already open. " + accountHolder.AccountNumber)
	}

	// @FIXME: Remove new line from data
	data[len(data)-1] = strings.Replace(data[len(data)-1], "\n", "", -1)

	// Create account
	accountHolderObject, err := setAccountDetails(data)
	if err != nil {
		return "", errors.New("accounts.openAccount: " + err.Error())
	}
	accountHolderDetailsObject, err := setAccountHolderDetails(data)
	if err != nil {
		return "", errors.New("accounts.openAccount: " + err.Error())
	}
	err = createAccount(&accountHolderObject, &accountHolderDetailsObject)
	if err != nil {
		return "", errors.New("accounts.openAccount: " + err.Error())
	}

	result = accountHolderObject.AccountNumber
	return
}
func updateAccountInformation(data []string) (result string, err error) {
	// Validate string against required info/length
	if len(data) < 14 {
		err = errors.New("accounts.openAccount: Not all fields present")
		//@TODO Add to documentation rather than returning here
		//result = "ERROR: acmt transactions must be as follows:acmt~AcmtType~AccountHolderGivenName~AccountHolderFamilyName~AccountHolderDateOfBirth~AccountHolderIdentificationNumber~AccountHolderContactNumber1~AccountHolderContactNumber2~AccountHolderEmailAddress~AccountHolderAddressLine1~AccountHolderAddressLine2~AccountHolderAddressLine3~AccountHolderPostalCode"
		return
	}

	// @FIXME: Remove new line from data
	data[len(data)-1] = strings.Replace(data[len(data)-1], "\n", "", -1)

	accountHolderDetailsObject, err := setAccountDetailsForUpdate(data)
	if err != nil {
		return "", errors.New("accounts.openAccount: " + err.Error())
	}

	err = updateAccount(&accountHolderDetailsObject)
	if err != nil {
		return "", errors.New("accounts.openAccount: " + err.Error())
	}

	result = "Update Successful"
	return

}
func closeAccount(data []string) (result string, err error) {
	// Validate string against required info/length
	if len(data) < 14 {
		err = errors.New("accounts.closeAccount: Not all fields present")
		return
	}

	// Check if account already exists, check on ID number
	accountHolder, _ := getAccountMeta(data[6])
	if accountHolder.AccountNumber == "" {
		return "", errors.New("accounts.closeAccount: Account does not exist. " + accountHolder.AccountNumber)
	}

	// @FIXME: Remove new line from data
	data[len(data)-1] = strings.Replace(data[len(data)-1], "\n", "", -1)

	// Delete account
	accountHolderObject, err := setAccountDetails(data)
	if err != nil {
		return "", errors.New("accounts.closeAccount: " + err.Error())
	}
	accountHolderDetailsObject, err := setAccountHolderDetails(data)
	if err != nil {
		return "", errors.New("accounts.closeAccount: " + err.Error())
	}
	err = deleteAccount(&accountHolderObject, &accountHolderDetailsObject)
	if err != nil {
		return "", errors.New("accounts.closeAccount: " + err.Error())
	}

	return
}

func setAccountDetails(data []string) (accountDetails AccountDetails, err error) {
	fmt.Println(data)
	if data[4] == "" {
		return AccountDetails{}, errors.New("accounts.setAccountDetails: Family name cannot be empty")
	}
	if data[3] == "" {
		return AccountDetails{}, errors.New("accounts.setAccountDetails: Given name cannot be empty")
	}
	nubanGenerator := nuban.NewNUBANGenerator()
	nuban := nubanGenerator.GenerateNUBAN()

	ukaccountgenerator := ukaccountgen.New().GenerateUKAccountNumber()

	accountDetails.AccountNumber = ukaccountgenerator
	accountDetails.BankNumber = nuban
	accountDetails.AccountHolderName = data[4] + "," + data[3] // Family Name, Given Name
	accountDetails.AccountBalance = decimal.NewFromFloat(OPENING_BALANCE)
	accountDetails.Overdraft = decimal.NewFromFloat(OPENING_OVERDRAFT)
	accountDetails.AvailableBalance = decimal.NewFromFloat(OPENING_BALANCE + OPENING_OVERDRAFT)

	return
}

func setSpecialAccountDetails(data []string) (specialAccountDetails SpecialAccountDetails, err error) {
	if len(data) < 3 {
		return SpecialAccountDetails{}, errors.New("accounts.setAccountHolderDetails: Not all field values present")
	}
	//@TODO: Test date parsing in format ddmmyyyy
	if data[4] == "" {
		return SpecialAccountDetails{}, errors.New("accounts.setAccountHolderDetails:  cannot be empty")
	}
	if data[3] == "" {
		return SpecialAccountDetails{}, errors.New("accounts.setAccountHolderDetails: accountName name cannot be empty")
	}

	// @TODO Integrity checks
	nubanGenerator := nuban.NewNUBANGenerator()
	nuban := nubanGenerator.GenerateNUBAN()

	ukaccountgenerator := ukaccountgen.New().GenerateUKAccountNumber()

	specialAccountDetails.AccountNumber = ukaccountgenerator
	specialAccountDetails.BankNumber = nuban
	specialAccountDetails.AccountHolderName = data[3] //AcccountName
	specialAccountDetails.AccountBalance = decimal.NewFromFloat(OPENING_BALANCE)
	specialAccountDetails.Overdraft = decimal.NewFromFloat(OPENING_OVERDRAFT)
	specialAccountDetails.AvailableBalance = decimal.NewFromFloat(OPENING_BALANCE + OPENING_OVERDRAFT)
	specialAccountDetails.AccountProfile.AccountName = data[3] //AcccountName
	specialAccountDetails.AccountProfile.Creator = data[4]
	specialAccountDetails.AccountProfile.Purpose = data[5]

	return
}
func setAccountHolderDetails(data []string) (accountHolderDetails AccountHolderDetails, err error) {
	if len(data) < 14 {
		return AccountHolderDetails{}, errors.New("accounts.setAccountHolderDetails: Not all field values present")
	}
	//@TODO: Test date parsing in format ddmmyyyy
	if data[4] == "" {
		return AccountHolderDetails{}, errors.New("accounts.setAccountHolderDetails: Family name cannot be empty")
	}
	if data[3] == "" {
		return AccountHolderDetails{}, errors.New("accounts.setAccountHolderDetails: Given name cannot be empty")
	}

	// @TODO Integrity checks
	nubanGenerator := nuban.NewNUBANGenerator()
	nuban := nubanGenerator.GenerateNUBAN()

	ukaccountgenerator := ukaccountgen.New().GenerateUKAccountNumber()

	accountHolderDetails.AccountNumber = ukaccountgenerator
	accountHolderDetails.BankNumber = nuban
	accountHolderDetails.GivenName = data[3]
	accountHolderDetails.FamilyName = data[4]
	accountHolderDetails.DateOfBirth = data[5]
	accountHolderDetails.IdentificationNumber = data[6]
	accountHolderDetails.ContactNumber1 = data[7]
	accountHolderDetails.ContactNumber2 = data[8]
	accountHolderDetails.EmailAddress = data[9]
	accountHolderDetails.AddressLine1 = data[10]
	accountHolderDetails.AddressLine2 = data[11]
	accountHolderDetails.AddressLine3 = data[12]
	accountHolderDetails.PostalCode = data[13]
	accountHolderDetails.Image = data[14]
	accountHolderDetails.IdentificationType = data[15]
	accountHolderDetails.Country = data[16]

	return
}

func setAccountDetailsForUpdate(data []string) (accountHolderDetails AccountHolderDetails, err error) {
	if len(data) < 14 {
		return AccountHolderDetails{}, errors.New("accounts.setAccountHolderDetails: Not all field values present")
	}
	//@TODO: Test date parsing in format ddmmyyyy
	if data[4] == "" {
		return AccountHolderDetails{}, errors.New("accounts.setAccountHolderDetails: Family name cannot be empty")
	}
	if data[3] == "" {
		return AccountHolderDetails{}, errors.New("accounts.setAccountHolderDetails: Given name cannot be empty")
	}

	accountHolderDetails.GivenName = data[3]
	accountHolderDetails.FamilyName = data[4]
	accountHolderDetails.DateOfBirth = data[5]
	accountHolderDetails.IdentificationNumber = data[6]
	accountHolderDetails.ContactNumber1 = data[7]
	accountHolderDetails.ContactNumber2 = data[8]
	accountHolderDetails.EmailAddress = data[9]
	accountHolderDetails.AddressLine1 = data[10]
	accountHolderDetails.AddressLine2 = data[11]
	accountHolderDetails.AddressLine3 = data[12]
	accountHolderDetails.PostalCode = data[13]
	accountHolderDetails.Image = data[14]
	accountHolderDetails.IdentificationType = data[15]
	accountHolderDetails.Country = data[16]
	accountHolderDetails.AccountNumber = data[17]

	return
}

// @TODO Remove this after testing, security risk

func fetchSingleAccount(data []string) (result string, err error) {
	// Fetch user account. Must be user logged in
	tokenUser, err := appauth.GetUserFromToken(data[0])
	if err != nil {
		return "", errors.New("accounts.fetchSingleAccount: " + err.Error())
	}
	account, err := getSingleAccountDetail(tokenUser)
	if err != nil {
		return "", errors.New("accounts.fetchSingleAccount: " + err.Error())
	}

	// Parse into nice result string
	jsonAccount, err := json.Marshal(account)
	if err != nil {
		return "", errors.New("accounts.fetchSingleAccount: " + err.Error())
	}

	result = string(jsonAccount)
	return
}

func fetchSingleAccountByID(data []string) (result string, err error) {
	// Format: token~acmt~1002~USERID
	userID := data[3]
	if userID == "" {
		return "", errors.New("accounts.fetchSingleAccountByID: User ID not present")
	}

	userAccountNumber, err := getSingleAccountNumberByID(userID)
	if err != nil {
		return "", errors.New("accounts.fetchSingleAccountByID: " + err.Error())
	}

	result = userAccountNumber
	return
}
func fetchAccountDetails(data []string) (result *BalanceEnquiry, err error) {
	accountNumber := data[3]
	if accountNumber == "" {
		return nil, errors.New("accounts.fetchSingleAccountByID: Account number not present")
	}

	balanceEnquiry, err := GetBalanceDetails(accountNumber)
	if err != nil {
		return nil, errors.New("accounts.fetchSingleAccountByID: " + err.Error())
	}

	return &balanceEnquiry, nil
}

func fetchAccountHistory(data []string) (result []Transaction, err error) {

	// Format: token~acmt~1002~USERID
	accountNumber := data[3]
	if accountNumber == "" {
		return nil, errors.New("accounts.fetchSingleAccountByID: Account number not present")
	}
	//check if  accounts is valid

	exists, err := CheckIfAccountNumberExists(accountNumber)
	if err != nil {
		return nil, errors.New("payments.fetchSingleAccountByID: " + err.Error())
	}
	// Check the result.
	if !exists {
		return nil, errors.New("payments.fetchSingleAccountByID: " + " Account Not valid")
	}
	history, err := GetAccountHistory(accountNumber)
	if err != nil {
		return nil, errors.New("accounts.fetchSingleAccountByID: " + err.Error())
	}

	return history, nil
}
func fetchOutflowHistory(data []string) (result []Transaction, err error) {

	// Format: token~acmt~1002~USERID
	accountNumber := data[3]
	if accountNumber == "" {
		return nil, errors.New("accounts.fetchSingleAccountByID: Account number not present")
	}
	//check if  accounts is valid

	exists, err := CheckIfAccountNumberExists(accountNumber)
	if err != nil {
		return nil, errors.New("payments.fetchSingleAccountByID: " + err.Error())
	}
	// Check the result.
	if !exists {
		return nil, errors.New("payments.fetchSingleAccountByID: " + " Account Not valid")
	}
	history, err := GetOutflowHistory(accountNumber)
	if err != nil {
		return nil, errors.New("accounts.fetchSingleAccountByID: " + err.Error())
	}

	return history, nil
}

// CheckIfValueExists checks if a given value is in the specified table and returns a boolean
func CheckIfAccountNumberExists(accountNumber string) (bool, error) {
	query := "SELECT COUNT(*) FROM accounts WHERE accountNumber = ?"
	// Declare a variable to store the count.
	var count int

	// Use the context.WithTimeout() function to create a context.Context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryRowContext() to execute the query and get the count.
	err := Config.Db.QueryRowContext(ctx, query, accountNumber).Scan(&count)
	if err != nil {
		// Print the error message for debugging purposes.
		fmt.Println("Error executing query:", err)
		return false, err
	}

	// If the count is greater than 0, the value exists in the database.
	return count > 0, nil
}
func Base64ToImage(encodedString string, outputPath string) error {
	// Decode Base64 string to binary data
	decoded, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return err
	}

	// Create an image from the binary data
	img, _, err := image.Decode(strings.NewReader(string(decoded)))
	if err != nil {
		return err
	}

	// Create a new file to save the image
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Save the image as JPEG
	err = jpeg.Encode(outputFile, img, nil)
	if err != nil {
		return err
	}

	return nil
}

func ImageToBase64(imagePath string) (string, error) {
	// Open the image file
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer imageFile.Close()

	// Decode the image file
	img, _, err := image.Decode(imageFile)
	if err != nil {
		return "", err
	}

	// Encode the image as Base64
	var encodedString string
	buffer := new(strings.Builder)
	err = jpeg.Encode(buffer, img, nil)
	if err != nil {
		return "", err
	}
	encodedString = base64.StdEncoding.EncodeToString([]byte(buffer.String()))

	return encodedString, nil
}
func ImageToBase64FromRequest(r *http.Request) (string, error) {
	// Parse the form data to get the image file
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		return "", fmt.Errorf("error parsing form: %v", err)
	}

	file, _, err := r.FormFile("profilePicture")
	if err != nil {
		return "", fmt.Errorf("error retrieving file from form: %v", err)
	}
	defer file.Close()

	// Decode the image file
	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("error decoding image: %v", err)
	}

	// Encode the image as Base64
	var encodedString string
	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, img, nil)
	if err != nil {
		return "", fmt.Errorf("error encoding image to Base64: %v", err)
	}
	encodedString = base64.StdEncoding.EncodeToString(buffer.Bytes())

	return encodedString, nil
}
