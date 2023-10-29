package service

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/wexinc/ps-tag-onboarding-go/internal/model"
	"github.com/wexinc/ps-tag-onboarding-go/internal/repository"
	"github.com/wexinc/ps-tag-onboarding-go/internal/utils"
	"github.com/wexinc/ps-tag-onboarding-go/log"
	"net/http"
	"strconv"
	"strings"
)

type IUserService interface {
	ListUsers(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	SaveUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type UserService struct {
	Repository        repository.IUserRepository
	ValidationService *UserValidationService
}

func (us *UserService) ListUsers(w http.ResponseWriter, r *http.Request) {

	userList, err := us.Repository.DbListUsers()

	if err != nil {
		log.Error.Println(err)
		utils.ResponseCustomError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.ResponseJson(w, http.StatusOK, userList)
}

func (us *UserService) GetUser(w http.ResponseWriter, r *http.Request) {

	userId := chi.URLParam(r, "userId")
	id, err := strconv.Atoi(userId)

	user, err := us.Repository.DbGetUser(int64(id))

	if err != nil {
		log.Error.Println(err)
		utils.ResponseCustomError(w, http.StatusBadRequest, err.Error())
	}

	utils.ResponseJson(w, http.StatusOK, user)
}

func (us *UserService) SaveUser(w http.ResponseWriter, r *http.Request) {

	var body model.User
	if err := utils.ParseJson(r, &body); err != nil {
		log.Error.Println(err)
		utils.ResponseError(w, http.StatusBadRequest)
		return
	}

	// validate user
	valiationErr, err := us.ValidationService.ValidateUser(&body)

	if err != nil {
		log.Error.Println(err)
		utils.ResponseCustomError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(valiationErr) > 0 {
		log.Error.Println(valiationErr)
		utils.ResponseCustomError(w, http.StatusBadRequest, strings.Join(valiationErr, ","))
		return
	}

	// create user
	id, err := us.Repository.DbCreateUser(&body)

	if err != nil {
		log.Error.Println(err)
		utils.ResponseError(w, http.StatusBadRequest)
		return
	}

	log.Info.Printf("User created with id %d", id)

	utils.ResponseJson(w, http.StatusCreated, fmt.Sprint("{'id':", id, "}"))

}

func (us *UserService) UpdateUser(w http.ResponseWriter, r *http.Request) {

	log.Info.Printf("User update service ")

	var body model.User
	if err := utils.ParseJson(r, &body); err != nil {
		log.Error.Println(err)
		utils.ResponseError(w, http.StatusBadRequest)
		return
	}

	log.Info.Printf("User to update: %v", body)

	//validateUser()

	u, err := us.Repository.DbUpdateUser(&body)

	if err != nil {
		log.Error.Println(err)
		utils.ResponseError(w, http.StatusBadRequest)
		return
	}

	log.Info.Printf("User updated with details %v", u)

	utils.ResponseJson(w, http.StatusOK, fmt.Sprint("{'id':", u.Id, "}"))

}

func (us *UserService) DeleteUser(w http.ResponseWriter, r *http.Request) {

	//validateUser()

	userId := chi.URLParam(r, "userId")
	id, err := strconv.Atoi(userId)

	_, err = us.Repository.DbDeleteUser(int64(id))

	if err != nil {
		log.Error.Println(err)
		utils.ResponseCustomError(w, http.StatusBadRequest, err.Error())
	}

	log.Info.Print("User deleted with id : ", userId)

	utils.ResponseJson(w, http.StatusOK, "")

}
