package controllers

import (
	"database/sql"
)

type BaseController struct {
	Db *sql.DB

}
