package app

import (
	"outshort/app/common"
	"outshort/app/links"
	"outshort/app/users"
)

type AppContext struct {
	DbConnection *common.DbConnection

	LinkStorage *links.Storage
	UserStorage *users.Storage

	LinksController *links.LinksController
	UsersController *users.UsersController
}

func NewAppContext() *AppContext {
	dbConnection := common.NewDbConnection()
	usersStorage := users.NewStorage(dbConnection)
	usersController := users.NewUsersController(usersStorage)
	linksStorage := links.NewStorage(dbConnection)
	linksController := links.NewLinksController(linksStorage)

	return &AppContext{
		DbConnection: dbConnection,

		LinkStorage: linksStorage,
		UserStorage: usersStorage,

		LinksController: linksController,
		UsersController: usersController,
	}
}

func (this *AppContext) Dispose() {
	this.DbConnection.Close()
}
