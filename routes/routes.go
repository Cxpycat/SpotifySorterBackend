package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/code", controllers.SendCode)
	}

	playlistRoutes := router.Group("/playlists")
	{
		playlistRoutes.GET("/", controllers.GetPlaylists)
		playlistRoutes.POST("/", controllers.CreatePlaylist)
	}

	return router
}
