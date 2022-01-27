package controllers

import (
	"github.com/development-raul/footy-predictor/src/domains/seasons"
	"github.com/development-raul/footy-predictor/src/services"
	"github.com/development-raul/footy-predictor/src/swaggertypes"
	"github.com/development-raul/footy-predictor/src/utils"
	"github.com/development-raul/footy-predictor/src/utils/resterror"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type seasonControllerInterface interface {
	Create(ctx *gin.Context)
	Find(ctx *gin.Context)
	List(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Sync(ctx *gin.Context)
}

type seasonController struct{}

var SeasonController seasonControllerInterface = &seasonController{}

// Create
// @Summary Create season
// @Description Endpoint used to create a new season record
// @ID v1-seasons-create
// @Produce json
// @Accept json
// @Tags Seasons
// @Param JSON request body seasons.Season true "Request Sample"
// @Success 201 {object} swaggertypes.NoErrorString
// @Failure 400 {object} swaggertypes.StandardBadRequestError
// @Failure 401 {object} swaggertypes.StandardUnauthorisedError
// @Failure 500 {object} swaggertypes.StandardInternalServerError
// @Router /seasons [post]
func (c *seasonController) Create(ctx *gin.Context) {
	var req seasons.Season
	if ok := utils.GinShouldPassAll(ctx, utils.GinShouldBind(&req), utils.GinShouldValidate(&req)); !ok {
		return
	}

	if err := services.SeasonService.Create(req.ID); err != nil {
		ctx.JSON(err.Code(), err)
		return
	}

	ctx.JSON(http.StatusCreated, swaggertypes.NoErrorString{
		Message: "SUCCESS",
		Code:    http.StatusCreated,
	})
}

// Find
// @Summary Find season
// @Description Retrieve a season identified by id
// @ID v1-seasons-find
// @Produce json
// @Tags Seasons
// @Param id path int true "Season ID"
// @Success 200 {object} swaggertypes.NoErrorI{data=seasons.Season}
// @Failure 400 {object} swaggertypes.StandardBadRequestError
// @Failure 401 {object} swaggertypes.StandardUnauthorisedError
// @Failure 500 {object} swaggertypes.StandardInternalServerError
// @Router /seasons/{id} [get]
func (c *seasonController) Find(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiErr := resterror.NewBadRequestError("INVALID_SEASON_ID")
		ctx.JSON(apiErr.Code(), apiErr)
		return
	}
	result, apiErr := services.SeasonService.Find(id)
	if apiErr != nil {
		ctx.JSON(apiErr.Code(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, swaggertypes.NoErrorData{
		Data: result,
		Code: http.StatusOK,
	})
}

// List
// @Summary List seasons
// @Description Retrieve all seasons
// @ID v1-seasons-list
// @Produce json
// @Tags Seasons
// @Param id query string false "filter by id"
// @Param order query string false "order direction" Enums(asc,desc)
// @Success 200 {object} swaggertypes.NoErrorI{data=[]seasons.Season}
// @Failure 400 {object} swaggertypes.StandardBadRequestError
// @Failure 401 {object} swaggertypes.StandardUnauthorisedError
// @Failure 500 {object} swaggertypes.StandardInternalServerError
// @Router /seasons [get]
func (c *seasonController) List(ctx *gin.Context) {
	var req seasons.ListSeasonInput

	if ok := utils.GinShouldPassAll(ctx,
		utils.GinShouldBind(&req),
		utils.GinShouldValidate(&req),
	); !ok {
		return
	}

	results, apiErr := services.SeasonService.List(&req)
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
// @Summary Delete season
// @Description Endpoint used to delete an existing season record
// @ID v1-seasons-delete
// @Produce json
// @Accept json
// @Tags Seasons
// @Param id path int true "Season ID"
// @Success 200 {object} swaggertypes.NoErrorString
// @Failure 400 {object} swaggertypes.StandardBadRequestError
// @Failure 401 {object} swaggertypes.StandardUnauthorisedError
// @Failure 500 {object} swaggertypes.StandardInternalServerError
// @Router /seasons/{id} [delete]
func (c *seasonController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiErr := resterror.NewBadRequestError("INVALID_SEASON_ID")
		ctx.JSON(apiErr.Code(), apiErr)
		return
	}

	if err := services.SeasonService.Delete(id); err != nil {
		ctx.JSON(err.Code(), err)
		return
	}

	ctx.JSON(http.StatusOK, swaggertypes.NoErrorString{
		Message: "SUCCESS",
		Code:    http.StatusOK,
	})
}

func (c *seasonController) Sync(ctx *gin.Context) {
	if err := services.SeasonService.Sync(); err != nil {
		ctx.JSON(err.Code(), err)
		return
	}
	ctx.JSON(http.StatusOK, swaggertypes.NoErrorString{
		Message: "SUCCESS",
		Code:    http.StatusOK,
	})
}
