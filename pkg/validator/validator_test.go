package validator

import (
	"testing"
)

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	IsActive bool   `json:"is_active"`
}

// Определение структуры для тестов
type TestStruct struct {
	Username string `validate:"required,alpha"`
	Email    string `validate:"required,email"`
	Phone    string `validate:"phone"`
	Password string `validate:"required,password"`
	UUID     string `validate:"uuid"`
}

func TestStruct2Validation(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateUserRequest
		wantErr bool
	}{
		{
			name: "valid input",
			input: CreateUserRequest{
				Email:    "test123@email.com",
				Password: "Password123",
				IsActive: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Struct(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Struct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStructValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   TestStruct
		wantErr bool
	}{
		{
			name: "missing required field",
			input: TestStruct{
				Email:    "john.doe@example.com",
				Phone:    "+1234567890",
				Password: "Password123",
				UUID:     "b00e1a65-f342-4496-bddc-acd438174c8d",
			},
			wantErr: true,
		},
		{
			name: "invalid email format",
			input: TestStruct{
				Username: "JohnDoe",
				Email:    "john.doe@example",
				Phone:    "+1234567890",
				Password: "Password123",
				UUID:     "b00e1a65-f342-4496-bddc-acd438174c8d",
			},
			wantErr: true,
		},
		{
			name: "invalid phone number",
			input: TestStruct{
				Username: "JohnDoe",
				Email:    "john.doe@example.com",
				Phone:    "123-456-7890",
				Password: "Password123",
				UUID:     "b00e1a65-f342-4496-bddc-acd438174c8d",
			},
			wantErr: true,
		},
		{
			name: "invalid password",
			input: TestStruct{
				Username: "JohnDoe",
				Email:    "john.doe@example.com",
				Phone:    "+1234567890",
				Password: "pass123",
				UUID:     "b00e1a65-f342-4496-bddc-acd438174c8d",
			},
			wantErr: true,
		},
		{
			name: "invalid UUID format",
			input: TestStruct{
				Username: "JohnDoe",
				Email:    "john.doe@example.com",
				Phone:    "+1234567890",
				Password: "Password123",
				UUID:     "invalid-uuid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Struct(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Struct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
