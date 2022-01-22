package app

import (
	"github.com/development-raul/footy-predictor/src/controllers"
	"net/http"

	"github.com/gin-gonic/gin"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func (app *App) SetupRoutes() {
	// Adding CORS
	app.Router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "https://localhost:8080")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, token, tableau_token")

		c.Next()
	})

	app.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app.Router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not found",
			"code":  404,
		})
	})

	v1Routes := app.Router.Group("/v1")

	v1Routes.GET("/", controllers.HealthController.Check)
	countryGroup := v1Routes.Group("/countries")
	{
		countryGroup.POST("", controllers.CountryController.Create)
		countryGroup.PUT("/:id", controllers.CountryController.Update)
		countryGroup.GET("", controllers.CountryController.List)
		countryGroup.GET("/:id", controllers.CountryController.Find)
		countryGroup.DELETE("/:id", controllers.CountryController.Delete)
	}
}
