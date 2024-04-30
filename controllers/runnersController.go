package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runners/models"
	"runners/services"

	"github.com/gin-gonic/gin"
)

const (
	ROLE_ADMIN  = "admin"
	ROLE_RUNNER = "runner"
)

type RunnersContoller struct {
	runnersService *services.RunnersService
	usersService   *services.UsersService
}

func NewRunnersController(runnersService *services.RunnersService, usersService *services.UsersService) *RunnersContoller {
	return &RunnersContoller{runnersService: runnersService, usersService: usersService}
}

func (rh RunnersContoller) CreateRunner(ctx *gin.Context) {

	// Authentication and Authorization
	accessToken := ctx.Request.Header.Get("Token")
	auth, resErr := rh.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN})
	if resErr != nil {
		ctx.JSON(resErr.Status, resErr)
		return
	}
	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var runner models.Runner
	err = json.Unmarshal(body, &runner)
	if err != nil {
		log.Println("Error while unmarshalling create runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response, responseErr := rh.runnersService.CreateRunner(&runner)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (rh RunnersContoller) UpdateRunner(ctx *gin.Context) {

	// Authentication and Authorization
	accessToken := ctx.Request.Header.Get("Token")
	auth, resErr := rh.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN})
	if resErr != nil {
		ctx.JSON(resErr.Status, resErr)
		return
	}
	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading update runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var runner models.Runner
	err = json.Unmarshal(body, &runner)
	if err != nil {
		log.Println("update runner request body update runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	responseErr := rh.runnersService.UpdateRunner(&runner)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (rh RunnersContoller) DeleteRunner(ctx *gin.Context) {

	// Authentication and Authorization
	accessToken := ctx.Request.Header.Get("Token")
	auth, resErr := rh.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN})
	if resErr != nil {
		ctx.JSON(resErr.Status, resErr)
		return
	}
	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	runnerId := ctx.Param("id")
	responseErr := rh.runnersService.DeleteRunner(runnerId)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (rh RunnersContoller) GetRunner(ctx *gin.Context) {

	// Authentication and Authorization
	accessToken := ctx.Request.Header.Get("Token")
	auth, resErr := rh.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN, ROLE_RUNNER})
	if resErr != nil {
		ctx.JSON(resErr.Status, resErr)
		return
	}
	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	runnerId := ctx.Param("id")
	response, responseErr := rh.runnersService.GetRunner(runnerId)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (rh RunnersContoller) GetRunnersBatch(ctx *gin.Context) {

	// Authentication and Authorization
	accessToken := ctx.Request.Header.Get("Token")
	auth, resErr := rh.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN, ROLE_RUNNER})
	if resErr != nil {
		ctx.JSON(resErr.Status, resErr)
		return
	}
	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	params := ctx.Request.URL.Query()
	country := params.Get("country")
	year := params.Get("year")

	response, responseErr := rh.runnersService.GetRunnersBatch(country, year)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
