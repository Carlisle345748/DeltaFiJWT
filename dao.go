package main

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct {
	gorm.Model
	Email     string
	FirstName string
	LastName  string
	Password  string
}

func InitDB() error {
	instance, err := gorm.Open(sqlite.Open("deltaFi.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	db = instance

	if err := db.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}

func Authenticate(email, password string) (*User, error) {
	user := &User{}
	db.First(user, "email = ?", email)
	if user.Password != password {
		return nil, errors.New("incorrect email or password")
	}
	return user, nil
}
