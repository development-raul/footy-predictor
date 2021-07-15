package controllers

import "github.com/gin-gonic/gin"

type countryControllerInterface interface {
	SyncCountries(ctx *gin.Context)
}

type countryController struct{}

var CountryController countryControllerInterface = &countryController{}

func (c *countryController) SyncCountries(ctx *gin.Context) {

}
