package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	authMiddleware, err := CreateJWTMiddleware()
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

	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.GET("/hello", HelloHandler)

	if err := r.Run(); err != nil {
		log.Fatalf("start server failed %v", err)
	}
}
