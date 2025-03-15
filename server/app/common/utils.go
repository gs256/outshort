package common

import (
	"math/rand/v2"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

const StringGenerationAlphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randRange(min int, max int) int {
	return rand.IntN(max-min) + min
}

func RandomString(length int) string {
	result := ""
	for range length {
		randomIndex := randRange(0, len(StringGenerationAlphabet)-1)
		result += string(StringGenerationAlphabet[randomIndex])
	}
	return result
}

func ValidateUrl(sourceUrl string) (string, bool) {
	parsed, err := url.ParseRequestURI(sourceUrl)
	if err != nil {
		return sourceUrl, false
	}
	return parsed.String(), true
}

func GetAuthTokenFromHeader(context *gin.Context) string {
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

func GenerateLinkUid() string {
	return RandomString(16)
}
