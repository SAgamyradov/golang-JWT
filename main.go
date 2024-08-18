package main

import (
	"go-jwt/db"
	"go-jwt/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()
	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRouter(router)
	routes.UserRouter(router)

	router.GET("/api-1", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"success": "Granted for api-1"})
	})
	router.GET("/api-2", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"success": "Granted for api-2"})
	})

	router.Run(":8080")
}
