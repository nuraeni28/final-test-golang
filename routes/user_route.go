// routes/user_routes.go
package routes

import (
	"MyGram/controllers"
	"MyGram/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine) {
	userGroup := r.Group("/users")

	userGroup.POST("/register", controllers.Register)
	userGroup.POST("/login", controllers.Login)

	authorized := userGroup.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.PUT("/:userId", controllers.UpdateUser)
		authorized.DELETE("/:userId", controllers.DeleteUser)
	}

}
