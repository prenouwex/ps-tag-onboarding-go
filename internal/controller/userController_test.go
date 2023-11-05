package controller

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/wexinc/ps-tag-onboarding-go/internal/model"
	//"github.com/wexinc/ps-tag-onboarding-go/internal/repository"
	"github.com/wexinc/ps-tag-onboarding-go/internal/service"
	"github.com/wexinc/ps-tag-onboarding-go/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	getUserService    func(msgId int64) (*model.User, utils.MessageErr)
	createUserService func(message *model.User) (*model.User, utils.MessageErr)
	updateUserService func(message *model.User) (*model.User, utils.MessageErr)
	deleteUserService func(msgId int64) utils.MessageErr
	getAllUserService func() ([]model.User, utils.MessageErr)
)

type serviceMock struct{}

func (sm *serviceMock) GetUser(msgId int64) (*model.User, utils.MessageErr) {
	return getUserService(msgId)
}

func (sm *serviceMock) SaveUser(message *model.User) (*model.User, utils.MessageErr) {
	return createUserService(message)
}
func (sm *serviceMock) UpdateUser(message *model.User) (*model.User, utils.MessageErr) {
	return updateUserService(message)
}
func (sm *serviceMock) DeleteUser(msgId int64) utils.MessageErr {
	return deleteUserService(msgId)
}
func (sm *serviceMock) GetAllUsers() ([]model.User, utils.MessageErr) {
	return getAllUserService()
}

// /////////////////////////////////////////////////////////////
// "GetUser" test cases
// /////////////////////////////////////////////////////////////
func TestUserController_GetUser_Success(t *testing.T) {
	// Given
	var service service.IUserService = &serviceMock{}
	var userController = UserController{service}

	getUserService = func(msgId int64) (*model.User, utils.MessageErr) {
		return &model.User{
			Id:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@gmail.com",
			Age:       30,
		}, nil
	}
	msgId := "1"         //this has to be a string, because is passed through the url
	r := chi.NewRouter() //chi.NewMux()

	// When
	req, _ := http.NewRequest(http.MethodGet, "/users/"+msgId, nil)
	rr := httptest.NewRecorder()
	r.Get("/users/{user_id}", userController.GetUser)
	r.ServeHTTP(rr, req)

	// Then
	var user model.User
	err := json.Unmarshal(rr.Body.Bytes(), &user)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "John", user.FirstName)
	assert.EqualValues(t, "Doe", user.LastName)
}

func TestGetMessage_Invalid_Id(t *testing.T) {
	// Given
	var service service.IUserService = &serviceMock{}
	var userController = UserController{service}

	getUserService = func(msgId int64) (*model.User, utils.MessageErr) {
		return &model.User{
			Id:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@gmail.com",
			Age:       30,
		}, nil
	}
	msgId := "abc" //this has to be a string, because is passed through the url
	r := chi.NewRouter()
	req, _ := http.NewRequest(http.MethodGet, "/users/"+msgId, nil)
	rr := httptest.NewRecorder()
	r.Get("/users/{user_id}", userController.GetUser)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.ApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "user id should be a number", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

func TestGetUser_User_Not_Found(t *testing.T) {
	// Given
	var service service.IUserService = &serviceMock{}
	var userController = UserController{service}
	getUserService = func(msgId int64) (*model.User, utils.MessageErr) {
		return nil, utils.NotFoundError("message not found")
	}
	msgId := "1" //valid id
	r := chi.NewRouter()
	req, _ := http.NewRequest(http.MethodGet, "/users/"+msgId, nil)
	rr := httptest.NewRecorder()
	r.Get("/users/{user_id}", userController.GetUser)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.ApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusNotFound, apiErr.Status())
	assert.EqualValues(t, "message not found", apiErr.Message())
	assert.EqualValues(t, "not_found", apiErr.Error())
}

func TestGetUser_User_Database_Error(t *testing.T) {
	// Given
	var userService service.IUserService = &serviceMock{}
	var userController = UserController{userService}
	getUserService = func(msgId int64) (*model.User, utils.MessageErr) {
		return nil, utils.InternalServerError("database error")
	}
	msgId := "1" //valid id
	r := chi.NewRouter()
	req, _ := http.NewRequest(http.MethodGet, "/users/"+msgId, nil)
	rr := httptest.NewRecorder()
	r.Get("/users/{user_id}", userController.GetUser)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.ApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
	assert.EqualValues(t, "database error", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
}

// /////////////////////////////////////////////////////////////
// "CreateUser" test cases
// /////////////////////////////////////////////////////////////
func TestCreateUser_Success(t *testing.T) {
	// Given
	var userService service.IUserService = &serviceMock{}
	var userController = UserController{userService}
	createUserService = func(message *model.User) (*model.User, utils.MessageErr) {
		return &model.User{
			Id:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@gmail.com",
			Age:       30,
		}, nil
	}
	jsonBody := `{"Id": 1, "first_name": "John", "last_name": "Doe", "email": "john.doe@gmail.com", "age": 30}`
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(jsonBody))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Post("/users", userController.SaveUser)
	r.ServeHTTP(rr, req)

	var message model.User
	err = json.Unmarshal(rr.Body.Bytes(), &message)
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.EqualValues(t, http.StatusCreated, rr.Code)
	assert.EqualValues(t, 1, message.Id)
	assert.EqualValues(t, "John", message.FirstName)
	assert.EqualValues(t, "Doe", message.LastName)
	assert.EqualValues(t, "john.doe@gmail.com", message.Email)
	assert.EqualValues(t, 30, message.Age)
}

