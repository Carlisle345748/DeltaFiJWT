package main

import (
	"log"
	"net/http"

	"deltaFiJWT/controller"
	"deltaFiJWT/dao"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	if err := dao.InitDB(); err != nil {
		log.Fatalf("initilize databse fialed: %v", err)
	}

	// Create server
	r := gin.Default()
	r.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Page not found",
		})
	})
	controller.SetUpRouter(r)

	// Run server
	if err := r.Run(); err != nil {
		log.Fatalf("start server failed %v", err)
	}
}
