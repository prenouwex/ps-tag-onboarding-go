package service

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/wexinc/ps-tag-onboarding-go/internal/log"
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
	log.Info.Println("err", err)
	log.Info.Println("validationErrors", validationErrors)

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
	log.Info.Println("err", err)
	log.Info.Println("validationErrors", validationErrors)

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
	log.Info.Println("err", err)
	log.Info.Println("validationErrors", validationErrors)

	// Then
	assert.True(t, strings.Contains(strings.Join(validationErrors, ","), ERROR_NAME_UNIQUE))
	if validationErrors == nil {
		t.Errorf("expected validation error, none received")
	}
}

//// ================ Mock Declaration ================= //
//// =================================================== //
//// =================================================== //
//
//// MockRepo is a struct that mocks UserRepository.
//type MockRepo struct{}
//
//// Repository mock method implementation.
//func (m *MockRepo) DbListUsers() ([]*model.User, error) {
//	// Implement your mock behavior here
//	return nil, nil // Return a mock GORM DB
//}
//func (m *MockRepo) DbCreateUser(user *model.User) (int64, error) {
//	// Implement your mock behavior here
//	return 0, nil // Return a mock GORM DB
//}
//func (m *MockRepo) DbGetUser(id int64) (*model.User, error) {
//	// Implement your mock behavior here
//	return nil, nil // Return a mock GORM DB
//}
//func (m *MockRepo) DbUpdateUser(user *model.User) (*model.User, error) {
//	// Implement your mock behavior here
//	return nil, nil // Return a mock GORM DB
//}
//func (m *MockRepo) DbDeleteUser(id int64) (*model.User, error) {
//	// Implement your mock behavior here
//	return nil, nil // Return a mock GORM DB
//}
//func (m *MockRepo) ExistsByFirstNameAndLastName(firstName string, lastName string) (bool, error) {
//	// Implement your mock behavior here
//	return true, nil // Return a mock GORM DB
//}
//
//// =================================================== //
