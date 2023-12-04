package data

import (
	"github.com/ebitezion/backend-framework/internal/validator"
)

// user registration validation functions
func ValidateAccountBioData(v *validator.Validator, data *AccountBioData) {
	v.Check(data.FirstName != "", "firstName", "must be provided")
	v.Check(data.Surname != "", "surname", "must be provided")
	v.Check(data.HomeAddress != "", "homeAddress", "must be provided")
	v.Check(data.City != "", "city", "must be provided")
	v.Check(data.PhoneNumber != "", "phoneNumber", "must be provided")

	// Validate Identity sub-struct
	v.Check(data.Identity.BVN != "", "identity.bvn", "must be provided")
	v.Check(data.Identity.Passport != "", "identity.passport", "must be provided")
	v.Check(data.Identity.UtilityBill != "", "identity.utilityBill", "must be provided")
	v.Check(data.Identity.Country != "", "identity.country", "must be provided")
	v.Check(data.Picture != "", "picture", "must be provided")
}
func ValidateUserInformation(v *validator.Validator, data *User) {
	// General validation
	v.Check(data.AccountNumber != "", "accountNumber", "must be provided")

}

// ValidateAccountID validates a given AccountID struct and checks to make sure it passes all cases
func ValidateAccountID(v *validator.Validator, data *AccountID) {
	// General validation
	v.Check(data.AccountID != "", "accountID", "must be provided")

}

// ValidateAuthLoginData validates a given AuthLoginData struct and checks to make sure it passes all cases
func ValidateAuthLoginData(v *validator.Validator, data *AuthLoginData) {
	// General validation
	v.Check(data.Username != "", "username", "must be provided")
	v.Check(data.Password != "", "password", "must be provided")

}

// ValidateBeneficiaryData validates a given Beneficiary struct
func ValidateBeneficiaryData(v *validator.Validator, data *Beneficiary) {
	// General validation
	v.Check(data.FullName != "", "fullName", "must be provided")
	v.Check(data.BankName != "", "bankName", "must be provided")
	v.Check(data.BankAccountNumber != "", "bankAccountNumber", "must be provided")
	v.Check(data.BankRoutingNumber != "", "bankRoutingNumber", "must be provided")
	v.Check(data.SwiftCode != "", "swiftCode", "must be provided")
}

// ValidateAuthCreateData validates a given AuthCreateData struct
func ValidateAuthCreateData(v *validator.Validator, data *AuthCreateData) {
	// General validation
	v.Check(data.Username != "", "username", "must be provided")
	v.Check(data.AccountNumber != "", "accountNumber", "must be provided")
	v.Check(data.Password != "", "password", "must be provided")
}

// ValidateCreditInitiationData validates a given PaymentInitiationData struct
func ValidateCreditInitiationData(v *validator.Validator, data *PaymentInitiationData) {
	// General validation
	v.Check(data.SendersAccountNumber != "", "sendersAccountNumber", "must be provided")
	v.Check(data.ReceiversAccountNumber != "", "receiversAccountNumber", "must be provided")
	v.Check(data.Amount != "", "amount", "must be provided")
}

// ValidateDepositInitiationData validates a given DepositInitiationData struct
func ValidateDepositInitiationData(v *validator.Validator, data *DepositInitiationData) {
	// General validation
	v.Check(data.AccountNumber != "", "accountNumber", "must be provided")
	v.Check(data.Amount != "", "amount", "must be provided")
}

// ValidateUser validates a given User struct
func ValidateUser(v *validator.Validator, data *User) {
	// General validation
	v.Check(data.AccountNumber != "", "accountNumber", "must be provided")
}

// ValidateNewAccountRequestData validates a given NewAccountRequest struct
func ValidateNewAccountRequestData(v *validator.Validator, data *NewAccountRequest) {
	// General validation

	v.Check(data.AccountHolderGivenName != "", "accountHolderGivenName", "must be provided")
	v.Check(data.AccountHolderFamilyName != "", "accountHolderFamilyName", "must be provided")
	v.Check(data.AccountHolderDateOfBirth != "", "accountHolderDateOfBirth", "must be provided")
	v.Check(data.AccountHolderIdentificationNum != "", "accountHolderIdentificationNumber", "must be provided")
	v.Check(data.AccountHolderEmailAddress != "", "accountHolderEmailAddress", "must be provided")
	v.Check(data.AccountHolderContactNumber1 != "", "accountHolderContactNumber1", "must be provided")
	v.Check(data.AccountHolderContactNumber2 != "", "accountHolderContactNumber2", "must be provided")
	v.Check(data.AccountHolderAddressLine1 != "", "accountHolderAddressLine1", "must be provided")
	v.Check(data.AccountHolderAddressLine2 != "", "accountHolderAddressLine2", "must be provided")
	v.Check(data.AccountHolderAddressLine3 != "", "accountHolderAddressLine3", "must be provided")
	v.Check(data.AccountHolderPostalCode != "", "accountHolderPostalCode", "must be provided")
	v.Check(data.AccountHolderIdentificationType != "", "accountHolderIdentificationType", "must be provided")
	v.Check(data.Country != "", "country", "must be provided")

}
