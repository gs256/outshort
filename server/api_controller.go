package main

import (
	"math/rand/v2"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const AliasLength = 5
const StringGenerationAlphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const AuthTokenLifetimeSec = 1 * 60 * 60

type ShortenRequest struct {
	Url string `json:"url"`
}

type CreateLinkRequest struct {
	Url      string `json:"url"`
	Name     string `json:"name"`
	Alias    string `json:"alias"`
	Lifetime int    `json:"lifetime"`
}

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Link struct {
	Uid          string    `json:"uid"`
	Alias        string    `json:"alias"`
	OriginalUrl  string    `json:"originalUrl"`
	Name         string    `json:"name"`
	Lifetime     int       `json:"lifetime"`
	CreationDate time.Time `json:"creationDate"`
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
	return RandomString(32)
}

func getAuthTokenFromHeader(context *gin.Context) string {
	authHeader := context.GetHeader("Authorization")
	split := strings.Split(authHeader, " ")
	if len(split) != 2 {
		return ""
	}
	if strings.ToLower(split[0]) != "bearer" {
		return ""
	}
	token := strings.TrimSpace(split[1])
	return token
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

func (this *ApiController) HandleQuickShorten(context *gin.Context) {
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
		alias := RandomString(AliasLength)
		_, err := this.storage.CreateQuickLink(originalUrl, alias)
		if err != nil && err.code == UniqueViolation {
			continue
		}
		context.JSON(http.StatusAccepted, gin.H{"alias": alias})
		break
	}
}

func (this *ApiController) HandleLinkCreate(context *gin.Context) {
	token := getAuthTokenFromHeader(context)
	user, err := this.storage.GetUserInfo(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var req CreateLinkRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body format"})
		return
	}
	if req.Alias != "" {
		if len(req.Alias) < 5 {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alias"})
			return
		}
		exists, err := this.storage.AliasAlreadyExists(req.Alias)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
			return
		}
		if exists {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Alias already exists"})
			return
		}
	}
	originalUrl, urlValid := validateUrl(req.Url)
	if !urlValid {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid url"})
		return
	}
	if req.Lifetime < 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lifetime"})
		return
	}
	if req.Alias == "" {
		for true {
			newAlias := RandomString(AliasLength)
			exists, err := this.storage.AliasAlreadyExists(req.Alias)
			if err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
				return
			}
			if exists {
				continue
			}
			req.Alias = newAlias
			break
		}
	}
	linkModel, err := this.storage.CreateLink(originalUrl, req.Name, req.Alias, req.Lifetime, user.Id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	link := ToLink(*linkModel)
	context.JSON(http.StatusAccepted, link)
}

func (this *ApiController) HandleLinksGetAll(context *gin.Context) {
	token := getAuthTokenFromHeader(context)
	user, err := this.storage.GetUserInfo(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	linkModels, err := this.storage.GetAllLinks(user.Id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user links"})
		return
	}
	links := ToLinks(linkModels)
	context.JSON(http.StatusOK, links)
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
	var req SignInRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body format"})
		return
	}
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	if len(username) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body format"})
		return
	}
	userId, err := this.storage.AuthenticateUser(username, password)
	if err != nil {
		if err.code == InvalidCredentials {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else if err.code == NotFound {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return
	}
	authToken := generateAuthToken()
	err = this.storage.CreateAuthToken(authToken, userId, AuthTokenLifetimeSec)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusAccepted, gin.H{"authToken": authToken})
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

func (this *ApiController) HandleSignOut(context *gin.Context) {
	token := getAuthTokenFromHeader(context)
	_ = this.storage.DeleteAuthToken(token)
	context.Status(http.StatusAccepted)
}

func (this *ApiController) HandleGetUserInfo(context *gin.Context) {
	token := getAuthTokenFromHeader(context)
	user, err := this.storage.GetUserInfo(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"username": user.Username,
	})
}
