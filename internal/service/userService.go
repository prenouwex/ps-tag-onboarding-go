package service

import (
	"errors"
	"github.com/wexinc/ps-tag-onboarding-go/internal/model"
	"github.com/wexinc/ps-tag-onboarding-go/internal/repository"
	"github.com/wexinc/ps-tag-onboarding-go/log"
	"strings"
)

type IUserService interface {
	GetAllUsers() ([]model.User, error)
	GetUser(id int64) (*model.User, error)
	SaveUser(user *model.User) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	DeleteUser(id int64) (*model.User, error)
}

type UserService struct {
	Repository        repository.IUserRepository
	ValidationService IUserValidationService
}

func (us *UserService) GetAllUsers() ([]model.User, error) {

	userList, err := us.Repository.DbListUsers()

	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	return userList, nil
}

func (us *UserService) GetUser(id int64) (*model.User, error) {

	user, err := us.Repository.DbGetUser(id)

	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	return user, nil
}

func (us *UserService) SaveUser(user *model.User) (*model.User, error) {

	// TODO : update validateUser to only return one list of errors
	// validate user
	valiationErr, err := us.ValidationService.ValidateUser(user)

	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	if len(valiationErr) > 0 {
		log.Error.Println(valiationErr)
		return nil, errors.New(strings.Join(valiationErr, ","))
	}

	// create user
	userCreated, err := us.Repository.DbCreateUser(user)

	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	log.Info.Printf("User created : %v", userCreated)

	return userCreated, nil

}

func (us *UserService) UpdateUser(user *model.User) (*model.User, error) {

	log.Info.Printf("User update service ")

	log.Info.Printf("User to update: %v", user)

	// load up existing user with same id
	current, err := us.Repository.DbGetUser(user.Id)
	if err != nil {
		return nil, err
	}
	current.FirstName = user.FirstName
	current.LastName = user.LastName
	current.Email = user.Email
	current.Age = user.Age

	// update user
	updateUser, err := us.Repository.DbUpdateUser(current)
	if err != nil {
		return nil, err
	}
	log.Info.Printf("User updated with details %v", updateUser)
	return updateUser, nil
}

func (us *UserService) DeleteUser(id int64) (*model.User, error) {
	//verify if user exist
	user, err := us.Repository.DbGetUser(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		log.Error.Printf("User with id '%v' cannot be found", id)
		return nil, err
	}
	deleteErr, err := us.Repository.DbDeleteUser(id)
	if err != nil {
		return nil, err
	}
	return deleteErr, nil
}
