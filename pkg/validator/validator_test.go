package validator_test

import (
	"providerHub/pkg/validator"
	"testing"
)

// Test struct to validate
type User struct {
	Username string `validate:"required,alpha"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,password"`
	Phone    string `validate:"phone"`
}

func TestStructValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   User
		wantErr bool
	}{
		{
			name: "valid input",
			input: User{
				Username: "JohnDoe",
				Email:    "john.doe@example.com",
				Password: "Password1",
				Phone:    "+1234567890",
			},
			wantErr: false,
		},
		{
			name: "missing required field",
			input: User{
				Email:    "john.doe@example.com",
				Password: "Password1",
				Phone:    "+1234567890",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			input: User{
				Username: "JohnDoe",
				Email:    "invalid-email",
				Password: "Password1",
				Phone:    "+1234567890",
			},
			wantErr: true,
		},
		{
			name: "invalid password",
			input: User{
				Username: "JohnDoe",
				Email:    "john.doe@example.com",
				Password: "password", // Missing upper case and digit
				Phone:    "+1234567890",
			},
			wantErr: true,
		},
		{
			name: "invalid phone",
			input: User{
				Username: "JohnDoe",
				Email:    "john.doe@example.com",
				Password: "Password1",
				Phone:    "invalid-phone", // Invalid phone format
			},
			wantErr: true,
		},
		{
			name: "empty phone (optional)",
			input: User{
				Username: "JohnDoe",
				Email:    "john.doe@example.com",
				Password: "Password1",
				Phone:    "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Struct(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Struct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
