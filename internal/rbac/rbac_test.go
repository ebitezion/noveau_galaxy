package rbac_test

import (
	"testing"

	"github.com/ebitezion/backend-framework/internal/rbac"
)

func TestRBAC(t *testing.T) {
	// Initialize RBAC
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

	// Retrieve and update a user
	userToUpdate, err := rbacSystem.GetUserByUsername("admin")
	if err != nil {
		t.Errorf("Error retrieving user: %v", err)
	}

	userToUpdate.Password = "newadminpass"
	err = rbacSystem.UpdateUser("admin", userToUpdate)
	if err != nil {
		t.Errorf("Error updating user: %v", err)
	}

	updatedUser, err := rbacSystem.GetUserByUsername("admin")//
	t.Log("Got",updatedUser,"Expected", )
	if err != nil {
		t.Errorf("Error retrieving updated user: %v", err)
	}

	// Check permissions for users
	adminPermissions := rbacSystem.CheckPermission(adminUser, "Delete")
	if !adminPermissions {
		t.Error("Admin should have permission to delete")
	}

	userPermissions := rbacSystem.CheckPermission(userUser, "Update")
	if userPermissions {
		t.Error("User shouldn't have permission to update")
	}

	// Update role privileges
	err = rbacSystem.UpdateRolePrivileges("Admin", []rbac.Privilege{"Create", "Read", "Update"})
	if err != nil {
		t.Errorf("Error updating role privileges: %v", err)
	}

	// Delete a user
	err = rbacSystem.DeleteUser("user")
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}

	_, err = rbacSystem.GetUserByUsername("user")
	if err == nil {
		t.Error("Deleted user shouldn't be found")
	}

	// Delete a role
	err = rbacSystem.DeleteRole("User")
	if err != nil {
		t.Errorf("Error deleting role: %v", err)
	}
}
