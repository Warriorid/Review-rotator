package handler

import (
	"net/http"
	"review-rotator/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) addTeam(c *gin.Context) {
	var input models.Team
	if err := c.BindJSON(&input); err != nil {
		errorResponse := models.ErrorResponse{}
		errorResponse.Error.Code = "INVALID_REQUEST"
		errorResponse.Error.Message = "Invalid request body"
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	team, err := h.service.Team.AddTeam(input)
	if err != nil {
		errorResponse := models.ErrorResponse{}
		errorCode := "INTERNAL_ERROR"
		statusCode := http.StatusInternalServerError

		switch err {
		case models.ErrTeamExists:
			errorCode = "TEAM_EXISTS"
			statusCode = http.StatusBadRequest
		}

		errorResponse.Error.Code = errorCode
		errorResponse.Error.Message = "team_name already exists"
		c.JSON(statusCode, errorResponse)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"team": team})
}

func (h *Handler) getTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		errorResponse := models.ErrorResponse{}
		errorResponse.Error.Code = "INVALID_REQUEST"
		errorResponse.Error.Message = "team_name parameter is required"
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	team, err := h.service.Team.GetTeam(teamName)
	if err != nil {
		errorResponse := models.ErrorResponse{}
		errorCode := "INTERNAL_ERROR"
		statusCode := http.StatusInternalServerError

		switch err {
		case models.ErrNotFound:
			errorCode = "NOT_FOUND"
			statusCode = http.StatusNotFound
		}

		errorResponse.Error.Code = errorCode
		errorResponse.Error.Message = err.Error()
		c.JSON(statusCode, errorResponse)
		return
	}

	c.JSON(http.StatusOK, team)
}