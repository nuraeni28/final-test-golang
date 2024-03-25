// routes/user_routes.go
package routes

import (
	"MyGram/controllers"
	"MyGram/middleware"

	"github.com/gin-gonic/gin"
)

func SetupSocialRoutes(r *gin.Engine) {
	userGroup := r.Group("/socialmedias")

	authorized := userGroup.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/", controllers.CreateSocialMediaProfile)
		authorized.GET("/", controllers.GetSocialMediaWithUser)
		authorized.PUT("/:socialmediaID", controllers.UpdateSocialMedia)
		authorized.DELETE("/:socialmediaID", controllers.DeleteSocialMedia)
	}

}
