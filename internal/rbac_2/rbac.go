// rbac.go
package rbac_2

import (
	"errors"
	"sync"

	"github.com/ebitezion/backend-framework/internal/configuration"
)

// Role represents a user role
type Role string

// Privilege represents a specific action a user can perform
type Privilege string

// User represents a user in the system
type User struct {
	Username string
	Password string
	Role     Role
}

var Config configuration.Configuration

func SetConfig(config *configuration.Configuration) {
	Config = *config
}

// RBAC represents the Role-Based Access Control system
type RBAC struct {
	Roles       map[Role][]Privilege // Map of role to privileges
	Users       map[string]User      // Map of username to User
	mutex       sync.RWMutex
	defaultRole Role // Default role assigned to new users
	// Database    *YourDatabaseType // Database connection
}

// NewRBACWithDB creates a new RBAC instance with a database connection
func NewRBACWithDB() *RBAC {
	return &RBAC{
		Roles:       make(map[Role][]Privilege),
		Users:       make(map[string]User),
		defaultRole: "user", // Set the default role for new users
		//	Database:    db,
	}
}

// SetDefaultRole sets the default role for new users
func (r *RBAC) SetDefaultRole(role Role) {
	r.defaultRole = role
}

// LoadUsersFromDB loads users from the database into the RBAC system
func (r *RBAC) LoadUsersFromDB() error {
	// Query users from the database and populate the Users map
	// users, err := Config.Db.Query()
	// if err != nil {
	//     return err
	// }
	// for _, user := range users {
	//     r.Users[user.Username] = user
	// }
	return nil
}

// AddRole adds a new role with associated privileges
func (r *RBAC) AddRole(role Role, privileges []Privilege) {
	r.Roles[role] = privileges
}

// AddUser adds a new user to the RBAC system and the database
func (r *RBAC) AddUser(user User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.Users[user.Username]; exists {
		return errors.New("user already exists")
	}

	// Save user to the database
	// Example: err := r.Database.SaveUser(user)
	// if err != nil {
	//     return err
	// }

	r.Users[user.Username] = user
	return nil
}

// GetUserByUsername retrieves a user by username from the RBAC system
func (r *RBAC) GetUserByUsername(username string) (User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// user, exists := r.Users[username]
	// if !exists {
	// 	return User{}, errors.New("user not found")
	// }
	type Request struct {
		Username string
		Role     string
	}

	// query := "SELECT Username, Role FROM `accounts_auth` WHERE `username` = ?"

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// //err := Config.Db.QueryRowContext(ctx, query, username).Scan(&accountNumber)
	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return errors.New("accounts.GetUserByUsername: Account not found")
	// 	}
	// 	return err
	// }

	return User{}, nil
}

// UpdateUser updates an existing user in the RBAC system and the database
func (r *RBAC) UpdateUser(username string, updatedUser User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.Users[username]; !exists {
		return errors.New("user not found")
	}

	// Update user in the database
	// Example: err := r.Database.UpdateUser(username, updatedUser)
	// if err != nil {
	//     return err
	// }

	r.Users[username] = updatedUser
	return nil
}

// DeleteUser removes a user from the RBAC system and the database
func (r *RBAC) DeleteUser(username string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.Users[username]; !exists {
		return errors.New("user not found")
	}

	// Delete user from the database
	// Example: err := r.Database.DeleteUser(username)
	// if err != nil {
	//     return err
	// }

	delete(r.Users, username)
	return nil
}

// CheckPermission checks if a user has the required privilege based on their role
func (r *RBAC) CheckPermission(username string, privilege Privilege) bool {
	user, err := r.GetUserByUsername(username)
	if err != nil {
		return false // User not found or error occurred
	}

	// Retrieve user's role and associated privileges
	userRole := user.Role
	privileges, exists := r.Roles[userRole]
	if !exists {
		return false // Role not found
	}

	// Check if the required privilege exists for the user's role
	for _, p := range privileges {
		if p == privilege {
			return true // User has the required privilege
		}
	}

	return false // User doesn't have the required privilege
}
