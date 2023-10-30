package repository

import (
	//"errors"
	//"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/wexinc/ps-tag-onboarding-go/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestUserRepo_Get(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	//gdb, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	gdb, err := gorm.Open(sqlite.Open("file:userdb?mode=memory&cache=shared"), &gorm.Config{}, &gorm.Config{})

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := UserRepository{DB: gdb}

	// AutoMigrate the User model
	gdb.AutoMigrate(&model.User{})

	tests := []struct {
		name    string
		s       UserRepository
		msgId   int64
		mock    func()
		want    *model.User
		wantErr bool
	}{
		//{
		//	//When everything works as expected
		//	name:  "OK",
		//	s:     s,
		//	msgId: 1,
		//	mock: func() {
		//		//We added one row
		//		rows := sqlmock.NewRows([]string{"Id", "FirstName", "LastName", "Email", "Age"}).AddRow(1, "John", "Doe", "john.doe@gmail.com", 30)
		//		mock.ExpectPrepare("SELECT (.+) FROM users").ExpectQuery().WithArgs(1).WillReturnRows(rows)
		//	},
		//	want: &model.User{
		//		Id:        1,
		//		FirstName: "John",
		//		LastName:  "Doe",
		//		Email:     "john.doe@gmail.com",
		//		Age:       30,
		//	},
		//},
		{
			//When the role tried to access is not found
			name:  "Not Found",
			s:     s,
			msgId: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"Id", "FirstName", "LastName", "Email", "Age"}) //observe that we didnt add any role here
				mock.ExpectPrepare("SELECT (.+) FROM messages").ExpectQuery().WithArgs(1).WillReturnRows(rows)
			},
			wantErr: true,
		},
		//{
		//	//When invalid statement is provided, ie the SQL syntax is wrong(in this case, we provided a wrong database)
		//	name:  "Invalid Prepare",
		//	s:     s,
		//	msgId: 1,
		//	mock: func() {
		//		rows := sqlmock.NewRows([]string{"Id", "Title", "Body", "CreatedAt"}).AddRow(1, "title", "body", created_at)
		//		mock.ExpectPrepare("SELECT (.+) FROM wrong_table").ExpectQuery().WithArgs(1).WillReturnRows(rows)
		//	},
		//	wantErr: true,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.DbGetUser(tt.msgId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error new = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
