package main

import (
	"math/rand/v2"
	"net/http"

	"github.com/gin-gonic/gin"
)

const AliasLength = 5
const AliasAlphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type ShortenRequest struct {
	Url string `json:"url"`
}

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

func (this *ApiController) HandleTest(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "test",
	})
}

func (this *ApiController) HandleShorten(context *gin.Context) {
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

func (this *ApiController) HandleRedirect(context *gin.Context) {
	alias := context.Param("alias")
	if alias == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Alias parameter required"})
	}
	location, err := this.storage.GetOriginalUrl(alias)
	if err != nil {
		if err.code == NotFound {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Alias not found"})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return
	}
	context.Redirect(http.StatusFound, location)
}
