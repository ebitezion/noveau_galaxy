package rbac_test

import (
	"testing"

	"github.com/ebitezion/backend-framework/internal/rbac"
)

func TestRBAC(t *testing.T) {
	// Initialize a new RBAC instance
	rbacSystem := rbac.NewRBAC()

	// Add roles with associated privileges
	rbacSystem.AddRole("Admin", []rbac.Privilege{"Create", "Read", "Update", "Delete"})
	rbacSystem.AddRole("User", []rbac.Privilege{"Read"})

	// Add users to the RBAC system
	adminUser := rbac.User{Username: "admin", Password: "adminpass", Roles: []rbac.Role{"Admin"}}
	userUser := rbac.User{Username: "user", Password: "userpass", Roles: []rbac.Role{"User"}}

	err := rbacSystem.AddUser(adminUser)
	if err != nil {
		t.Errorf("Error adding admin user: %v", err)
	}

	err = rbacSystem.AddUser(userUser)
	if err != nil {
		t.Errorf("Error adding user user: %v", err)
	}

	// Save RBAC configuration to a file
	if err := rbacSystem.Save("rbac_test_config.json"); err != nil {
		t.Errorf("Error saving RBAC configuration: %v", err)
	}

	// Load RBAC configuration from the file
	loadedRBACSystem, err := rbac.Load("rbac_test_config.json")
	if err != nil {
		t.Errorf("Error loading RBAC configuration: %v", err)
	} else {
		// Check if the loaded RBAC system has the privilege for a user
		user := rbac.User{Username: "admin", Password: "adminpass", Roles: []rbac.Role{"Admin"}}
		permission := "Read"
		hasPermission := loadedRBACSystem.CheckPermission(user, rbac.Privilege(permission))
		expectedPermission := true // Change this according to your expectation
		if hasPermission != expectedPermission {
			t.Errorf("Expected permission '%s' for user 'admin', got %v", permission, hasPermission)
		}
	}
}
