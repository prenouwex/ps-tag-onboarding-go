package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mocksRepo "github.com/wexinc/ps-tag-onboarding-go/internal/mocks/repository"
	mocks "github.com/wexinc/ps-tag-onboarding-go/internal/mocks/service"
	"github.com/wexinc/ps-tag-onboarding-go/internal/model"
	"github.com/wexinc/ps-tag-onboarding-go/internal/utils"
	"net/http"
	"testing"
)

// Test use cases on IUserService class using mockery mocks
func TestUserService_Mockery_GetUser(t *testing.T) {
	// Given
	user := model.User{
		Id:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@gmail.com",
		Age:       30,
	}
	mockUserRepository := new(mocksRepo.IUserRepository)
	mockUserRepository.On("DbGetUser", mock.AnythingOfType("int64")).Return(&user, nil)
	mockUserValidationService := new(mocks.IUserValidationService)
	userService := UserService{mockUserRepository, mockUserValidationService}

	// When
	userRet, err := userService.GetUser(1)

	// Then
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, userRet.Id)
	assert.EqualValues(t, "John", userRet.FirstName)
	assert.EqualValues(t, "Doe", userRet.LastName)
}

// Test the not found functionality
func TestUserService_Mockery_GetUserNotFoundID(t *testing.T) {
	// Given
	mockUserRepository := new(mocksRepo.IUserRepository)
	mockUserRepository.On("DbGetUser", mock.AnythingOfType("int64")).Return(
		nil,
		utils.InternalServerError("the id is not found"))
	mockUserValidationService := new(mocks.IUserValidationService)
	userService := UserService{mockUserRepository, mockUserValidationService}

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
func TestMessagesService_Mockery_SaveUser_Success(t *testing.T) {
	// Given
	user := model.User{
		Id:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@gmail.com",
		Age:       30,
	}
	mockUserRepository := new(mocksRepo.IUserRepository)
	mockUserRepository.On("DbCreateUser", mock.Anything).Return(
		&user,
		nil)
	mockUserValidationService := new(mocks.IUserValidationService)
	mockUserValidationService.On("ValidateUser", mock.Anything).Return(nil)
	userService := UserService{mockUserRepository, mockUserValidationService}
	request := &model.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@gmail.com",
		Age:       30,
	}

	// When
	userRet, err := userService.SaveUser(request)

	// Then
	assert.NotNil(t, userRet)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, userRet.Id)
	assert.EqualValues(t, "John", userRet.FirstName)
	assert.EqualValues(t, "Doe", userRet.LastName)
	assert.EqualValues(t, "john.doe@gmail.com", userRet.Email)
	assert.EqualValues(t, 30, userRet.Age)
}

func TestUserService_Mockery_SaveUser_Invalid_Request(t *testing.T) {
	// Given
	user := model.User{
		Id:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@gmail.com",
		Age:       30,
	}
	mockUserRepository := new(mocksRepo.IUserRepository)
	mockUserRepository.On("DbCreateUser", mock.Anything).Return(
		&user,
		nil)
	mockUserValidationService := new(mocks.IUserValidationService)
	mockUserValidationService.On("ValidateUser", mock.Anything).Return([]string{"invalid_request"})
	userService := UserService{mockUserRepository, mockUserValidationService}

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
