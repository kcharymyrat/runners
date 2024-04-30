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

type ResultsController struct {
	resultsService *services.ResultsService
	usersService   *services.UsersService
}

func NewResultsController(resultsService *services.ResultsService, usersService *services.UsersService) *ResultsController {
	return &ResultsController{resultsService: resultsService, usersService: usersService}
}

func (rh ResultsController) CreateResult(ctx *gin.Context) {

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
		log.Println("Error while reading create result request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var result models.Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("create result request body creates result request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response, responseErr := rh.resultsService.CreateResult(&result)
	if responseErr != nil {
		log.Println("err", err)
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (rh ResultsController) DeleteResult(ctx *gin.Context) {

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

	resultId := ctx.Param("id")
	responseErr := rh.resultsService.DeleteResult(resultId)
	if responseErr != nil {
		log.Println("", responseErr)
		ctx.JSON(responseErr.Status, responseErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}
