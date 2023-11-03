package repository

import (
	"fmt"
	"github.com/wexinc/ps-tag-onboarding-go/internal/model"
	"github.com/wexinc/ps-tag-onboarding-go/internal/utils"
	"gorm.io/gorm"
)

const (
	USER_NOT_FOUND = "user not found with id %v"
	NO_USER_FOUND  = "no users found"
)

type IUserRepository interface {
	DbListUsers() ([]model.User, utils.MessageErr)
	DbCreateUser(user *model.User) (*model.User, utils.MessageErr)
	DbGetUser(id int64) (*model.User, utils.MessageErr)
	DbUpdateUser(user *model.User) (*model.User, utils.MessageErr)
	DbDeleteUser(id int64) utils.MessageErr
	ExistsByFirstNameAndLastName(firstName string, lastName string) (bool, utils.MessageErr)
	// Add other necessary GORM methods here
}

type UserRepository struct {
	DB *gorm.DB
}

func (ur *UserRepository) DbListUsers() ([]model.User, utils.MessageErr) {

	var users []model.User

	if err := ur.DB.Find(&users).Error; err != nil {
		return nil, utils.InternalServerError(err.Error())
	}

	if len(users) > 0 {
		return users, nil
	}

	return nil, utils.NotFoundError(NO_USER_FOUND)

}

func (ur *UserRepository) DbCreateUser(user *model.User) (*model.User, utils.MessageErr) {

	if err := ur.DB.Save(user).Error; err != nil {
		return nil, utils.InternalServerError(err.Error())
	}

	return user, nil
}

func (ur *UserRepository) DbGetUser(id int64) (*model.User, utils.MessageErr) {

	var user model.User

	if err := ur.DB.Where(model.User{Id: id}).Take(&user).Error; err != nil {
		return nil, utils.InternalServerError(err.Error())
	}

	//if err := ur.DB.Where("id = ?", id).First(&user).Error; err != nil {
	//	return nil, utils.InternalServerError(err.Error())
	//}

	if user.Id == id {
		return &user, nil
	}

	return nil, utils.NotFoundError(fmt.Sprintf(USER_NOT_FOUND, id))
}

func (ur *UserRepository) DbUpdateUser(user *model.User) (*model.User, utils.MessageErr) {

	if err := ur.DB.Model(&model.User{}).Where("Id = ?", user.Id).Updates(user).Error; err != nil {
		return nil, utils.InternalServerError(err.Error())
	}

	return user, nil
}

func (ur *UserRepository) DbDeleteUser(id int64) utils.MessageErr {

	ur.DB.Delete(&model.User{Id: id})

	return nil
}

func (ur *UserRepository) ExistsByFirstNameAndLastName(firstName string, lastName string) (bool, utils.MessageErr) {

	// Query to find users with the specified first and last names
	var users []model.User
	if err := ur.DB.Where("first_name = ? AND last_name = ?", firstName, lastName).Find(&users).Error; err != nil {
		return false, utils.InternalServerError(err.Error())
	}

	if len(users) > 0 {
		return true, nil
	}

	return false, nil

}
