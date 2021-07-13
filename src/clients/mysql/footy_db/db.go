package footy_db

import (
	"fmt"
	"github.com/development-raul/footy-predictor/src/zlog"
	"github.com/jmoiron/sqlx"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var (
	Client *sqlx.DB
)

func ConnectToDatabase(dbUser, dbPass, dbHost, dbPort, dbName string) *sqlx.DB {
	var err error
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	Client, err = sqlx.Connect("mysql", connectionString)
	if err != nil {
		zlog.Logger.Panicw("database connection failed", "error", err)
	}
	Client.SetMaxIdleConns(50)
	Client.SetMaxOpenConns(100)
	Client.SetConnMaxLifetime(time.Second * 10)

	return Client
}
