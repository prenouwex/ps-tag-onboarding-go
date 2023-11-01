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

// ListUsers returns a list of users.
//
// @Summary List users
// @Description Retrieve a list of users
// @ID list_users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} User "List of users"
// @Failure 400 MessageErr "Invalid request"
// @Failure 500 MessageErr "Server error"
// @Router /users [get]
func (uc *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {

	userList, err := uc.UserService.GetAllUsers()

	if err != nil {
		log.Error.Println(err)
		utils.ResponseCustomError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.ResponseJson(w, http.StatusOK, userList)
}

// @Summary Get a user
// @Description GetUser handles GET requests for retrieving a user by ID.
// @ID update_user
// @Tags users
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param user body User true "User data"
// @Success 200 {object} User "User successfully updated"
// @Failure 400 {object} MessageErr "Invalid ID or payload"
// @Failure 404 {object} MessageErr "User not found"
// @Failure 500 {object} MessageErr "Server error"
// @Router /users/{userId} [put]
func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {

	userId := chi.URLParam(r, "userId")
	id, _ := strconv.Atoi(userId)

	user, err := uc.UserService.GetUser(int64(id))

	if err != nil {
		log.Error.Println(err)
		utils.ResponseCustomError(w, http.StatusBadRequest, err.Error())
	}

	utils.ResponseJson(w, http.StatusOK, user)
}

// SaveUser creates a new user
//
// This will create a new user based on the information provided in the request body.
//
// swagger:route POST /users createUser
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

// UpdateUser updates a user by ID.
//
// @Summary Update a user
// @Description Update a user by ID
// @ID update_user
// @Tags users
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param user body User true "User data"
// @Success 200 {object} User "User successfully updated"
// @Failure 400 {object} MessageErr "Invalid ID or payload"
// @Failure 404 {object} MessageErr "User not found"
// @Failure 500 {object} MessageErr "Server error"
// @Router /users/{userId} [put]
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

// DeleteUser deletes a user by ID.
//
// @Summary Delete a user
// @Description Delete a user by ID
// @ID delete_user
// @Tags users
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {string} string "User successfully deleted"
// @Failure 400 {object} MessageErr "Invalid ID"
// @Failure 404 {object} MessageErr "User not found"
// @Failure 500 {object} MessageErr "Server error"
// @Router /users/{userId} [delete]
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {

	userId := chi.URLParam(r, "userId")
	id, _ := strconv.Atoi(userId)

	_, err := uc.UserService.DeleteUser(int64(id))

	if err != nil {
		log.Error.Println(err)
		utils.ResponseCustomError(w, err.Status(), err.Message())
	}

	log.Info.Print("User deleted with id : ", userId)

	utils.ResponseJson(w, http.StatusOK, map[string]string{"status": "deleted"})

}