func TestCreateUser_Invalid_Json(t *testing.T) {
	// Given
	var userService service.IUserService = &serviceMock{}
	var userController = UserController{userService}
	inputJson := `{"Id": 1, "first_name": 456, "last_name": "Doe", "email": "john.doe@gmail.com", "age": 30}`
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Post("/users", userController.SaveUser)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.ApiErrFromBytes(rr.Body.Bytes())

	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "json: cannot unmarshal number into Go struct field User.first_name of type string", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

func TestCreateUser_Empty_FirstName(t *testing.T) {
	// Given
	var userService service.IUserService = &serviceMock{}
	var userController = UserController{userService}
	createUserService = func(message *model.User) (*model.User, utils.MessageErr) {
		return nil, utils.UnprocessibleEntityError("Please enter a valid firstname")
	}
	inputJson := `{"Id": 1, "first_name": "", "last_name": "Doe", "email": "john.doe@gmail.com", "age": 30}`
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Post("/users", userController.SaveUser)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.ApiErrFromBytes(rr.Body.Bytes())

	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "Please enter a valid firstname", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

func TestCreateUser_Empty_Lastname(t *testing.T) {
	// Given
	var userService service.IUserService = &serviceMock{}
	var userController = UserController{userService}
	createUserService = func(message *model.User) (*model.User, utils.MessageErr) {
		return nil, utils.UnprocessibleEntityError("Please enter a valid lastname")
	}
	inputJson := `{"Id": 1, "first_name": "John", "last_name": "", "email": "john.doe@gmail.com", "age": 30}`
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Post("/users", userController.SaveUser)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.ApiErrFromBytes(rr.Body.Bytes())

	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "Please enter a valid lastname", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

// /////////////////////////////////////////////////////////////
// "UpdateUser" test cases
// /////////////////////////////////////////////////////////////
func TestUpdateUser_Success(t *testing.T) {
	// Given
	var userService service.IUserService = &serviceMock{}
	var userController = UserController{userService}
	updateUserService = func(message *model.User) (*model.User, utils.MessageErr) {
		return &model.User{
			Id:        1,
			FirstName: "Johnny",
			LastName:  "Dover",
			Email:     "johnny.dover@gmail.com",
			Age:       37,
		}, nil
	}
	jsonBody := `{"Id": 1, "first_name": "Johnny", "last_name": "Dover", "email": "johnny.dover@gmail.com", "age": 37}`
	r := chi.NewRouter()
	id := "1"
	req, err := http.NewRequest(http.MethodPut, "/users/"+id, bytes.NewBufferString(jsonBody))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Put("/users/{user_id}", userController.UpdateUser)
	r.ServeHTTP(rr, req)

	var message model.User
	err = json.Unmarshal(rr.Body.Bytes(), &message)
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, 1, message.Id)
	assert.EqualValues(t, 1, message.Id)
	assert.EqualValues(t, "Johnny", message.FirstName)
	assert.EqualValues(t, "Dover", message.LastName)
	assert.EqualValues(t, "johnny.dover@gmail.com", message.Email)
	assert.EqualValues(t, 37, message.Age)
}

func TestUpdateUser_Invalid_Id(t *testing.T) {
	// Given
	var userService service.IUserService = &serviceMock{}
	var userController = UserController{userService}
	jsonBody := `{"Id": "abc", "first_name": "Johnny", "last_name": "Dover", "email": "johnny.dover@gmail.com", "age": 37}`
	r := chi.NewRouter()
	id := "abc"
	req, err := http.NewRequest(http.MethodPut, "/users/"+id, bytes.NewBufferString(jsonBody))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Put("/users/{user_id}", userController.UpdateUser)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.ApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "json: cannot unmarshal string into Go struct field User.id of type int64", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

func TestUpdateUser_Invalid_Json(t *testing.T) {
	// Given
	var userService service.IUserService = &serviceMock{}
	var userController = UserController{userService}
	inputJson := `{"Id": 1, "first_name": 456, "last_name": "Doe", "email": "john.doe@gmail.com", "age": 30}`
	r := chi.NewRouter()
	id := "1"
	req, err := http.NewRequest(http.MethodPut, "/users/"+id, bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Put("/users/{user_id}", userController.UpdateUser)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.ApiErrFromBytes(rr.Body.Bytes())

	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "json: cannot unmarshal number into Go struct field User.first_name of type string", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

func TestUpdateUser_Empty_Firstname(t *testing.T) {
	// Given
	var userService service.IUserService = &serviceMock{}
	var userController = UserController{userService}
	updateUserService = func(message *model.User) (*model.User, utils.MessageErr) {
		return nil, utils.BadRequestError("Please enter a valid firstname")
	}
	inputJson := `{"Id": 1, "first_name": "", "last_name": "Doe", "email": "john.doe@gmail.com", "age": 30}`
	id := "1"
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodPut, "/users/"+id, bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Put("/users/{user_id}", userController.UpdateUser)
	r.ServeHTTP(rr, req)
	apiErr, err := utils.ApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "Please enter a valid firstname", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

func TestUpdateUser_Empty_Lastname(t *testing.T) {
	// Given
	var userService service.IUserService = &serviceMock{}
	var userController = UserController{userService}
	updateUserService = func(message *model.User) (*model.User, utils.MessageErr) {
		return nil, utils.BadRequestError("Please enter a valid lastname")
	}
	inputJson := `{"Id": 1, "first_name": "John", "last_name": "", "email": "john.doe@gmail.com", "age": 30}`
	id := "1"
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodPut, "/users/"+id, bytes.NewBufferString(inputJson))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Put("/users/{user_id}", userController.UpdateUser)
	r.ServeHTTP(rr, req)
	apiErr, err := utils.ApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "Please enter a valid lastname", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

func TestUpdateUser_Error_Updating(t *testing.T) {
	// Given
	var userService service.IUserService = &serviceMock{}
	var userController = UserController{userService}
	updateUserService = func(message *model.User) (*model.User, utils.MessageErr) {
		return nil, utils.InternalServerError("error when updating user")
	}
	jsonBody := `{"Id": 1, "first_name": "Johnny", "last_name": "Dover", "email": "johnny.dover@gmail.com", "age": 37}`
	r := chi.NewRouter()
	id := "1"
	req, err := http.NewRequest(http.MethodPut, "/users/"+id, bytes.NewBufferString(jsonBody))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Put("/users/{user_id}", userController.UpdateUser)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.ApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)

	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
	assert.EqualValues(t, "error when updating user", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
}

// /////////////////////////////////////////////////////////////
// "DeleteUser" test cases
// /////////////////////////////////////////////////////////////
func TestDeleteUser_Success(t *testing.T) {
	// Given
	var userService service.IUserService = &serviceMock{}
	var userController = UserController{userService}
	deleteUserService = func(msg int64) utils.MessageErr {
		return nil
	}
	r := chi.NewRouter()
	id := "1"
	req, err := http.NewRequest(http.MethodDelete, "/users/"+id, nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Delete("/users/{user_id}", userController.DeleteUser)
	r.ServeHTTP(rr, req)

	var response = make(map[string]string)
	theErr := json.Unmarshal(rr.Body.Bytes(), &response)
	if theErr != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, response["status"], "deleted")
}

func TestDeleteUser_Invalid_Id(t *testing.T) {
	// Given
	var userService service.IUserService = &serviceMock{}
	var userController = UserController{userService}

	r := chi.NewRouter()
	id := "abc"
	req, err := http.NewRequest(http.MethodDelete, "/users/"+id, nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Delete("/users/{user_id}", userController.DeleteUser)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.ApiErrFromBytes(rr.Body.Bytes())
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "user id should be a number", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

func TestDeleteUser_Failure(t *testing.T) {
	// Given
	var userService service.IUserService = &serviceMock{}
	var userController = UserController{userService}
	deleteUserService = func(msg int64) utils.MessageErr {
		return utils.InternalServerError("error deleting message")
	}
	r := chi.NewRouter()
	id := "1"
	req, err := http.NewRequest(http.MethodDelete, "/users/"+id, nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Delete("/users/{user_id}", userController.DeleteUser)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.ApiErrFromBytes(rr.Body.Bytes())
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
	assert.EqualValues(t, "error deleting message", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
}

// /////////////////////////////////////////////////////////////
// "GetAllUsers" test cases
// /////////////////////////////////////////////////////////////
func TestGetAllUsers_Success(t *testing.T) {
	// Given
	var userService service.IUserService = &serviceMock{}
	var userController = UserController{userService}
	getAllUserService = func() ([]model.User, utils.MessageErr) {
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
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Get("/users", userController.ListUsers)
	r.ServeHTTP(rr, req)

	var users []model.User
	theErr := json.Unmarshal(rr.Body.Bytes(), &users)
	if theErr != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	assert.Nil(t, err)
	assert.NotNil(t, users)
	assert.EqualValues(t, users[0].Id, 1)
	assert.EqualValues(t, users[0].FirstName, "John")
	assert.EqualValues(t, users[0].LastName, "Doe")
	assert.EqualValues(t, users[1].Id, 2)
	assert.EqualValues(t, users[1].FirstName, "Johnny")
	assert.EqualValues(t, users[1].LastName, "Dover")
}

func TestGetAllUsers_Failure(t *testing.T) {
	// Given
	var userService service.IUserService = &serviceMock{}
	var userController = UserController{userService}
	getAllUserService = func() ([]model.User, utils.MessageErr) {
		return nil, utils.InternalServerError("error getting messages")
	}
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.Get("/users", userController.ListUsers)
	r.ServeHTTP(rr, req)

	apiErr, err := utils.ApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, "error getting messages", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
}
