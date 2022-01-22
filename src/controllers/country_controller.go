package controllers

import (
	"github.com/development-raul/footy-predictor/src/domains/countries"
	"github.com/development-raul/footy-predictor/src/services"
	"github.com/development-raul/footy-predictor/src/swaggertypes"
	"github.com/development-raul/footy-predictor/src/utils"
	"github.com/development-raul/footy-predictor/src/utils/resterror"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type countryControllerInterface interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Find(ctx *gin.Context)
	List(ctx *gin.Context)
	Delete(ctx *gin.Context)
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
// @Param JSON request body countries.CountryInput true "Request Sample"
// @Success 201 {object} swaggertypes.NoErrorString
// @Failure 400 {object} swaggertypes.StandardBadRequestError
// @Failure 401 {object} swaggertypes.StandardUnauthorisedError
// @Failure 500 {object} swaggertypes.StandardInternalServerError
// @Router /countries [post]
func (c *countryController) Create(ctx *gin.Context) {
	var req countries.CountryInput
	if ok := utils.GinShouldPassAll(ctx, utils.GinShouldBind(&req), utils.GinShouldValidate(&req)); !ok {
		return
	}

	if err := services.CountryService.Create(&req); err != nil {
		ctx.JSON(err.Code(), err)
		return
	}

	ctx.JSON(http.StatusCreated, swaggertypes.NoErrorString{
		Message: "SUCCESS",
		Code:    http.StatusCreated,
	})
}

// Update
// @Summary Update country
// @Description Endpoint used to update an existing country record
// @ID v1-countries-update
// @Produce json
// @Accept json
// @Tags Countries
// @Param id path int true "Country ID"
// @Param JSON request body countries.UpdateCountryInput true "Request Sample"
// @Success 200 {object} swaggertypes.NoErrorString
// @Failure 400 {object} swaggertypes.StandardBadRequestError
// @Failure 401 {object} swaggertypes.StandardUnauthorisedError
// @Failure 500 {object} swaggertypes.StandardInternalServerError
// @Router /countries/{id} [put]
func (c *countryController) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiErr := resterror.NewBadRequestError("INVALID_COUNTRY_ID")
		ctx.JSON(apiErr.Code(), apiErr)
		return
	}

	var req countries.UpdateCountryInput
	if ok := utils.GinShouldPassAll(ctx, utils.GinShouldBind(&req), utils.GinShouldValidate(&req)); !ok {
		return
	}

	if err := services.CountryService.Update(&req, id); err != nil {
		ctx.JSON(err.Code(), err)
		return
	}

	ctx.JSON(http.StatusOK, swaggertypes.NoErrorString{
		Message: "SUCCESS",
		Code:    http.StatusOK,
	})
}

// Find
// @Summary Find country
// @Description Retrieve a country identified by id
// @ID v1-countries-find
// @Produce json
// @Tags Countries
// @Param id path int true "Country ID"
// @Success 200 {object} swaggertypes.NoErrorI{data=countries.CountryOutput}
// @Failure 400 {object} swaggertypes.StandardBadRequestError
// @Failure 401 {object} swaggertypes.StandardUnauthorisedError
// @Failure 500 {object} swaggertypes.StandardInternalServerError
// @Router /countries/{id} [get]
func (c *countryController) Find(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiErr := resterror.NewBadRequestError("INVALID_COUNTRY_ID")
		ctx.JSON(apiErr.Code(), apiErr)
		return
	}
	result, apiErr := services.CountryService.Find(id, 0)
	if err != nil {
		ctx.JSON(apiErr.Code(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, swaggertypes.NoErrorData{
		Data: result,
		Code: http.StatusOK,
	})
}

// List
// @Summary List countries
// @Description Retrieve all countries
// @ID v1-countries-list
// @Produce json
// @Tags Countries
// @Param code query string false "filter by code"
// @Param name query string false "filter by name"
// @Param active query bool false "filter by status" Enums(true,false)
// @Param order query string false "order direction" Enums(asc,desc)
// @Param order_by query string false "order field" Enums(id,code,name,active)
// @Param page query integer false "page number"
// @Param per_page query integer false "records per page"
// @Success 200 {object} swaggertypes.PaginatedData{data=pagination.PaginatedResponse{data=[]countries.CountryOutput}}
// @Failure 400 {object} swaggertypes.StandardBadRequestError
// @Failure 401 {object} swaggertypes.StandardUnauthorisedError
// @Failure 500 {object} swaggertypes.StandardInternalServerError
// @Router /countries [get]
func (c *countryController) List(ctx *gin.Context) {
	var req countries.ListCountryInput

	if ok := utils.GinShouldPassAll(ctx,
		utils.GinShouldBind(&req),
		utils.GinShouldValidate(&req),
	); !ok {
		return
	}

	results, apiErr := services.CountryService.List(&req)
	if apiErr != nil {
		ctx.JSON(apiErr.Code(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, swaggertypes.NoErrorData{
		Data: results,
		Code: http.StatusOK,
	})
}

// Delete
// @Summary Delete country
// @Description Endpoint used to delete an existing country record
// @ID v1-countries-delete
// @Produce json
// @Accept json
// @Tags Countries
// @Param id path int true "Country ID"
// @Success 200 {object} swaggertypes.NoErrorString
// @Failure 400 {object} swaggertypes.StandardBadRequestError
// @Failure 401 {object} swaggertypes.StandardUnauthorisedError
// @Failure 500 {object} swaggertypes.StandardInternalServerError
// @Router /countries/{id} [delete]
func (c *countryController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiErr := resterror.NewBadRequestError("INVALID_COUNTRY_ID")
		ctx.JSON(apiErr.Code(), apiErr)
		return
	}

	if err := services.CountryService.Delete(id); err != nil {
		ctx.JSON(err.Code(), err)
		return
	}

	ctx.JSON(http.StatusOK, swaggertypes.NoErrorString{
		Message: "SUCCESS",
		Code:    http.StatusOK,
	})
}

func (c *countryController) SyncCountries(ctx *gin.Context) {

}
