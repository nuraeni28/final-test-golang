// routes/user_routes.go
package routes

import (
	"MyGram/controllers"
	"MyGram/middleware"

	"github.com/gin-gonic/gin"
)

func SetupCommentRoutes(r *gin.Engine) {
	userGroup := r.Group("/comments")

	authorized := userGroup.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/", controllers.CreateComment)
		authorized.GET("/", controllers.GetComment)
		authorized.PUT("/:commentId", controllers.UpdateComment)
		authorized.DELETE("/:commentId", controllers.DeleteComment)
	}

}
