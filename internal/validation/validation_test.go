package validation

import (
	"context"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type customValidationTests struct {
	have any
	want bool
	name string
}

var testValidator *Validator

func TestMain(m *testing.M) {
	ctx := context.Background()
	validate := validator.New(validator.WithRequiredStructEnabled())
	testValidator = New(ctx, validate).(*Validator)
	m.Run()
}

func TestPasswordTagValidation(t *testing.T) {
	passwordTests := []customValidationTests{
		{
			have: "123",
			want: false,
			name: "failure case: Testing password validation with a < 8 character only numbers password",
		},
		{
			have: "@##$#",
			want: false,
			name: "failure case: Testing password validation with a < 8 character only symbols password",
		},
		{
			have: "ABCDE",
			want: false,
			name: "failure case: Testing password validation with a < 8 character only uppercased letters password",
		},
		{
			have: "abc",
			want: false,
			name: "failure case: Testing password validation with a < 8 character only lowercased letters password",
		},
		{
			have: "123456789",
			want: false,
			name: "failure case: Testing password validation with a 8+ character only numbers password",
		},
		{
			have: "@@@$$$%^^#$#@$#",
			want: false,
			name: "failure case: Testing password validation with a 8+ character only symbols password",
		},
		{
			have: "EUSOUASENHATESTELONGA",
			want: false,
			name: "failure case: Testing password validation with a 8+ character only uppercased letters password",
		},
		{
			have: "eusouasenhatestelonga",
			want: false,
			name: "failure case: Testing password validation with a 8+ character only lowercased letters password",
		},
		{
			have: "Testando123",
			want: false,
			name: "failure case: Testing password validation with a 8+ character no symbols password",
		},
		{
			have: "testando123@",
			want: false,
			name: "failure case: Testing password validation with a 8+ character no uppercased letters password",
		},
		{
			have: "TESTANDO123@",
			want: false,
			name: "failure case: Testing password validation with a 8+ character no lowercased letters password",
		},
		{
			have: "TESTanDO@#$#",
			want: false,
			name: "failure case: Testing password validation with a 8+ character no numbers password",
		},
		{
			have: "Testando123@",
			want: true,
			name: "success case: Testing a valid password with 8+ characters and at least one symbol, one uppercased letter, one lowercased letter and one number",
		},
	}

	for _, testCase := range passwordTests {
		t.Logf("Running passwordTagValidation %s\n", testCase.name)
		err := testValidator.validate.Var(testCase.have, "password")

		if testCase.want {
			assert.NoError(t, err, "Unexpected error for testCase: %v", testCase)
		} else {
			assert.Error(t, err, "Expected error for testCase: %v", testCase)
		}
	}
}

func TestStructValidation(t *testing.T) {
	type arguments struct {
		validate ValidateProvider
		data     any
	}

	structTests := []struct {
		name      string
		arguments arguments
		want      []errorResponse
	}{
		{
			name: "success case: Testing struct validation using all utilized tags in the application",
			arguments: arguments{
				validate: testValidator.validate,
				data: struct {
					Name     string `json:"name" validate:"required,min=3"`
					Email    string `json:"email" validate:"required,email"`
					Password string `json:"password" validate:"required,password"`
				}{
					Name:     "Testing name",
					Email:    "testing@testing.com",
					Password: "Testing123@123",
				},
			},
			want: []errorResponse{},
		},
		{
			name: "failure case: Testing struct validation errors using all utilized tags in the application",
			arguments: arguments{
				validate: testValidator.validate,
				data: struct {
					Name     string `json:"name" validate:"required,min=3"`
					Email    string `json:"email" validate:"required,email"`
					Password string `json:"password" validate:"required,password"`
				}{
					Name:     "T",
					Email:    "testing",
					Password: "test",
				},
			},
			want: []errorResponse{
				{
					Error:        true,
					FailedField:  "name",
					Tag:          "min",
					ErrorMessage: "field 'name' must be at least 3 characters long",
				},
				{
					Error:        true,
					FailedField:  "email",
					Tag:          "email",
					ErrorMessage: "field 'email' must be a valid email address",
				},
				{
					Error:        true,
					FailedField:  "password",
					Tag:          "password",
					ErrorMessage: "field 'password' must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character",
				},
			},
		},
	}

	for _, testCase := range structTests {
		t.Run("", func(t *testing.T) {
			t.Logf("Running structValidation %s\n", testCase.name)
			response := structValidation(testCase.arguments.validate, testCase.arguments.data)

			assert.ElementsMatch(t, testCase.want, response, "error lists do not match")
		})
	}
}
