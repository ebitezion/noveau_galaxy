// rbac.go

package rbac_2

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"
)

// Role represents a user role
type Role string

// Privilege represents a specific action a user can perform
type Privilege string

// User represents a user in the system
type User struct {
	Username string
	Password string
	Roles    []Role
}

// RBAC represents the Role-Based Access Control system
type RBAC struct {
	RolePrivileges map[Role][]Privilege
	Users          map[string]User // Map of username to User
	mutex          sync.RWMutex
}

// NewRBAC creates a new RBAC instance
func NewRBAC() *RBAC {
	return &RBAC{
		RolePrivileges: make(map[Role][]Privilege),
		Users:          make(map[string]User),
	}
}

// AddRole adds a new role with associated privileges
func (r *RBAC) AddRole(role Role, privileges []Privilege) {
	r.RolePrivileges[role] = privileges
}

// AddUser adds a new user to the RBAC system
func (r *RBAC) AddUser(user User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.Users[user.Username]; exists {
		return errors.New("user already exists")
	}

	r.Users[user.Username] = user
	return nil
}

// GetUserByUsername retrieves a user by username
func (r *RBAC) GetUserByUsername(username string) (User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.Users[username]
	if !exists {
		return User{}, errors.New("user not found")
	}

	return user, nil
}

// UpdateUser updates an existing user
func (r *RBAC) UpdateUser(username string, updatedUser User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.Users[username]; !exists {
		return errors.New("user not found")
	}

	r.Users[username] = updatedUser
	return nil
}

// DeleteUser removes a user from the RBAC system
func (r *RBAC) DeleteUser(username string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.Users[username]; !exists {
		return errors.New("user not found")
	}

	delete(r.Users, username)
	return nil
}

// UpdateRolePrivileges updates privileges for a role
func (r *RBAC) UpdateRolePrivileges(role Role, updatedPrivileges []Privilege) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.RolePrivileges[role]; !exists {
		return errors.New("role not found")
	}

	r.RolePrivileges[role] = updatedPrivileges
	return nil
}

// DeleteRole removes a role from the RBAC system
func (r *RBAC) DeleteRole(role Role) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.RolePrivileges[role]; !exists {
		return errors.New("role not found")
	}

	delete(r.RolePrivileges, role)
	return nil
}

// CheckPermission checks if a user with certain roles has the required privilege
func (r *RBAC) CheckPermission(user User, privilege Privilege) bool {
	// Retrieve user roles from the RBAC system
	userRoles, err := r.GetUserByUsername(user.Username)
	if err != nil {
		return false // User not found or error occurred
	}

	for _, role := range userRoles.Roles {
		// Retrieve privileges for each role from the RBAC system
		privileges, exists := r.RolePrivileges[role]
		if exists {
			// Check if the required privilege exists for the user's role
			for _, p := range privileges {
				if p == privilege {
					return true // User has the required privilege
				}
			}
		}
	}

	return false // User doesn't have the required privilege
}

// Save saves RBAC configuration to a JSON file
func (r *RBAC) Save(filename string) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	data, err := json.Marshal(r)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

// Load loads RBAC configuration from a JSON file
func Load(filename string) (*RBAC, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var loadedRBACSystem RBAC
	if err := json.Unmarshal(data, &loadedRBACSystem); err != nil {
		return nil, err
	}

	return &loadedRBACSystem, nil
}
