package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
)

// Define a Permissions slice, which we will use to will hold the permission codes (like
// "account:read" and "account:write") for a single user.
// "account:read" and "account:write") for a single user.
type Permissions []string
type UserPermissions struct {
	ID         string
	Permission string
}

// Add a helper method to check whether the Permissions slice contains a specific
// permission code.
func (p Permissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}

// Define the PermissionModel type.
type PermissionModel struct {
	DB *sql.DB
}

// The GetAllForUser() method returns all permission codes for a specific user in a
// Permissions slice. The code in this method should feel very familiar --- it uses the
// standard pattern that we've already seen before for retrieving multiple data rows in
// an SQL query.
func (m PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	query := `
	SELECT permissions.code
	FROM permissions
	INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
	INNER JOIN accounts_auth ON users_permissions.user_id = accounts_auth.id
	WHERE accounts_auth.id = ?`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var permissions Permissions
	for rows.Next() {
		var permission string
		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return permissions, nil
}

// Add the provided permission codes for a specific user. Notice that we're using a
// variadic parameter for the codes so that we can assign multiple permissions in a
// single call.

func (m PermissionModel) AddForUser(userID int64, codes string) error {

	query := `
        INSERT INTO users_permissions
        SELECT ?, permissions.id FROM permissions WHERE permissions.code IN (?)
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userID, codes)

	if err != nil {
		// Check if the error is a MySQL duplicate entry error
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			// Duplicate entry error, return a custom error
			return errors.New("User already Has This Priviledge")
		}

		// Return the original error for other types of errors
		return err
	}

	return nil
}
func (m PermissionModel) GetAllPermissions() ([]UserPermissions, error) {
	query := `
	SELECT id, code
	FROM permissions
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []UserPermissions
	for rows.Next() {
		var permission UserPermissions
		err := rows.Scan(&permission.ID, &permission.Permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return permissions, nil
}
