package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	storage := Storage{}
	storage.Initialize()
	defer storage.Dispose()

	apiController := ApiController{}
	apiController.Initialize(&storage)

	router.GET("/api/v1/test", apiController.HandleTest)
	router.POST("/api/v1/shorten", apiController.HandleShorten)
	router.GET("/api/v1/redirect/:alias", apiController.HandleRedirect)

	router.Run()
}
