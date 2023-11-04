package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/wexinc/ps-tag-onboarding-go/internal/log"
	"github.com/wexinc/ps-tag-onboarding-go/internal/model"
	"github.com/wexinc/ps-tag-onboarding-go/internal/service"
	"github.com/wexinc/ps-tag-onboarding-go/internal/utils"
	"net/http"
	"strconv"
)

type IUserController interface {
	ListUsers(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	SaveUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type UserController struct {
	UserService service.IUserService
}

// GetUser Get a list of all users
//
// This will returns a list of users.
//
// swagger:route GET /users/ getAllUser
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Responses:
//
//	201: User
//	400: MessageErr
//	500: MessageErr
func (uc *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {

	userList, err := uc.UserService.GetAllUsers()

	if err != nil {
		log.Error.Println(err)
		utils.ResponseCustomError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.ResponseJson(w, http.StatusOK, userList)
}

// GetUser Get a user
//
// This will handle GET requests for retrieving a user by ID.
//
// swagger:route GET /users/{userId} getUser
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Responses:
//
//	200: User
//	400: MessageErr
//	404: MessageErr
//	500: MessageErr
func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {

	userId := chi.URLParam(r, "userId")
	id, _ := strconv.Atoi(userId)

	user, err := uc.UserService.GetUser(int64(id))

	if err != nil {
		log.Error.Println(err)
		utils.ResponseCustomError(w, err.Status(), err.Error())
	}

	utils.ResponseJson(w, http.StatusOK, user)
}

// SaveUser creates a new user
//
// This will create a new user based on the information provided in the request body.
//
// swagger:route POST /users saveUser
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Responses:
//
//	201: User
//	400: MessageErr
//	500: MessageErr
//
// responses.createUserCreated.headers.body.type: UserResponse
func (uc *UserController) SaveUser(w http.ResponseWriter, r *http.Request) {

	var body model.User
	if err := utils.ParseJson(r, &body); err != nil {
		log.Error.Println(err)
		utils.ResponseError(w, http.StatusBadRequest)
		return
	}

	// create user
	user, err := uc.UserService.SaveUser(&body)

	if err != nil {
		log.Error.Println(err)
		utils.ResponseCustomError(w, err.Status(), err.Message())
		return
	}

	log.Info.Printf("User created : %v", user)

	utils.ResponseJson(w, http.StatusCreated, user)

}

// UpdateUser Updates a user by ID
//
// This will Update a user.
//
// swagger:route PUT /users/{userId} updateUser
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Responses:
//
//	200: User
//	400: MessageErr
//	500: MessageErr
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {

	log.Info.Printf("User update service ")

	var body model.User
	if err := utils.ParseJson(r, &body); err != nil {
		log.Error.Println(err)
		utils.ResponseCustomError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Info.Printf("User to update: %v", body)

	user, err := uc.UserService.UpdateUser(&body)

	if err != nil {
		log.Error.Println(err)
		utils.ResponseCustomError(w, err.Status(), err.Message())
		return
	}

	log.Info.Printf("User updated : %v", user)

	utils.ResponseJson(w, http.StatusOK, user)

}

// DeleteUser Deletes a user by ID
//
// This will delete a user.
//
// swagger:route DELETE /users/{userId} deleteUser
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Responses:
//
//	200: User
//	400: MessageErr
//	404: MessageErr
//	500: MessageErr
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {

	userId := chi.URLParam(r, "userId")
	id, _ := strconv.Atoi(userId)

	err := uc.UserService.DeleteUser(int64(id))

	if err != nil {
		log.Error.Println(err)
		utils.ResponseCustomError(w, err.Status(), err.Message())
	}

	log.Info.Print("User deleted with id : ", userId)

	utils.ResponseJson(w, http.StatusOK, map[string]string{"status": "deleted"})

}

// swagger:parameters getUser deleteUser updateUser
type UserPathParam struct {
	// in: path
	UserId string `json:"userId"`
}

// swagger:parameters updateUser saveUser
type UserBodyParam struct {
	// in:body
	User model.User
}
