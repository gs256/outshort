package main

import (
	"fmt"
	"net/http"

	"math/rand/v2"

	"github.com/gin-gonic/gin"
)

type ShortenRequest struct {
	Url string `json:"url"`
}

func main() {
	router := gin.Default()

	router.GET("/api/v1/test", handleTestGet)
	router.POST("/api/v1/shorten", handleShortenPost)

	router.Run()
}

func handleTestGet(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func handleShortenPost(context *gin.Context) {
	var req ShortenRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body format"})
		return
	}
	alias := generateAlias()
	fmt.Printf("shortening: %s (alias is `%s`)\n", req.Url, alias)
}

const AliasLength = 5
const AliasAlphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateAlias() string {
	result := ""
	for i := 0; i < AliasLength; i++ {
		randomIndex := randRange(0, len(AliasAlphabet)-1)
		result += string(AliasAlphabet[randomIndex])
	}
	return result
}

func randRange(min int, max int) int {
	return rand.IntN(max-min) + min
}
