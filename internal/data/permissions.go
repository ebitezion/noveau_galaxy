package data

import (
	"context"
	"database/sql"
	"time"
)

// Define a Permissions slice, which we will use to will hold the permission codes (like
// "movies:read" and "movies:write") for a single user.
// "movies:read" and "movies:write") for a single user.
type Permissions []string

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
	INNER JOIN users ON users_permissions.user_id = users.id
	WHERE users.id = ?`
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
	// Query to retrieve permission id
	query := `
        SELECT id FROM permissions WHERE code = ?
    `

	var permissionID int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Scan the result into permissionID
	err := m.DB.QueryRowContext(ctx, query, codes).Scan(&permissionID)
	if err != nil {
		return err
	}

	// Query to insert into users_permissions
	query = `
        INSERT INTO users_permissions (permission_id, user_id) VALUES (?, ?)
    `

	// Arguments for the query
	args := []interface{}{permissionID, userID}

	// Use the existing context for the ExecContext call
	_, err = m.DB.ExecContext(ctx, query, args...)
	return err
}
