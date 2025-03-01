package main

import (
	"math/rand/v2"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

const AliasLength = 5
const StringGenerationAlphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const AuthTokenLifetimeSec = 1 * 60 * 60

type ShortenRequest struct {
	Url string `json:"url"`
}

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func randomString(length int) string {
	result := ""
	for range length {
		randomIndex := randRange(0, len(StringGenerationAlphabet)-1)
		result += string(StringGenerationAlphabet[randomIndex])
	}
	return result
}

func randRange(min int, max int) int {
	return rand.IntN(max-min) + min
}

func validateUrl(sourceUrl string) (string, bool) {
	parsed, err := url.ParseRequestURI(sourceUrl)
	if err != nil {
		return sourceUrl, false
	}
	return parsed.String(), true
}

func generateAuthToken() string {
	return randomString(32)
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
	originalUrl, urlValid := validateUrl(req.Url)
	if !urlValid {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid url"})
		return
	}
	for true {
		alias := randomString(AliasLength)
		_, err := this.storage.CreateLink(originalUrl, alias)
		if err != nil && err.code == UniqueViolation {
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

func (this *ApiController) HandleSignIn(context *gin.Context) {
}

func (this *ApiController) HandleSignUp(context *gin.Context) {
	var req SignUpRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body format"})
		return
	}
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	if len(username) < 2 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Username must be at least 2 characters"})
		return
	}
	if len(password) < 6 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters"})
		return
	}
	userId, err := this.storage.CreateUser(username, password)
	if err != nil {
		if err.code == UniqueViolation {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	authToken := generateAuthToken()
	err = this.storage.CreateAuthToken(authToken, userId, AuthTokenLifetimeSec)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusAccepted, gin.H{"authToken": authToken})
}
