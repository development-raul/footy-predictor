package app

import (
	"fmt"
	"github.com/development-raul/footy-predictor/src/clients/mysql/footy_db"
	"github.com/development-raul/footy-predictor/src/zlog"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type App struct {
	FootyDB *sqlx.DB
	Router  *gin.Engine
}

type Credentials struct {
	DBUser  string
	DBPass  string
	DBHost  string
	DBPort  string
	DBName  string
	AppPort string
}

func StartApplication(c *Credentials) {
	// Connect to DB
	footyDB := footy_db.ConnectToDatabase(c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
	application := &App{
		FootyDB: footyDB,
		Router:  gin.Default(),
	}

	application.SetupRoutes()

	if err := application.Router.Run(fmt.Sprintf(":%s", c.AppPort)); err != nil {
		zlog.Logger.Panicw("failed to start application", "error", err)
	}
}
