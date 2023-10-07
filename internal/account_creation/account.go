package accountcreation

import (
	"database/sql"
	"fmt"

	"github.com/ebitezion/backend-framework/internal/data"
)

// cconnection to DB resources
type AccountModel struct {
	DB *sql.DB
}

// Register simply saves the biodata info of the user and return an account no or an error
func Register(biodata data.AccountBioData) (accountNo string, err error) {
	fmt.Println("sal data: ", biodata)
	query := ` INSERT INTO biodata
    (surname,firstName)
    VALUES ( ?, ?)`

	//Insert to DB
	err = data.Models{}.Insert(query, []interface{}{biodata.Surname, biodata.FirstName})

	if err != nil {
		return "", err
	}

	//Handle errora

	//on successful insertion, generate account no,

	//save on account_table

	//Handle error
	accountNo = "12345"
	// return the generated account
	return accountNo, nil
}
