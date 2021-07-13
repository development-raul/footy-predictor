package main

import (
	"github.com/development-raul/footy-predictor/src/app"
	"github.com/development-raul/footy-predictor/src/docs"
	"github.com/development-raul/footy-predictor/src/zlog"
	"os"
)

// @securityDefinitions.apikey bearerAuth
// @in header
// @name Authorization
func main() {
	// Swagger 2.0 Meta Information
	docs.SwaggerInfo.Title = "EMPRIS"
	docs.SwaggerInfo.Description = "EMPRIS - Go API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:5000"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	appCred := app.Credentials{
		DBName:  os.Getenv("DB_NAME"),
		DBUser:  os.Getenv("DB_USER"),
		DBPass:  os.Getenv("DB_PASS"),
		DBHost:  os.Getenv("DB_HOST"),
		DBPort:  os.Getenv("DB_PORT"),
		AppPort: os.Getenv("APP_PORT"),
	}

	app.StartApplication(&appCred)
	// Start application will panic so if we got here all is good :)
	zlog.Logger.Infow("application started")
}
