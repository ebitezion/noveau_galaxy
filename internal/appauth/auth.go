package appauth

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"time"

	"gopkg.in/redis.v3"

	"github.com/ebitezion/backend-framework/internal/configuration"
	"github.com/twinj/uuid"
)

const (
	TOKEN_TTL = time.Hour // One hour
)

var Config configuration.Configuration

func SetConfig(config *configuration.Configuration) {
	Config = *config
}

func ProcessAppAuth(data []string) (result string, err error) {
	//@TODO: Change from []string to something more solid, struct/interface/key-pair
	if len(data) < 3 {
		return "", errors.New("appauth.ProcessAppAuth: Not all required fields present")
	}
	switch data[2] {
	// Auth an existing account
	case "1":
		// TOKEN~appauth~1
		if len(data) < 3 {
			return "", errors.New("appauth.ProcessAppAuth: Not all required fields present")
		}
		err := CheckToken(data[0])
		if err != nil {
			return "", err
		}
		return result, nil
	// Log in
	case "2":
		if len(data) < 5 {
			return "", errors.New("appauth.ProcessAppAuth: Not all required fields present")
		}
		result, err = CreateToken(data[3], data[4])
		if err != nil {
			return "", err
		}
		return result, nil
	// Create an account
	case "3":
		if len(data) < 5 {
			return "", errors.New("appauth.ProcessAppAuth: Not all required fields present")
		}
		result, err = CreateUserPassword(data[3], data[4], data[5])
		if err != nil {
			return "", err
		}
		return result, nil
	// Remove an account
	case "4":
		if len(data) < 5 {
			return "", errors.New("appauth.ProcessAppAuth: Not all required fields present")
		}
		result, err = RemoveUserPassword(data[3], data[4])
		if err != nil {
			return "", err
		}
		return result, nil
	}
	return "", errors.New("appauth.ProcessAppAuth: No valid option chosen")
}

