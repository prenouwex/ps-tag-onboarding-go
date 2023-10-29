package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/wexinc/ps-tag-onboarding-go/database"
	"github.com/wexinc/ps-tag-onboarding-go/internal/model"
	"github.com/wexinc/ps-tag-onboarding-go/internal/repository"
	"github.com/wexinc/ps-tag-onboarding-go/internal/service"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BuildRouter() *chi.Mux {
	db := database.CreateNewGormDB()
	userRepository := repository.UserRepository{DB: db}
	userValidation := service.UserValidationService{&userRepository}
	userService := service.UserService{&userRepository, &userValidation}
	r := chi.NewRouter()
	userRoutes := UserRoutes{&userService}
	userRoutes.UserRoutes(r)
	return r
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		rec            *httptest.ResponseRecorder
		req            *http.Request
		reqPath        string
		body           io.Reader
		expectedBody   string
		expectedHeader string
	}{
		{
			name:         "OK_1",
			method:       http.MethodGet,
			rec:          httptest.NewRecorder(),
			reqPath:      "/users/1",
			expectedBody: `{"id":1,"first_name":"John","last_name":"Doe","email":"john.doe@yahoo.com","age":34}`,
		},
		{
			name:         "OK_2",
			method:       http.MethodGet,
			rec:          httptest.NewRecorder(),
			reqPath:      "/users/2",
			expectedBody: `{"id":2,"first_name":"Zenia","last_name":"Brennan","email":"ultrices.vivamus.rhoncus@yahoo.ca","age":34}`,
		},
		{
			name:         "BAD_REQUEST",
			method:       http.MethodGet,
			rec:          httptest.NewRecorder(),
			reqPath:      "/users/bad",
			expectedBody: `{"code":400,"message":"User not found"}null`,
		},
	}

	testServer := httptest.NewServer(BuildRouter())
	defer testServer.Close()

	verify(t, tests, testServer)
}

func TestGetAllUser(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		rec            *httptest.ResponseRecorder
		req            *http.Request
		reqPath        string
		body           io.Reader
		expectedBody   string
		expectedHeader string
	}{
		{
			name:         "OK",
			method:       http.MethodGet,
			rec:          httptest.NewRecorder(),
			reqPath:      "/users",
			expectedBody: `[{"id":1,"first_name":"John","last_name":"Doe","email":"john.doe@yahoo.com","age":34},{"id":2,"first_name":"Zenia","last_name":"Brennan","email":"ultrices.vivamus.rhoncus@yahoo.ca","age":34},{"id":3,"first_name":"Branden","last_name":"Spears","email":"non.lobortis@hotmail.net","age":34},{"id":4,"first_name":"Alice","last_name":"Wallace","email":"at@protonmail.couk","age":34},{"id":5,"first_name":"Ira","last_name":"Francis","email":"in.lobortis.tellus@protonmail.ca","age":34}]`,
		},
	}

	testServer := httptest.NewServer(BuildRouter())
	defer testServer.Close()

	verify(t, tests, testServer)
}

func TestSaveUser(t *testing.T) {

	user := &model.User{
		FirstName: "Nic",
		LastName:  "Raboy",
		Email:     "nic.raboy_t@gmail.com",
		Age:       45,
	}

	existingUser := &model.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "nic.raboy_t@gmail.com",
		Age:       45,
	}

	jsonUser, _ := json.Marshal(user)
	jsonExistingUser, _ := json.Marshal(existingUser)

	tests := []struct {
		name           string
		method         string
		rec            *httptest.ResponseRecorder
		req            *http.Request
		reqPath        string
		body           io.Reader
		expectedBody   string
		expectedHeader string
	}{
		{
			name:         "SAVE_OK",
			method:       http.MethodPost,
			rec:          httptest.NewRecorder(),
			reqPath:      "/users",
			body:         bytes.NewBuffer(jsonUser),
			expectedBody: `"{'id':6}"`,
		},
		{
			name:         "SAVE_FAILED",
			method:       http.MethodPost,
			rec:          httptest.NewRecorder(),
			reqPath:      "/users",
			body:         bytes.NewBuffer(jsonExistingUser),
			expectedBody: `{"code":400,"message":"User with the same first and last name already exists"}`,
		},
	}

	testServer := httptest.NewServer(BuildRouter())
	defer testServer.Close()

	verify(t, tests, testServer)
}

func TestUpdateUser(t *testing.T) {

	user := &model.User{
		Id:        1,
		FirstName: "Nic",
		LastName:  "Raboy",
		Email:     "nic.raboy_t@gmail.com",
		Age:       45,
	}

	jsonUser, _ := json.Marshal(user)

	tests := []struct {
		name           string
		method         string
		rec            *httptest.ResponseRecorder
		req            *http.Request
		reqPath        string
		body           io.Reader
		expectedBody   string
		expectedHeader string
	}{
		{
			name:    "UPDATE_OK",
			method:  http.MethodPut,
			rec:     httptest.NewRecorder(),
			reqPath: fmt.Sprint("/users/", user.Id),
			body:    bytes.NewBuffer(jsonUser), //bytes.NewBuffer([]byte("{'id':6,'first_name':'Ben','last_name':'Jefferson','email':'t.jefferson@yahoo.com','age':39}")), //json.Marshal(user.User{}))
			//req:          httptest.NewRequest("GET", "/users", nil),
			expectedBody: `"{'id':1}"`,
		},
	}

	testServer := httptest.NewServer(BuildRouter())
	defer testServer.Close()

	verify(t, tests, testServer)
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		rec            *httptest.ResponseRecorder
		req            *http.Request
		reqPath        string
		body           io.Reader
		expectedBody   string
		expectedHeader string
	}{
		{
			name:    "DELETE_OK_1",
			method:  http.MethodDelete,
			rec:     httptest.NewRecorder(),
			reqPath: "/users/1",
			//body:		  bytes.NewBuffer([]byte("")),//json.Marshal(user.User{}))
			//req:          httptest.NewRequest("GET", "/users/1", nil),
			expectedBody: `""`,
		},
	}

	testServer := httptest.NewServer(BuildRouter())
	defer testServer.Close()

	verify(t, tests, testServer)
}

func verify(t *testing.T, tests []struct {
	name           string
	method         string
	rec            *httptest.ResponseRecorder
	req            *http.Request
	reqPath        string
	body           io.Reader
	expectedBody   string
	expectedHeader string
}, testServer *httptest.Server) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest(test.method, testServer.URL+test.reqPath, test.body)
			if err != nil {
				t.Fatal(err)
			}

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatal(err)
			}

			respBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			assert.Equal(t, test.expectedBody, string(respBody))
		})
	}
}
