package database

import (
	"github.com/wexinc/ps-tag-onboarding-go/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateNewGormDB() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("file:userdb?mode=memory&cache=shared"), &gorm.Config{}, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}

	// Insert some users
	var users = []*model.User{
		{Id: 1, FirstName: "John", LastName: "Doe", Email: "john.doe@yahoo.com", Age: 34},
		{Id: 2, FirstName: "Zenia", LastName: "Brennan", Email: "ultrices.vivamus.rhoncus@yahoo.ca", Age: 34},
		{Id: 3, FirstName: "Branden", LastName: "Spears", Email: "non.lobortis@hotmail.net", Age: 34},
		{Id: 4, FirstName: "Alice", LastName: "Wallace", Email: "at@protonmail.couk", Age: 34},
		{Id: 5, FirstName: "Ira", LastName: "Francis", Email: "in.lobortis.tellus@protonmail.ca", Age: 34},
	}
	for _, u := range users {
		if err := db.Save(u).Error; err != nil {
			panic(err)
		}
	}

	return db

}
