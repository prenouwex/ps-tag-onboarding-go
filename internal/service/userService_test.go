package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/wexinc/ps-tag-onboarding-go/internal/model"
	"github.com/wexinc/ps-tag-onboarding-go/internal/repository"
	"github.com/wexinc/ps-tag-onboarding-go/internal/utils"
	"net/http"
	"testing"
)

func TestUserService_GetUser(t *testing.T) {
	// Given
	var repo repository.IUserRepository = &MockRepo{}
	var userValidation IUserValidationService = &MockValidation{}
	userService := UserService{repo, userValidation}
	getUserDomain = func(userId int64) (*model.User, utils.MessageErr) {
		return &model.User{
			Id:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@gmail.com",
			Age:       30,
		}, nil
	}

	// When
	user, err := userService.GetUser(1)

	// Then
	fmt.Println("this is the message: ", user)
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "John", user.FirstName)
	assert.EqualValues(t, "Doe", user.LastName)
}

// Test the not found functionality
func TestUserService_GetUserNotFoundID(t *testing.T) {
	// Given
	var repo repository.IUserRepository = &MockRepo{}
	var userValidation IUserValidationService = &MockValidation{}
	userService := UserService{repo, userValidation}
	getUserDomain = func(userId int64) (*model.User, utils.MessageErr) {
		return nil, utils.InternalServerError("the id is not found")
	}

	// When
	user, err := userService.GetUser(1)

	// Then
	assert.Nil(t, user)
	assert.NotNil(t, err)
	//assert.EqualValues(t, http.StatusNotFound, err.Status())
	assert.EqualValues(t, "the id is not found", err.Message())
	//assert.EqualValues(t, "not_found", err.Error())
}

///////////////////////////////////////////////////////////////
// 				"CreateMessage" test cases
///////////////////////////////////////////////////////////////

// Here we call the domain method, so we must mock it
func TestMessagesService_SaveUser_Success(t *testing.T) {
	// Given
	var repo repository.IUserRepository = &MockRepo{}
	var userValidation IUserValidationService = &MockValidation{}
	userService := UserService{repo, userValidation}
	createUserDomain = func(user *model.User) (*model.User, utils.MessageErr) {
		return &model.User{
			Id:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@gmail.com",
			Age:       30,
		}, nil
	}
	getValidation = func(user *model.User) []string {
		return nil
	}
	request := &model.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@gmail.com",
		Age:       30,
	}

	// When
	user, err := userService.SaveUser(request)

	// Then
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "John", user.FirstName)
	assert.EqualValues(t, "Doe", user.LastName)
	assert.EqualValues(t, "john.doe@gmail.com", user.Email)
	assert.EqualValues(t, 30, user.Age)
}

func TestUserService_SaveUser_Invalid_Request(t *testing.T) {
	// Given
	var repo repository.IUserRepository = &MockRepo{}
	var userValidation IUserValidationService = &MockValidation{}
	userService := UserService{repo, userValidation}
	getValidation = func(user *model.User) []string {
		return []string{"invalid_request"}
	}
	tests := []struct {
		request    *model.User
		statusCode int
		errMsg     string
		errErr     string
	}{
		{
			request: &model.User{
				FirstName: "",
				LastName:  "Doe",
				Email:     "john.doe@gmail.com",
				Age:       30,
			},
			statusCode: http.StatusBadRequest,
			errMsg:     "invalid_request",
			errErr:     "bad_request",
		},
		{
			request: &model.User{
				FirstName: "John",
				LastName:  "",
				Email:     "john.doe@gmail.com",
				Age:       30,
			},
			statusCode: http.StatusBadRequest,
			errMsg:     "invalid_request",
			errErr:     "bad_request",
		},
	}
	for _, tt := range tests {
		// When
		msg, err := userService.SaveUser(tt.request)

		// Then
		assert.Nil(t, msg)
		assert.NotNil(t, err)
		assert.EqualValues(t, tt.errMsg, err.Message())
		assert.EqualValues(t, tt.statusCode, err.Status())
		assert.EqualValues(t, tt.errErr, err.Error())
	}
}

