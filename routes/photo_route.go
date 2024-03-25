// routes/user_routes.go
package routes

import (
	"MyGram/controllers"
	"MyGram/middleware"

	"github.com/gin-gonic/gin"
)

func SetupPhotoRoutes(r *gin.Engine) {
	userGroup := r.Group("/photos")

	authorized := userGroup.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/", controllers.CreatePost)
		authorized.GET("/", controllers.GetPosts)
		authorized.PUT("/:photoID", controllers.EditPhoto)
		authorized.DELETE("/:photoID", controllers.DeletePhoto)
	}

}
