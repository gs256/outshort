package main

import (
	"math/rand/v2"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

type ApiController struct {
	storage *Storage
}

func (this *ApiController) Initialize(storage *Storage) {
	this.storage = storage
}

func (this *ApiController) HandleTestGet(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (this *ApiController) HandleShortenPost(context *gin.Context) {
	var req ShortenRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body format"})
		return
	}
	for true {
		alias := generateAlias()
		_, err := this.storage.CreateLink(req.Url, alias)
		if err != nil && err.code == AliasAlreadyExists {
			continue
		}
		context.JSON(http.StatusAccepted, gin.H{"alias": alias})
		break
	}
}
