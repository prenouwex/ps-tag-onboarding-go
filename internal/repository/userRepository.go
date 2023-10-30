package repository

import (
	"errors"
	"github.com/wexinc/ps-tag-onboarding-go/internal/model"
	"gorm.io/gorm"
)

const (
	USER_NOT_FOUND = "User not found"
	NO_USER_FOUND  = "no users found"
)

type IUserRepository interface {
	DbListUsers() ([]model.User, error)
	DbCreateUser(user *model.User) (*model.User, error)
	DbGetUser(id int64) (*model.User, error)
	DbUpdateUser(user *model.User) (*model.User, error)
	DbDeleteUser(id int64) (*model.User, error)
	ExistsByFirstNameAndLastName(firstName string, lastName string) (bool, error)
	// Add other necessary GORM methods here
}

type UserRepository struct {
	DB *gorm.DB
}

func (ur *UserRepository) DbListUsers() ([]model.User, error) {

	var users []model.User

	if err := ur.DB.Find(&users).Error; err != nil {
		return nil, errors.New(err.Error())
	}

	if len(users) > 0 {
		return users, nil
	}

	return nil, errors.New(NO_USER_FOUND)

}

func (ur *UserRepository) DbCreateUser(user *model.User) (*model.User, error) {

	if err := ur.DB.Save(user).Error; err != nil {
		return nil, errors.New(err.Error())
	}

	return user, nil
}

func (ur *UserRepository) DbGetUser(id int64) (*model.User, error) {

	var user model.User

	if err := ur.DB.Where(model.User{Id: id}).Take(&user).Error; err != nil {
		return nil, errors.New(err.Error())
	}

	if user.Id == id {
		return &user, nil
	}

	return nil, errors.New(USER_NOT_FOUND)
}

func (ur *UserRepository) DbUpdateUser(user *model.User) (*model.User, error) {

	if err := ur.DB.Model(&model.User{}).Where("Id = ?", user.Id).Updates(user).Error; err != nil {
		return nil, errors.New(err.Error())
	}

	return user, nil
}

func (ur *UserRepository) DbDeleteUser(id int64) (*model.User, error) {

	if err := ur.DB.Delete(&model.User{Id: id}).Error; err != nil {
		return nil, errors.New(err.Error())
	} else {
		return nil, nil
	}
}

func (ur *UserRepository) ExistsByFirstNameAndLastName(firstName string, lastName string) (bool, error) {

	// Query to find users with the specified first and last names
	var users []model.User
	if err := ur.DB.Where("first_name = ? AND last_name = ?", firstName, lastName).Find(&users).Error; err != nil {
		return false, errors.New(err.Error())
	}

	if len(users) > 0 {
		return true, nil
	}

	return false, nil

}
