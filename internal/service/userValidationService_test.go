package service

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/wexinc/ps-tag-onboarding-go/internal/model"
	"strings"
	"testing"
)

func TestAgeInvalid(t *testing.T) {
	// Given
	js := `{"id":2,"first_name":"Zenia","last_name":"Brennan","email":"ultrices.vivamus.rhoncus@yahoo.ca","age":9}`
	userValidation := UserValidationService{&MockRepo{}}
	var user model.User
	if err := json.Unmarshal([]byte(js), &user); err != nil {
		t.Errorf("failed to unmarshal lead to JSON: %v", err.Error())
	}

	// When
	validationErrors, err := userValidation.ValidateUser(&user)
	fmt.Println("err", err)
	fmt.Println("validationErrors", validationErrors)

	// Then
	assert.True(t, strings.Contains(strings.Join(validationErrors, ","), ERROR_AGE_MINIMUM))
	if validationErrors == nil {
		t.Errorf("expected validation error, none received")
	}
}

func TestEmailInvalid(t *testing.T) {
	// Given
	js := `{"id":2,"first_name":"Zenia","last_name":"Brennan","email":"bad_email","age":19}`
	userValidation := UserValidationService{&MockRepo{}}
	var user model.User
	if err := json.Unmarshal([]byte(js), &user); err != nil {
		t.Errorf("failed to unmarshal lead to JSON: %v", err.Error())
	}

	// When
	validationErrors, err := userValidation.ValidateUser(&user)
	fmt.Println("err", err)
	fmt.Println("validationErrors", validationErrors)

	// Then
	assert.True(t, strings.Contains(strings.Join(validationErrors, ","), ERROR_EMAIL_FORMAT))
	if validationErrors == nil {
		t.Errorf("expected validation error, none received")
	}
}

func TestFirstNameLastNameAlreadyExists(t *testing.T) {
	// Given
	js := `{"id":2,"first_name":"John","last_name":"Doe","email":"john.doe_t@gmail.com","age":19}`
	userValidation := UserValidationService{&MockRepo{}}
	var user model.User
	if err := json.Unmarshal([]byte(js), &user); err != nil {
		t.Errorf("failed to unmarshal lead to JSON: %v", err.Error())
	}

	// When
	validationErrors, err := userValidation.ValidateUser(&user)
	fmt.Println("err", err)
	fmt.Println("validationErrors", validationErrors)

	// Then
	assert.True(t, strings.Contains(strings.Join(validationErrors, ","), ERROR_NAME_UNIQUE))
	if validationErrors == nil {
		t.Errorf("expected validation error, none received")
	}
}
