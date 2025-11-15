package handler

import (
	"net/http"
	"review-rotator/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) setUserActive(c *gin.Context) {
    adminToken := c.GetHeader("X-Admin-Token")
    if adminToken == "" {
        errorResponse := models.ErrorResponse{}
        errorResponse.Error.Code = "UNAUTHORIZED"
        errorResponse.Error.Message = "No/invalid admin token"
        c.JSON(http.StatusUnauthorized, errorResponse)
        return
    }

    var input models.SetUserActiveRequest
	if err := c.BindJSON(&input); err != nil {
		errorResponse := models.ErrorResponse{}
		errorResponse.Error.Code = "INVALID_REQUEST"
		errorResponse.Error.Message = "Invalid request body"
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

    user, err := h.service.User.SetUserActive(input.UserID, input.IsActive)
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

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *Handler) getUserReviews(c *gin.Context) {
    userID := c.Query("user_id")
    if userID == "" {
        errorResponse := models.ErrorResponse{}
        errorResponse.Error.Code = "INVALID_REQUEST"
        errorResponse.Error.Message = "user_id parameter is required"
        c.JSON(http.StatusBadRequest, errorResponse)
        return
    }

    result, err := h.service.User.GetUserReviewRequests(userID)
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

    c.JSON(http.StatusOK, result)
}