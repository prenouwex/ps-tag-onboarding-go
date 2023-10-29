package model

type User struct {
	Id        int64  `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email,validateEmail"`
	Age       int64  `json:"age" validate:"required,validateAge"`
}
