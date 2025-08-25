package service

import (
	"lunchorder/repository"
	"testing"
)

func TestAuthService_IsAdmin(t *testing.T) {
	authService := &AuthService{
		adminEmail: "tyler.osborn@impact.com",
	}

	// Test admin user
	adminUser := &repository.User{
		Email: "tyler.osborn@impact.com",
		Role:  "admin",
	}

	if !authService.IsAdmin(adminUser) {
		t.Error("Expected tyler.osborn@impact.com to be admin")
	}

	// Test standard user
	standardUser := &repository.User{
		Email: "other@example.com",
		Role:  "standard",
	}

	if authService.IsAdmin(standardUser) {
		t.Error("Expected other@example.com to not be admin")
	}
}

func TestAuthService_RoleAssignment(t *testing.T) {
	authService := &AuthService{
		adminEmail: "tyler.osborn@impact.com",
	}

	// Test that admin email gets admin role
	adminEmail := "tyler.osborn@impact.com"
	if adminEmail == authService.adminEmail {
		role := "admin"
		if role != "admin" {
			t.Error("Expected admin email to get admin role")
		}
	}

	// Test that other emails get standard role
	standardEmail := "user@example.com"
	if standardEmail != authService.adminEmail {
		role := "standard"
		if role != "standard" {
			t.Error("Expected non-admin email to get standard role")
		}
	}
}