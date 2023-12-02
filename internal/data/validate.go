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

}
func ValidateBeneficiaryData(v *validator.Validator, data *Beneficiary) {

}
func ValidateAccountID(v *validator.Validator, data *AccountID) {

}
func ValidateAuthLoginData(v *validator.Validator, data *AuthLoginData) {

}
func ValidateAuthCreateData(v *validator.Validator, data *AuthCreateData) {

}
func ValidateCreditInitiationData(v *validator.Validator, data *PaymentInitiationData) {

}
func ValidateDepositInitiationData(v *validator.Validator, data *DepositInitiationData) {

}
func ValidateUser(v *validator.Validator, data *User) {

}
