package main

import (
	"github.com/development-raul/footy-predictor/src/app"
	"github.com/development-raul/footy-predictor/src/docs"
	"os"
)

// @title Footy Predictor API
// @version 1.0
// @description Endpoints details for Footy Predictor API.
// @termsOfService http://swagger.io/terms/

// @contact.name Raul Brindus
// @contact.url http://www.swagger.io/support
// @contact.email raul.brindus@gmail.com

// @license.name Proprietary
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5000
// @BasePath /v1
// @query.collection.format multi

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Swagger 2.0 Meta Information
	docs.SwaggerInfo.Title = "Footy Predictor API"
	docs.SwaggerInfo.Description = "Endpoints details for Footy Predictor API."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:5000"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

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
}
