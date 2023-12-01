// rbac.go

package rbac

import (
	"errors"
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
	rolePrivileges map[Role][]Privilege
	users          map[string]User // Map of username to User
	mutex          sync.RWMutex
}

// NewRBAC creates a new RBAC instance
func NewRBAC() *RBAC {
	return &RBAC{
		rolePrivileges: make(map[Role][]Privilege),
		users:          make(map[string]User),
	}
}

// AddRole adds a new role with associated privileges
func (r *RBAC) AddRole(role Role, privileges []Privilege) {
	r.rolePrivileges[role] = privileges
}

// AddUser adds a new user to the RBAC system
func (r *RBAC) AddUser(user User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[user.Username]; exists {
		return errors.New("user already exists")
	}

	r.users[user.Username] = user
	return nil
}

// GetUserByUsername retrieves a user by username
func (r *RBAC) GetUserByUsername(username string) (User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[username]
	if !exists {
		return User{}, errors.New("user not found")
	}

	return user, nil
}

// UpdateUser updates an existing user
func (r *RBAC) UpdateUser(username string, updatedUser User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[username]; !exists {
		return errors.New("user not found")
	}

	r.users[username] = updatedUser
	return nil
}

// DeleteUser removes a user from the RBAC system
func (r *RBAC) DeleteUser(username string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[username]; !exists {
		return errors.New("user not found")
	}

	delete(r.users, username)
	return nil
}

// UpdateRolePrivileges updates privileges for a role
func (r *RBAC) UpdateRolePrivileges(role Role, updatedPrivileges []Privilege) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.rolePrivileges[role]; !exists {
		return errors.New("role not found")
	}

	r.rolePrivileges[role] = updatedPrivileges
	return nil
}

// DeleteRole removes a role from the RBAC system
func (r *RBAC) DeleteRole(role Role) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.rolePrivileges[role]; !exists {
		return errors.New("role not found")
	}

	delete(r.rolePrivileges, role)
	return nil
}

// CheckPermission checks if a user with certain roles has the required privilege
func (r *RBAC) CheckPermission(user User, privilege Privilege) bool {
	for _, role := range user.Roles {
		privileges, exists := r.rolePrivileges[role]
		if exists {
			for _, p := range privileges {
				if p == privilege {
					return true // User has the required privilege
				}
			}
		}
	}

	return false // User doesn't have the required privilege
}
