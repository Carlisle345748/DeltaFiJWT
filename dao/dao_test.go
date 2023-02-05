package dao

import (
	"deltaFiJWT/types"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
	"testing"
)

func setupTest(t *testing.T) {
	if err := InitTestDB(); err != nil {
		t.Error(err)
	}
	if err := DB.Migrator().DropTable(&User{}); err != nil {
		t.Error(err)
	}
	if err := DB.AutoMigrate(&User{}); err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		if err := DB.Migrator().DropTable(&User{}); err != nil {
			t.Error(err)
		}
	})
}

func TestCreateUser(t *testing.T) {
	setupTest(t)
	input := types.CreateUserInput{
		Email:     "jackson@example.com",
		FirstName: "jackson",
		LastName:  "wang",
		Password:  "123",
	}

	// Test creat user
	user, err := CreateUser(input)
	if err != nil {
		t.Error(err)
	}
	if user.Email != input.Email || user.FirstName != input.FirstName ||
		user.LastName != input.LastName || user.Password != input.Password {
		t.Error("create user failed: data is not match")
	}

	got := &User{}
	result := DB.First(got, user.ID)
	if result.Error != nil {
		t.Error(result.Error)
	}

	if got.Email != input.Email || got.FirstName != input.FirstName ||
		got.LastName != input.LastName || got.Password != input.Password {
		t.Error("create user failed: data is not match")
	}

	// Test create user with the same email
	_, err = CreateUser(input)
	if err.(sqlite3.Error).Code != sqlite3.ErrConstraint {
		t.Error("unique constraint failed")
	}
}

func TestUpdateUser(t *testing.T) {
	setupTest(t)
	user := &User{
		Email:     "jackson@example.com",
		FirstName: "jackson",
		LastName:  "wang",
		Password:  "123",
	}
	result := DB.Create(user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	input := types.UpdateUserInput{
		ID:        user.ID,
		FirstName: "amy",
		LastName:  "june",
	}
	err := UpdateUser(input)
	if err != nil {
		t.Error(err)
	}

	got := &User{}
	result = DB.First(got, user.ID)
	if result.Error != nil {
		t.Error(result.Error)
	}

	if got.FirstName != input.FirstName || got.LastName != input.LastName {
		t.Error("update user failed: data is not match")
	}
}

func TestDeleteUser(t *testing.T) {
	setupTest(t)
	user := &User{
		Email:     "jackson@example.com",
		FirstName: "jackson",
		LastName:  "wang",
		Password:  "123",
	}
	result := DB.Create(user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	err := DeleteUser(user.ID)
	if err != nil {
		t.Error(err)
	}

	got := &User{}
	result = DB.First(got, user.ID)
	if result.Error != gorm.ErrRecordNotFound {
		t.Errorf("delete user failed: got error is %v", result.Error)
	}
}

func TestGetUser(t *testing.T) {
	setupTest(t)
	user := &User{
		Email:     "jackson@example.com",
		FirstName: "jackson",
		LastName:  "wang",
		Password:  "123",
	}
	result := DB.Create(user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	got, err := GetUser(user.ID)
	if err != nil {
		t.Error(err)
	}

	if got.ID != user.ID {
		t.Error("get user failed: data is not match")
	}
}
