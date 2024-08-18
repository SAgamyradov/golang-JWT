package routes

import (
	"go-jwt/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("user/signup", controllers.Signup())
	incomingRoutes.GET("user/login", controllers.Login())
}
