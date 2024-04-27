package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runners/models"

	"github.com/gin-gonic/gin"
)

type ResultsController struct {
	resultsService *services.ResultsService
}

func (rh ResultsController) CreateResult(ctx *gin.Context) {
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
	resultId := ctx.Param("id")
	responseErr := rh.resultsService.DeleteResult(resultId)
	if responseErr != nil {
		log.Println("", responseErr)
		ctx.JSON(responseErr.Status, responseErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}
