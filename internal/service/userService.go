package service

import (
	"github.com/wexinc/ps-tag-onboarding-go/internal/log"
	"github.com/wexinc/ps-tag-onboarding-go/internal/model"
	"github.com/wexinc/ps-tag-onboarding-go/internal/repository"
	"github.com/wexinc/ps-tag-onboarding-go/internal/utils"
	"strings"
)

type IUserService interface {
	GetAllUsers() ([]model.User, utils.MessageErr)
	GetUser(id int64) (*model.User, utils.MessageErr)
	SaveUser(user *model.User) (*model.User, utils.MessageErr)
	UpdateUser(user *model.User) (*model.User, utils.MessageErr)
	DeleteUser(id int64) utils.MessageErr
}

type UserService struct {
	Repository        repository.IUserRepository
	ValidationService IUserValidationService
}

func (us *UserService) GetAllUsers() ([]model.User, utils.MessageErr) {

	userList, err := us.Repository.DbListUsers()

	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	return userList, nil
}

func (us *UserService) GetUser(id int64) (*model.User, utils.MessageErr) {

	user, err := us.Repository.DbGetUser(id)

	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	return user, nil
}

func (us *UserService) SaveUser(user *model.User) (*model.User, utils.MessageErr) {

	// validate user
	validationErr := us.ValidationService.ValidateUser(user)

	if len(validationErr) > 0 {
		log.Error.Println(validationErr)
		return nil, utils.BadRequestError(strings.Join(validationErr, ","))
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

func (us *UserService) UpdateUser(user *model.User) (*model.User, utils.MessageErr) {

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

func (us *UserService) DeleteUser(id int64) utils.MessageErr {
	//verify if user exist
	user, err := us.Repository.DbGetUser(id)
	if err != nil {
		return err
	}
	if user == nil {
		log.Error.Printf("User with id '%v' cannot be found", id)
		return utils.NotFoundError("user not found")
	}
	err = us.Repository.DbDeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}
