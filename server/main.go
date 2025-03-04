package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		AllowAllOrigins:  true,
		AllowCredentials: true,
	}))

	storage := Storage{}
	storage.Initialize()
	defer storage.Dispose()

	apiController := ApiController{}
	apiController.Initialize(&storage)

	router.GET("/api/v1/test", apiController.HandleTest)
	router.GET("/api/v1/redirect/:alias", apiController.HandleRedirect)
	router.POST("/api/v1/auth/sign-in", apiController.HandleSignIn)
	router.POST("/api/v1/auth/sign-up", apiController.HandleSignUp)
	router.POST("/api/v1/auth/sign-out", apiController.HandleSignOut)
	router.GET("/api/v1/auth/user-info", apiController.HandleGetUserInfo)
	router.POST("/api/v1/links/quick-shorten", apiController.HandleShorten)

	router.Run(":8249")
}
