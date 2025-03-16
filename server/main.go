package main

import (
	"outshort/app/common"
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

	dbConnection := common.NewDbConnection()
	defer dbConnection.Close()

	usersController, userStorage := createUsersController(dbConnection)
	linksController := createLinksController(dbConnection)

	authRequired := users.AuthRequired(userStorage)

	router.POST("/api/v1/auth/sign-in", usersController.HandleSignIn)
	router.POST("/api/v1/auth/sign-up", usersController.HandleSignUp)
	router.POST("/api/v1/auth/sign-out", usersController.HandleSignOut)

	router.GET("/api/v1/users/me", authRequired, usersController.HandleGetUserInfo)

	router.GET("/api/v1/redirect/:alias", linksController.HandleRedirect)
	router.POST("/api/v1/links/quick-shorten", linksController.HandleQuickShorten)

	router.POST("/api/v1/links/create", authRequired, linksController.HandleLinkCreate)
	router.POST("/api/v1/links/update/:uid", authRequired, linksController.HandleLinkUpdate)
	router.GET("/api/v1/links/all", authRequired, linksController.HandleLinksGetAll)

	router.Run(":8249")
}

func createUsersController(dbConnection *common.DbConnection) (*users.UsersController, *users.Storage) {
	usersStorage := users.Storage{}
	usersStorage.Initialize(dbConnection)

	usersController := users.UsersController{}
	usersController.Initialize(&usersStorage)

	return &usersController, &usersStorage
}

func createLinksController(dbConnection *common.DbConnection) *links.LinksController {
	linksStorage := links.Storage{}
	linksStorage.Initialize(dbConnection)

	linksController := links.LinksController{}
	linksController.Initialize(&linksStorage)

	return &linksController
}
