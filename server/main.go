package main

import (
	"outshort/app/links"
	"outshort/app/users"

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

	usersController, disposeUsers := createUsersController()
	linksController, disposeLinks := createLinksController()

	router.POST("/api/v1/auth/sign-in", usersController.HandleSignIn)
	router.POST("/api/v1/auth/sign-up", usersController.HandleSignUp)
	router.POST("/api/v1/auth/sign-out", usersController.HandleSignOut)
	router.GET("/api/v1/auth/user-info", usersController.HandleGetUserInfo)

	router.GET("/api/v1/redirect/:alias", linksController.HandleRedirect)
	router.POST("/api/v1/links/quick-shorten", linksController.HandleQuickShorten)
	router.POST("/api/v1/links/create", linksController.HandleLinkCreate)
	router.POST("/api/v1/links/update/:uid", linksController.HandleLinkUpdate)
	router.GET("/api/v1/links/all", linksController.HandleLinksGetAll)

	router.Run(":8249")

	disposeUsers()
	disposeLinks()
}

func createUsersController() (*users.UsersController, func()) {
	usersStorage := users.Storage{}
	usersStorage.Initialize()

	usersController := users.UsersController{}
	usersController.Initialize(&usersStorage)

	return &usersController, usersStorage.Dispose
}

func createLinksController() (*links.LinksController, func()) {
	linksStorage := links.Storage{}
	linksStorage.Initialize()

	linksController := links.LinksController{}
	linksController.Initialize(&linksStorage)

	return &linksController, linksStorage.Dispose
}
