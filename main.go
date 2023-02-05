package main

import (
	"deltaFiJWT/jwt"
	"log"
	"net/http"

	"deltaFiJWT/controller"
	"deltaFiJWT/dao"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := dao.InitDB(); err != nil {
		log.Fatalf("initilize databse fialed: %v", err)
	}

	r := gin.Default()

	authMiddleware, err := jwt.CreateJWTMiddleware()
	if err != nil {
		log.Fatalf("JWT middleware intilization error: %v", err)
	}

	r.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Page not found",
		})
	})

	r.POST("/login", authMiddleware.LoginHandler)
	r.GET("/logout", authMiddleware.LogoutHandler)

	auth := r.Group("/")
	auth.Use(authMiddleware.MiddlewareFunc())
	auth.GET("/refreshToken", authMiddleware.RefreshHandler)
	auth.GET("/hello", controller.HelloHandler)
	auth.PUT("/user", controller.CreateUserHandler)
	auth.POST("/user", controller.UpdateUserHandler)
	auth.DELETE("/user", controller.DeleteUserHandler)

	if err := r.Run(); err != nil {
		log.Fatalf("start server failed %v", err)
	}
}