func CreateUserPassword(accountNumber string, password string, userame string) (result string, err error) {
	//check if account is valid

	rows, err := Config.Db.Query("SELECT `accountNumber` FROM `accounts` WHERE `accountNumber` = ?", accountNumber)
	if err != nil {
		return "", errors.New("appauth.CreateUserPassword: Error with select query. " + err.Error())
	}
	defer rows.Close()

	// @TODO Must be easy way to get row count returned
	count := 0
	for rows.Next() {
		count++
	}

	if count == 0 {
		return "", errors.New("appauth.CreateUserPassword: Invalid AccountNumber")
	}

	//TEST 0~appauth~3~181ac0ae-45cb-461d-b740-15ce33e4612f~testPassword
	// Generate hash
	hasher := sha512.New()
	hasher.Write([]byte(password))
	hash := hex.EncodeToString(hasher.Sum(nil))

	// Check for existing account
	rows, err = Config.Db.Query("SELECT `accountNumber` FROM `accounts_auth` WHERE `accountNumber` = ?", accountNumber)
	if err != nil {
		return "", errors.New("appauth.CreateUserPassword: Error with select query. " + err.Error())
	}
	defer rows.Close()

	// @TODO Must be easy way to get row count returned
	count = 0
	for rows.Next() {
		count++
	}

	if count > 0 {
		return "", errors.New("appauth.CreateUserPassword: Account already exists")
	}

	// Prepare statement for inserting data
	insertStatement := "INSERT INTO accounts_auth (`accountNumber`, `password`,`username`) "
	insertStatement += "VALUES(?, ?,?)"
	stmtIns, err := Config.Db.Prepare(insertStatement)
	if err != nil {
		return "", errors.New("appauth.CreateUserPassword: Error with insert. " + err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtIns.Exec(accountNumber, hash, userame)

	if err != nil {
		return "", errors.New("appauth.CreateUserPassword: Could not save account. " + err.Error())
	}

	result = "Successfully created account"
	return
}

func RemoveUserPassword(user string, hashedPassword string) (result string, err error) {
	// Check for existing account
	rows, err := Config.Db.Query("SELECT `accountNumber` FROM `accounts_auth` WHERE `accountNumber` = ?", user)
	if err != nil {
		return "", errors.New("appauth.RemoveUserPassword: Error with select query. " + err.Error())
	}
	defer rows.Close()

	// @TODO Must be easy way to get row count returned
	count := 0
	for rows.Next() {
		count++
	}

	if count == 0 {
		return "", errors.New("appauth.RemoveUserPassword: Account auth does not exists")
	}

	// Prepare statement for inserting data
	delStatement := "DELETE FROM accounts_auth WHERE `accountNumber` = ? AND `password` = ? "
	stmtDel, err := Config.Db.Prepare(delStatement)
	if err != nil {
		return "", errors.New("appauth.RemoveUserPassword: Error with delete. " + err.Error())
	}
	defer stmtDel.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtDel.Exec(user, hashedPassword)

	if err != nil {
		return "", errors.New("appauth.RemoveUserPassword: Could not delete account. " + err.Error())
	}

	result = "Successfully deleted account"
	return
}

func CreateToken(user string, password string) (token string, err error) {
	//check if account is valid

	//check if password is correct
	rows, err := Config.Db.Query("SELECT `password` FROM `accounts_auth` WHERE `accountNumber` = ?", user)
	if err != nil {
		return "", errors.New("appauth.CreateToken: Error with select query. " + err.Error())
	}
	defer rows.Close()

	count := 0
	hashedPassword := ""
	for rows.Next() {
		if err := rows.Scan(&hashedPassword); err != nil {
			return "", errors.New("appauth.CreateToken: Could not retreive account details")
		}
		count++
	}

	// Generate hash
	hasher := sha512.New()
	hasher.Write([]byte(password))
	hash := hex.EncodeToString(hasher.Sum(nil))

	if hash != hashedPassword {
		return "", errors.New("appauth.CreateToken: Authentication credentials invalid")
	}

	newUuid := uuid.NewV4()
	token = newUuid.String()

	// @TODO Remove all tokens for this user
	err = Config.Redis.Set(token, user, TOKEN_TTL).Err()
	if err != nil {
		return "", errors.New("appauth.CreateToken: Could not set token. " + err.Error())
	}

	return
}

// CheckIfValueExists checks if a given value is in the specified table and returns a boolean
func CheckIfValueExists(Query string, Args []interface{}) (bool, error) {

	// Declare a variable to store the count.
	var count int

	// Use the context.WithTimeout() function to create a context.Context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryRowContext() to execute the query and get the count.
	err := Config.Db.QueryRowContext(ctx, Query, Args...).Scan(&count)
	if err != nil {
		return false, err
	}

	// If the count is greater than 0, the invitee is already in the database.
	return count > 0, nil
}
func RemoveToken(token string) (result string, err error) {
	//TEST 0~appauth~480e67e3-e2c9-48ee-966c-8d251474b669
	_, err = Config.Redis.Del(token).Result()

	if err == redis.Nil {
		return "", errors.New("appauth.RemoveToken: Token not found. " + err.Error())
	} else if err != nil {
		return "", errors.New("appauth.RemoveToken: Could not remove token. " + err.Error())
	} else {
		result = "Token removed"
	}

	return
}

func CheckToken(token string) (err error) {
	//TEST 0~appauth~480e67e3-e2c9-48ee-966c-8d251474b669
	user, err := Config.Redis.Get(token).Result()

	if err == redis.Nil {
		return errors.New("appauth.CheckToken: Token not found. " + err.Error())
	} else if err != nil {
		return errors.New("appauth.CheckToken: Could not get token. " + err.Error())
	} else {
		// Extend token
		err := Config.Redis.Set(user, token, TOKEN_TTL).Err()
		if err != nil {
			return errors.New("appauth.CheckToken: Could not extend token. " + err.Error())
		}
	}

	return nil
}

func GetUserFromToken(token string) (user string, err error) {
	//TEST 0~appauth~~181ac0ae-45cb-461d-b740-15ce33e4612f~testPassword
	user, err = Config.Redis.Get(token).Result()
	if err != nil {
		return "", errors.New("appauth.GetUserFromToken: Could not get token. " + err.Error())
	}

	// If valid then extend
	if user != "" {
		err := Config.Redis.Set(user, token, TOKEN_TTL).Err()
		if err != nil {
			return "", errors.New("appauth.GetUserFromToken: Could not extend token. " + err.Error())
		}
	}

	return
}

func AssignRolePrivilege(adminUsername string, targetUsername string, targetRole string, targetPrivilege string) error {
	// Open a connection to your MySQL database
	// db, err := sql.Open("mysql", "username:password@tcp(your-mysql-host:port)/your-database")
	// if err != nil {
	// 	return err
	// }
	// defer db.Close()

	// Check if the admin exists and has sufficient privileges
	var adminRole string
	err := Config.Db.QueryRow("SELECT role FROM users WHERE username = ?", adminUsername).Scan(&adminRole)
	if err != nil {
		return err
	}

	if adminRole != "admin" {
		return errors.New("Admin does not have sufficient privileges")
	}

	// Get the target user's ID
	var targetID int
	err = Config.Db.QueryRow("SELECT user_id FROM users WHERE username = ?", targetUsername).Scan(&targetID)
	if err != nil {
		return err
	}

	if targetID == 0 {
		return errors.New("Target user not found")
	}

	// Assign the role to the target user
	_, err = Config.Db.Exec("UPDATE users SET role = ? WHERE user_id = ?", targetRole, targetID)
	if err != nil {
		return err
	}

	// Get the privilege ID for the specified privilege and role
	var privilegeID int
	err = Config.Db.QueryRow("SELECT privilege_id FROM privileges WHERE role = ? AND privilege_name = ?", targetRole, targetPrivilege).Scan(&privilegeID)
	if err != nil {
		return err
	}

	if privilegeID == 0 {
		return errors.New("Privilege not found for the specified role")
	}

	// Assign the privilege to the target role
	_, err = Config.Db.Exec("INSERT INTO role_privileges (role_id, privilege_id) VALUES (?, ?)", targetID, privilegeID)
	if err != nil {
		return err
	}

	return nil
}
