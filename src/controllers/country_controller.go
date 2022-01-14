package controllers

import (
	"github.com/development-raul/footy-predictor/src/domains/countries"
	"github.com/development-raul/footy-predictor/src/services"
	"github.com/development-raul/footy-predictor/src/swaggertypes"
	"github.com/development-raul/footy-predictor/src/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type countryControllerInterface interface {
	Create(ctx *gin.Context)
	SyncCountries(ctx *gin.Context)
}

type countryController struct{}

var CountryController countryControllerInterface = &countryController{}

// Create
// @Summary Create country
// @Description Endpoint used to create a new country record
// @ID v1-countries-create
// @Produce json
// @Accept json
// @Tags Countries
// @Param JSON request body countries.CreateCountryInput true "Request Sample"
// @Success 201 {object} swaggertypes.NoErrorString
// @Failure 400 {object} swaggertypes.StandardBadRequestError
// @Failure 401 {object} swaggertypes.StandardUnauthorisedError
// @Failure 500 {object} swaggertypes.StandardInternalServerError
// @Router /countries [post]
func (c *countryController) Create(ctx *gin.Context) {
	var req countries.CreateCountryInput
	if ok := utils.GinShouldPassAll(ctx, utils.GinShouldBind(&req), utils.GinShouldValidate(&req)); !ok {
		return
	}

	if err := services.CountryService.Create(&req); err != nil {
		ctx.JSON(err.Code(), err)
		return
	}

	ctx.JSON(http.StatusCreated, swaggertypes.NoErrorString{
		Message: "SUCCESS",
		Code: http.StatusCreated,
	})
}

func (c *countryController) SyncCountries(ctx *gin.Context) {

}


