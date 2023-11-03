package repository

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/wexinc/ps-tag-onboarding-go/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reflect"
	"regexp"
	"testing"
)

func NewDbMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	sqlDB, mock, err := sqlmock.New()
	// using postgres driver as sqlite gorm one doesn't seem to work
	mockDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	return mockDB, mock, err
}

func TestUserRepo_GetId(t *testing.T) {

	mockDB, mock, err := NewDbMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	s := UserRepository{DB: mockDB}

	tests := []struct {
		name    string
		s       *UserRepository
		msgId   int64
		mock    func()
		want    *model.User
		wantErr bool
	}{
		{
			//When everything works as expected
			name:  "OK",
			s:     &s,
			msgId: 1,
			mock: func() {
				// mock
				rows := sqlmock.NewRows([]string{"Id", "FirstName", "LastName", "Email", "Age"}).AddRow(1, "John", "Doe", "john.doe@gmail.com", 30)

				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "users" WHERE "users"."id" = $1 LIMIT 1`)).
					WithArgs(1).
					WillReturnRows(rows)

			},
			want: &model.User{
				Id:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@gmail.com",
				Age:       30,
			},
		},
		{
			//When the user queried is not found
			name:  "Not Found",
			s:     &s,
			msgId: 1,
			mock: func() {
				// mock
				rows := sqlmock.NewRows([]string{"Id", "FirstName", "LastName", "Email", "Age"}).AddRow(2, "John", "Doe", "john.doe@gmail.com", 30)

				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "users" WHERE "users"."id" = $1 LIMIT 1`)).
					WithArgs(1).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
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
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Test failed: %v", err)
			}
		})
	}
}

func TesUserRepo_Create(t *testing.T) {
	mockDB, mock, err := NewDbMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	s := UserRepository{DB: mockDB}

	tests := []struct {
		name    string
		s       *UserRepository
		request *model.User
		mock    func()
		want    *model.User
		wantErr bool
	}{
		{
			name: "OK",
			s:    &s,
			request: &model.User{
				//Id:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@gmail.com",
				Age:       30,
			},
			mock: func() {

				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO "users" ("first_name","last_name","email","age") VALUES ($1,$2,$3,$4)`)).
					WillReturnRows(rows)
				mock.ExpectCommit()

			},
			want: &model.User{
				Id:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@gmail.com",
				Age:       30,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.DbCreateUser(tt.request)
			if (err != nil) != tt.wantErr {
				fmt.Println("this is the error message: ", err.Message())
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TesUserRepo_Update(t *testing.T) {
	mockDB, mock, err := NewDbMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	s := UserRepository{DB: mockDB}

	tests := []struct {
		name    string
		s       *UserRepository
		request *model.User
		mock    func()
		want    *model.User
		wantErr bool
	}{
		{
			name: "OK",
			s:    &s,
			request: &model.User{
				Id:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@gmail.com",
				Age:       30,
			},
			mock: func() {
				// mock
				sqlmock.NewRows([]string{"Id", "FirstName", "LastName", "Email", "Age"}).AddRow(1, "John", "Doe", "john.doe@gmail.com", 30)

				mock.ExpectBegin()
				mock. /*ExpectPrepare("UPDATE messages").*/ ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "id"=$1,"first_name"=$2,"last_name"=$3,"email"=$4,"age"=$5 WHERE Id = $6`)).WithArgs(1, "Johnny", "Dover", "Johnny.Dover@gmail.com", 37, 1).WillReturnResult(sqlmock.NewResult(0, 1))
				//mock.ExpectQuery(regexp.QuoteMeta(
				//	`UPDATE "users" SET "id"=$1,"first_name"=$2,"last_name"=$3,"email"=$4,"age"=$5 WHERE Id = $6`)) //.WithArgs("Johnny", "Dover", "Johnny.Dover@gmail.com", 37, 1) /*.
				//WillReturnRows(rows)*/
				mock.ExpectCommit()

			},
			want: &model.User{
				Id:        1,
				FirstName: "Johnny",
				LastName:  "Dover",
				Email:     "Johnny.Dover@gmail.com",
				Age:       37,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.DbUpdateUser(tt.request)
			if (err != nil) != tt.wantErr {
				fmt.Println("this is the error message: ", err.Message())
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TesUserRepo_GetAll(t *testing.T) {
	mockDB, mock, err := NewDbMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	s := UserRepository{DB: mockDB}

	tests := []struct {
		name    string
		s       UserRepository
		msgId   int64
		mock    func()
		want    []model.User
		wantErr bool
	}{
		{
			//When everything works as expected
			name: "OK",
			s:    s,
			mock: func() {
				//We added two rows
				rows := sqlmock.NewRows([]string{"Id", "FirstName", "LastName", "Email", "Age"}).AddRow(1, "John", "Doe", "john.doe@gmail.com", 30).AddRow(2, "Johnny", "Dover", "Johnny.Dover@gmail.com", 37)
				mock. /*ExpectPrepare("SELECT (.+) FROM messages").*/ ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).WillReturnRows(rows)
			},
			want: []model.User{
				{
					Id:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@gmail.com",
					Age:       30,
				},
				{
					Id:        2,
					FirstName: "Johnny",
					LastName:  "Dover",
					Email:     "Johnny.Dover@gmail.com",
					Age:       37,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.DbListUsers()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error new = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TesUserRepo_Delete(t *testing.T) {

	mockDB, mock, err := NewDbMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	s := UserRepository{DB: mockDB}

	tests := []struct {
		name    string
		s       *UserRepository
		msgId   int64
		mock    func()
		want    *model.User
		wantErr bool
	}{
		{
			//When everything works as expected
			name:  "OK",
			s:     &s,
			msgId: 1,
			mock: func() {
				mock.ExpectExec(`DELETE FROM "users" WHERE id = $1`).WithArgs(1) //.WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := tt.s.DbDeleteUser(tt.msgId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error new = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