///////////////////////////////////////////////////////////////
// 				"UpdateUser" test cases
///////////////////////////////////////////////////////////////

func TestUserService_UpdateUser_Success(t *testing.T) {
	// Given
	var repo repository.IUserRepository = &MockRepo{}
	var userValidation IUserValidationService = &MockValidation{}
	userService := UserService{repo, userValidation}
	getUserDomain = func(userId int64) (*model.User, utils.MessageErr) {
		return &model.User{
			Id:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@gmail.com",
			Age:       30,
		}, nil
	}
	updateUserDomain = func(user *model.User) (*model.User, utils.MessageErr) {
		return &model.User{
			Id:        1,
			FirstName: "Johnny",
			LastName:  "Dover",
			Email:     "john.doe@gmail.com",
			Age:       30,
		}, nil
	}
	getValidation = func(user *model.User) []string {
		return nil
	}
	request := &model.User{
		FirstName: "Johnny",
		LastName:  "Dover",
		Email:     "john.doe@gmail.com",
		Age:       30,
	}

	// When
	user, err := userService.UpdateUser(request)

	// Then
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "Johnny", user.FirstName)
	assert.EqualValues(t, "Dover", user.LastName)
}

// Test error scenarios where an error can occur when trying to fetch the user to update,
// anything from a timeout error to a not found error.
func TestUserService_UpdateUser_Failure_Getting_Former_Message(t *testing.T) {
	// Given
	var repo repository.IUserRepository = &MockRepo{}
	var userValidation IUserValidationService = &MockValidation{}
	userService := UserService{repo, userValidation}
	getUserDomain = func(userId int64) (*model.User, utils.MessageErr) {
		return nil, utils.InternalServerError("error getting message")
	}
	getValidation = func(user *model.User) []string {
		return nil
	}
	request := &model.User{
		FirstName: "Johnny",
		LastName:  "Dover",
		Email:     "john.doe@gmail.com",
		Age:       30,
	}

	// When
	msg, err := userService.UpdateUser(request)

	// Then
	assert.Nil(t, msg)
	assert.NotNil(t, err)
	assert.EqualValues(t, "error getting message", err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "server_error", err.Error())
}

///////////////////////////////////////////////////////////////
// 				"DeleteUser" test cases
///////////////////////////////////////////////////////////////

func TestUserService_DeleteUser_Success(t *testing.T) {
	// Given
	var repo repository.IUserRepository = &MockRepo{}
	var userValidation IUserValidationService = &MockValidation{}
	userService := UserService{repo, userValidation}
	getUserDomain = func(userId int64) (*model.User, utils.MessageErr) {
		return &model.User{
			Id:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@gmail.com",
			Age:       30,
		}, nil
	}
	deleteUserDomain = func(userId int64) (*model.User, utils.MessageErr) {
		return nil, nil
	}
	getValidation = func(user *model.User) []string {
		return nil
	}

	// When
	_, err := userService.DeleteUser(1)

	// Then
	assert.Nil(t, err)
}

// It can range from a 500 error to a 404 error, we didnt mock deleting the message because we will not get there
func TestMessagesService_DeleteMessage_Error_Getting_Message(t *testing.T) {
	// Given
	var repo repository.IUserRepository = &MockRepo{}
	var userValidation IUserValidationService = &MockValidation{}
	userService := UserService{repo, userValidation}
	getUserDomain = func(userId int64) (*model.User, utils.MessageErr) {
		return nil, utils.InternalServerError("Something went wrong getting message")
	}
	getValidation = func(user *model.User) []string {
		return nil
	}

	// When
	_, err := userService.DeleteUser(1)

	// Then
	assert.NotNil(t, err)
	assert.EqualValues(t, "Something went wrong getting message", err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "server_error", err.Error())
}

