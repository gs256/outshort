package users

import (
	"net/http"
	"outshort/app/common"
	"strings"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
	storage *Storage
}

func NewUsersController(storage *Storage) *UsersController {
	return &UsersController{
		storage: storage,
	}
}

func (this *UsersController) HandleSignIn(context *gin.Context) {
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
		if err.Code == common.ErrorInvalidCredentials {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else if err.Code == common.ErrorNotFound {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return
	}
	authToken := common.GenerateAuthToken()
	err = this.storage.CreateAuthToken(authToken, userId, common.AuthTokenLifetimeSec)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusAccepted, gin.H{"authToken": authToken})
}

func (this *UsersController) HandleSignUp(context *gin.Context) {
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
		if err.Code == common.ErrorUniqueViolation {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	authToken := common.GenerateAuthToken()
	err = this.storage.CreateAuthToken(authToken, userId, common.AuthTokenLifetimeSec)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusAccepted, gin.H{"authToken": authToken})
}

func (this *UsersController) HandleSignOut(context *gin.Context) {
	token := common.GetAuthTokenFromHeader(context)
	_ = this.storage.DeleteAuthToken(token)
	context.Status(http.StatusAccepted)
}

func (this *UsersController) HandleGetUserInfo(context *gin.Context) {
	user := context.MustGet("user").(UserModel)
	context.JSON(http.StatusOK, gin.H{
		"username": user.Username,
	})
}
