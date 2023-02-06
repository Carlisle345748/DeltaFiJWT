package dao

import (
	"testing"

	"deltaFiJWT/types"

	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) {
	assert.NoError(t, InitTestDB())

	assert.NoError(t, DB.Migrator().DropTable(&User{}))

	assert.NoError(t, DB.AutoMigrate(&User{}))

	t.Cleanup(func() {
		assert.NoError(t, DB.Migrator().DropTable(&User{}))
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
	assert.NoError(t, err)
	assert.Equal(t, user.Email, input.Email)
	assert.Equal(t, user.FirstName, input.FirstName)
	assert.Equal(t, user.LastName, input.LastName)
	assert.Equal(t, user.Password, input.Password)

	got := &User{}
	result := DB.First(got, user.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, got.Email, input.Email)
	assert.Equal(t, got.FirstName, input.FirstName)
	assert.Equal(t, got.LastName, input.LastName)
	assert.Equal(t, got.Password, input.Password)

	// Test create user with the same email
	_, err = CreateUser(input)
	assert.Error(t, err)
	assert.Equal(t, err.(sqlite3.Error).Code, sqlite3.ErrConstraint)
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
	assert.NoError(t, result.Error)

	input := types.UpdateUserInput{
		ID:        user.ID,
		FirstName: "amy",
		LastName:  "june",
	}
	err := UpdateUser(input)
	assert.NoError(t, err)

	got := &User{}
	result = DB.First(got, user.ID)
	assert.NoError(t, result.Error)

	assert.Equal(t, got.FirstName, input.FirstName)
	assert.Equal(t, got.LastName, input.LastName)
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
	assert.NoError(t, result.Error)

	err := DeleteUser(user.ID)
	assert.NoError(t, err)

	got := &User{}
	result = DB.First(got, user.ID)
	assert.Error(t, result.Error)
	assert.Equal(t, result.Error, gorm.ErrRecordNotFound)
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
	assert.NoError(t, result.Error)

	got, err := GetUser(user.ID)
	assert.NoError(t, err)

	assert.Equal(t, got.ID, user.ID)
}
