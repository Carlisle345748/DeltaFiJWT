package dao

import (
	"errors"

	"deltaFiJWT/types"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Email     string `gorm:"unique;not null" json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `gorm:"not null" json:"-"`
}

// InitDB initialize database
func InitDB() error {
	instance, err := gorm.Open(sqlite.Open("deltaFi.DB"), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = instance
	if err := DB.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}

// InitTestDB initialize test database
func InitTestDB() error {
	instance, err := gorm.Open(sqlite.Open("deltaFi_test.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = instance
	return nil
}

// Authenticate user by email and password
func Authenticate(email, password string) (*User, error) {
	user := &User{}
	result := DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("incorrect email or password")
		}
		return nil, result.Error
	}
	if user.Password != password {
		return nil, errors.New("incorrect email or password")
	}
	return user, nil
}

// CreateUser create new user. Return error if the email is already exist
func CreateUser(input types.CreateUserInput) (*User, error) {
	user := User{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Password:  input.Password,
	}
	result := DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUser query a specific user by ID
func GetUser(ID uint) (*User, error) {
	user := &User{}
	result := DB.First(user, ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// UpdateUser update user firstname and lastname
func UpdateUser(input types.UpdateUserInput) error {
	result := DB.Model(&User{}).Where("id = ?", input.ID).Updates(map[string]interface{}{
		"first_name": input.FirstName,
		"last_name":  input.LastName,
	})
	return result.Error
}

// DeleteUser delete user by ID
func DeleteUser(ID uint) error {
	result := DB.Delete(&User{}, ID)
	return result.Error
}
