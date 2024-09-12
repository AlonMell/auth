package validator

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestStruct struct {
	A1 string `validate:"required,alpha"`
	A2 string `validate:"required,alphanum"`
	A3 string `validate:"required,email"`
	A4 string `validate:"required,phone"`
	A5 string `validate:"required,password"`
	A6 string `validate:"required,uuid"`
	B1 string `validate:"alpha"`
	B2 string
	B4 string `json:"b4" validate:"required,password"`
	B5 string `json:"b5" validate:""`
}

func TestStructValidation(t *testing.T) {
	tests := []struct {
		name  string
		input TestStruct
		err   error
	}{
		{
			name: "Valid data",
			input: TestStruct{
				A1: "ValidName",
				A2: "Valid123",
				A3: "test@example.com",
				A4: "+1234567890",
				A5: "StrongPass1",
				A6: uuid.New().String(),
				B1: "OptionalAlpha",
				B2: "OptionalField",
				B4: "AnotherPass1",
				B5: "Optional",
			},
			err: nil,
		},
		{
			name: "Missing required fields",
			input: TestStruct{
				A1: "",
				A2: "",
				A3: "",
				A4: "",
				A5: "",
				A6: "",
			},
			err: ErrRequiredField,
		},
		{
			name: "Invalid email format",
			input: TestStruct{
				A1: "ValidName",
				A2: "Valid123",
				A3: "invalid-email",
				A4: "+1234567890",
				A5: "StrongPass1",
				A6: uuid.New().String(),
			},
			err: ErrNotValidValue,
		},
		{
			name: "Invalid phone format",
			input: TestStruct{
				A1: "ValidName",
				A2: "Valid123",
				A3: "test@example.com",
				A4: "123456", // Invalid phone
				A5: "StrongPass1",
				A6: uuid.New().String(),
			},
			err: ErrNotValidValue,
		},
		{
			name: "Invalid password (missing digit)",
			input: TestStruct{
				A1: "ValidName",
				A2: "Valid123",
				A3: "test@example.com",
				A4: "+1234567890",
				A5: "StrongPass", // Missing digit
				A6: uuid.New().String(),
			},
			err: ErrNotValidPassword,
		},
		{
			name: "Invalid UUID format",
			input: TestStruct{
				A1: "ValidName",
				A2: "Valid123",
				A3: "test@example.com",
				A4: "+1234567890",
				A5: "StrongPass1",
				A6: "invalid-uuid",
			},
			err: ErrNotValidValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Struct(tt.input)

			if assert.Errorf(t, err, tt.err.Error()) {
				assert.True(t, errors.Is(err, tt.err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
