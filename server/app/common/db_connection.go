package common

import "database/sql"

const DbFile string = "env/database.db"

type DbConnection struct {
	db *sql.DB
}

func (this *DbConnection) Initialize() {
	db, err := sql.Open("sqlite3", DbFile)
	if err != nil {
		panic(err.Error())
	}
	this.db = db
}

func (this *DbConnection) Close() {
	this.db.Close()
}

func (this *DbConnection) Database() *sql.DB {
	return this.db
}

func NewDbConnection() *DbConnection {
	dbConnection := DbConnection{}
	dbConnection.Initialize()
	return &dbConnection
}
