package main

import (
	"github.com/gin-gonic/gin"
)

type ShortenRequest struct {
	Url string `json:"url"`
}

func main() {
	router := gin.Default()

	storage := Storage{}
	storage.Initialize()
	defer storage.Dispose()

	apiController := ApiController{}
	apiController.Initialize(&storage)

	router.GET("/api/v1/test", apiController.HandleTestGet)
	router.POST("/api/v1/shorten", apiController.HandleShortenPost)

	router.Run()
}
