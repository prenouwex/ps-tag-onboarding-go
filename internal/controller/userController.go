package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/wexinc/ps-tag-onboarding-go/internal/model"
	"github.com/wexinc/ps-tag-onboarding-go/internal/service"
	"github.com/wexinc/ps-tag-onboarding-go/internal/utils"
	"github.com/wexinc/ps-tag-onboarding-go/log"
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

func (uc *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {

	userList, err := uc.UserService.GetAllUsers()

	if err != nil {
		log.Error.Println(err)
		utils.ResponseCustomError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.ResponseJson(w, http.StatusOK, userList)
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {

	userId := chi.URLParam(r, "userId")
	id, err := strconv.Atoi(userId)

	user, err := uc.UserService.GetUser(int64(id))

	if err != nil {
		log.Error.Println(err)
		utils.ResponseCustomError(w, http.StatusBadRequest, err.Error())
	}

	utils.ResponseJson(w, http.StatusOK, user)
}

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
		utils.ResponseCustomError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Info.Printf("User created : %v", user)

	utils.ResponseJson(w, http.StatusCreated, user)

}

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
		utils.ResponseError(w, http.StatusBadRequest)
		return
	}

	log.Info.Printf("User updated : %v", user)

	utils.ResponseJson(w, http.StatusOK, user)

}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {

	userId := chi.URLParam(r, "userId")
	id, err := strconv.Atoi(userId)

	_, err = uc.UserService.DeleteUser(int64(id))

	if err != nil {
		log.Error.Println(err)
		utils.ResponseCustomError(w, http.StatusBadRequest, err.Error())
	}

	log.Info.Print("User deleted with id : ", userId)

	utils.ResponseJson(w, http.StatusOK, map[string]string{"status": "deleted"})

}
