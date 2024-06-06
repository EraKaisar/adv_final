// internal/data/models
package data

import (
	"EPLgateway/auth-service/internal/data"
	"database/sql"
	_ "errors"
	"log"
	"os"
)

type Models struct {
	Users       data.UserModel
	Tokens      data.TokenModel
	Permissions data.PermissionModel
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{

		Users: data.UserModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Tokens: data.TokenModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Permissions: data.PermissionModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
	}
}
