package controller

import (
	"fmt"
	"net/http"

	"deltaFiJWT/dao"
	"deltaFiJWT/types"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess        = 0
	CodeUserNotFound   = 1
	CodeInvalidInput   = 2
	CodeCreateUserFail = 3
	CodeUpdateUserFail = 4
	CodeDeleteUserFail = 6
)

func HelloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userId := claims[jwt.IdentityKey]
	user, err := dao.GetUser(userId.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    CodeUserNotFound,
			"message": fmt.Sprintf("find user fail: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     CodeSuccess,
		"message":  "success",
		"greeting": fmt.Sprintf("Hello %s %s", user.FirstName, user.LastName),
	})
}

func CreateUserHandler(c *gin.Context) {
	input := types.CreateUserInput{}
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    CodeInvalidInput,
			"message": err.Error(),
		})
		return
	}
	user, err := dao.CreateUser(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    CodeCreateUserFail,
			"message": fmt.Sprintf("create user failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    CodeSuccess,
		"message": "success",
		"user":    user,
	})
}

func UpdateUserHandler(c *gin.Context) {
	input := types.UpdateUserInput{}
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    CodeInvalidInput,
			"message": err.Error(),
		})
		return
	}

	if err := dao.UpdateUser(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    CodeUpdateUserFail,
			"message": fmt.Sprintf("update user failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    CodeSuccess,
		"message": "success",
	})
}

func DeleteUserHandler(c *gin.Context) {
	input := types.DeleteUserInput{}
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    CodeInvalidInput,
			"message": err.Error(),
		})
		return
	}

	if err := dao.DeleteUser(input.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    CodeDeleteUserFail,
			"message": fmt.Sprintf("delete user failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    CodeSuccess,
		"message": "success",
	})
}