///////////////////////////////////////////////////////////////
// 				"DeleteUser" test cases
///////////////////////////////////////////////////////////////

func TestUserService_GetAllUser(t *testing.T) {
	// Given
	var repo repository.IUserRepository = &MockRepo{}
	var userValidation IUserValidationService = &MockValidation{}
	userService := UserService{repo, userValidation}
	getAllUsersDomain = func() ([]model.User, utils.MessageErr) {
		return []model.User{
			{
				Id:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@gmail.com",
				Age:       30,
			},
			{
				Id:        2,
				FirstName: "Johnny",
				LastName:  "Dover",
				Email:     "john.doe@gmail.com",
				Age:       30,
			},
		}, nil
	}

	// When
	users, err := userService.GetAllUsers()

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, users)
	assert.EqualValues(t, users[0].Id, 1)
	assert.EqualValues(t, users[0].FirstName, "John")
	assert.EqualValues(t, users[0].LastName, "Doe")
	assert.EqualValues(t, users[1].Id, 2)
	assert.EqualValues(t, users[1].FirstName, "Johnny")
	assert.EqualValues(t, users[1].LastName, "Dover")
}

func TestUserService_GetAllUsers_Error_Getting_Users(t *testing.T) {
	// Given
	var repo repository.IUserRepository = &MockRepo{}
	var userValidation IUserValidationService = &MockValidation{}
	userService := UserService{repo, userValidation}
	getAllUsersDomain = func() ([]model.User, utils.MessageErr) {
		return nil, utils.InternalServerError("error getting messages")
	}
	messages, err := userService.GetAllUsers()
	assert.NotNil(t, err)
	assert.Nil(t, messages)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "error getting messages", err.Message())
	assert.EqualValues(t, "server_error", err.Error())
}

// =================================================== //
// ================ Mock Declaration ================= //
// =================================================== //

var (
	getUserDomain    func(userId int64) (*model.User, utils.MessageErr)
	createUserDomain func(user *model.User) (*model.User, utils.MessageErr)
	//createMessageDomain  func(msg *domain.Message) (*domain.Message, error_utils.MessageErr)
	updateUserDomain  func(user *model.User) (*model.User, utils.MessageErr)
	deleteUserDomain  func(userId int64) (*model.User, utils.MessageErr)
	getAllUsersDomain func() ([]model.User, utils.MessageErr)

	getValidation func(user *model.User) []string
)

// MockRepo is a struct that mocks UserRepository.
type MockRepo struct{}

// Repository mock method implementation.
func (m *MockRepo) DbListUsers() ([]model.User, utils.MessageErr) {
	// Implement your mock behavior here
	return getAllUsersDomain() // Return a mock GORM DB
}
func (m *MockRepo) DbCreateUser(user *model.User) (*model.User, utils.MessageErr) {
	// Implement your mock behavior here
	return createUserDomain(user) // Return a mock GORM DB
}
func (m *MockRepo) DbGetUser(id int64) (*model.User, utils.MessageErr) {
	// Implement your mock behavior here
	return getUserDomain(id) // Return a mock GORM DB
}
func (m *MockRepo) DbUpdateUser(user *model.User) (*model.User, utils.MessageErr) {
	// Implement your mock behavior here
	return updateUserDomain(user) // Return a mock GORM DB
}
func (m *MockRepo) DbDeleteUser(id int64) (*model.User, utils.MessageErr) {
	// Implement your mock behavior here
	return deleteUserDomain(id) // Return a mock GORM DB
}
func (m *MockRepo) ExistsByFirstNameAndLastName(firstName string, lastName string) (bool, utils.MessageErr) {
	// Implement your mock behavior here
	return true, nil // Return a mock GORM DB
}

type MockValidation struct{}

func (m *MockValidation) ValidateUser(user *model.User) []string {
	return getValidation(user)
}

// =================================================== //
