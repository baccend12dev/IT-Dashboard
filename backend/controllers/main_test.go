package controllers_test

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"testing"
)

// generateTestToken seeds a test user with the specified role and returns a valid JWT token
func generateTestToken(t *testing.T, role string) string {
	hashedPassword, err := utils.HashPassword("testpassword")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	user := models.User{
		Username: "test_" + role,
		Password: hashedPassword,
		Role:     role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		t.Fatalf("failed to generate JWT token: %v", err)
	}

	return token
}
