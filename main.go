// main.go
package main

import (
	"MyGram/config"
	"MyGram/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.StartDB()
	routes.SetupUserRoutes(r)
	routes.SetupPhotoRoutes(r)
	routes.SetupCommentRoutes(r)
	routes.SetupSocialRoutes(r)

	// Run the server
	r.Run(":8080")
}
