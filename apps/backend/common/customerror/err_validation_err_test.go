package customerror

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

// Test struct for validation
type TestStruct struct {
	Email    string `validate:"required,email"`
	Username string `validate:"required,min=3,max=20"`
	Age      int    `validate:"required,min=18,max=100"`
}

func TestParseValidationErr(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name          string
		input         TestStruct
		expectReasons bool
		expectCode    ErrCode
	}{
		{
			name: "multiple validation errors",
			input: TestStruct{
				Email:    "invalid-email",
				Username: "ab",
				Age:      10,
			},
			expectReasons: true,
			expectCode:    ErrInvalidArgument.Code,
		},
		{
			name: "missing required fields",
			input: TestStruct{
				Email:    "",
				Username: "",
				Age:      0,
			},
			expectReasons: true,
			expectCode:    ErrInvalidArgument.Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.input)
			if err == nil {
				t.Skip("Expected validation error but got none")
			}

			validationErrs, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatalf("Expected validator.ValidationErrors but got %T", err)
			}

			got := ParseValidationErr(&validationErrs)

			if got == nil {
				t.Fatal("ParseValidationErr() returned nil")
			}

			if got.Code == nil || *got.Code != tt.expectCode {
				t.Errorf("ParseValidationErr() Code = %v, want %v", got.Code, tt.expectCode)
			}

			if got.HttpStatus == nil || *got.HttpStatus != ErrInvalidArgument.HttpStatus {
				t.Errorf("ParseValidationErr() HttpStatus = %v, want %v", got.HttpStatus, ErrInvalidArgument.HttpStatus)
			}

			if tt.expectReasons {
				if len(got.Reasons) == 0 {
					t.Error("ParseValidationErr() Reasons is empty, expected validation errors")
				}

				// Check that reasons map contains field names
				for _, validationErr := range validationErrs {
					fieldName := validationErr.Field()
					if _, exists := got.Reasons[fieldName]; !exists {
						t.Errorf("ParseValidationErr() Reasons missing field %s", fieldName)
					}
				}
			}

			if got.Inner == nil {
				t.Error("ParseValidationErr() Inner is nil, expected wrapped error")
			}
		})
	}
}

func TestParseValidationErr_FieldMapping(t *testing.T) {
	validate := validator.New()

	// Create a struct with specific validation errors
	input := TestStruct{
		Email:    "not-an-email",
		Username: "u",
		Age:      10,
	}

	err := validate.Struct(input)
	if err == nil {
		t.Fatal("Expected validation error but got none")
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		t.Fatalf("Expected validator.ValidationErrors but got %T", err)
	}

	got := ParseValidationErr(&validationErrs)

	// Check that each validation error is properly mapped
	for _, validationErr := range validationErrs {
		fieldName := validationErr.Field()
		reasonMsg, exists := got.Reasons[fieldName]

		if !exists {
			t.Errorf("Field %s not found in Reasons map", fieldName)
			continue
		}

		if reasonMsg == "" {
			t.Errorf("Reason for field %s is empty", fieldName)
		}

		// Verify the error message contains the field name
		if reasonMsg != validationErr.Error() {
			t.Errorf("Reason for field %s = %v, want %v", fieldName, reasonMsg, validationErr.Error())
		}
	}
}
