package db

import (
	"database/sql"
)

var Db *sql.DB

func Connect(){
	var err error
	Db, err = sql.Open("postgres", "postgres://stylianosgiakoumidis:asdASDqwe123!!!@localhost/journal?sslmode=disable")
	if err != nil {
		panic("Database connection failed. Error: " + err.Error())
	}

	if err = Db.Ping(); err != nil{
		panic("Database cannot be pinged. Error: " + err.Error())
	}
}

