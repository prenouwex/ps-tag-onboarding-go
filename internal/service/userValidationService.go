package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/wexinc/ps-tag-onboarding-go/internal/log"
	"github.com/wexinc/ps-tag-onboarding-go/internal/model"
	"github.com/wexinc/ps-tag-onboarding-go/internal/repository"
	"strings"
)

const (
	ERROR_AGE_MINIMUM          = "User does not meet minimum age requirement"
	ERROR_EMAIL_FORMAT         = "User email must be properly formatted"
	ERROR_NAME_UNIQUE          = "User with the same first and last name already exists"
	RESPONSE_USER_NOT_FOUND    = "User not found"
	RESPONSE_VALIDATION_FAILED = "User did not pass validation"
)

type IUserValidationService interface {
	ValidateUser(user *model.User) ([]string, error)
}

type UserValidationService struct {
	Repository repository.IUserRepository //*repository.UserRepository
}

func (uvs *UserValidationService) ValidateUser(user *model.User) ([]string, error) {

	validate := validator.New()
	validationErr := []string{}

	// Register custom validation function with the validator
	validate.RegisterValidation("validateAge", uvs.validateAge)
	validate.RegisterValidation("validateEmail", uvs.validateEmail)

	// Test custom field validation
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {

			switch err.Tag() {

			case "validateAge":
				log.Error.Print("validateAge error:", err.Error())
				validationErr = append(validationErr, ERROR_AGE_MINIMUM)
			case "validateEmail", "email":
				log.Error.Print("validateEmail error:", err.Error())
				validationErr = append(validationErr, ERROR_EMAIL_FORMAT)
			default:
				log.Error.Print("validation error:", err.Error())
				validationErr = append(validationErr, err.Error())

			}
		}
	}

	// validate firstName and lastName
	if !uvs.validateFirstNameLastName(user.FirstName, user.LastName) {
		log.Error.Print("validateFirstNameLastName error")
		validationErr = append(validationErr, ERROR_NAME_UNIQUE)
	}

	if len(validationErr) > 0 {
		log.Info.Println(RESPONSE_VALIDATION_FAILED)
		return validationErr, nil
	} else {
		log.Info.Println("Validation successful")
	}

	return nil, nil
}

func (uvs *UserValidationService) validateAge(fl validator.FieldLevel) bool {
	age := fl.Field().Int()
	log.Info.Println("in validateAge")
	return age >= 10
}

func (uvs *UserValidationService) validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	return len(strings.TrimSpace(email)) != 0 && strings.Contains(email, "@")
}

func (uvs *UserValidationService) validateFirstNameLastName(firstName string, lastName string) bool {
	exist, err := uvs.Repository.ExistsByFirstNameAndLastName(firstName, lastName)
	if err != nil {
		panic(err)
	}
	return !exist
}
