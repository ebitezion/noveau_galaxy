package data

import "database/sql"

type AccountBioData struct {
	Surname     string `json:"surname"`
	FirstName   string `json:"firstName"`
	HomeAddress string `json:"homeAddress"`
	City        string `json:"city"`
	PhoneNumber string `json:"phoneNumber"`
	Identity    struct {
		BVN         string `json:"bvn"`
		Passport    string `json:"passport"`
		UtilityBill string `json:"utilityBill"`
		Country     string `json:"country"`
	} `json:"identity"`
	Picture string `json:"picture"`
}

// cconnection to DB resources
type AccountModel struct {
	DB *sql.DB
}
