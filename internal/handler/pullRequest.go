package handler

import (
	"net/http"
	"review-rotator/internal/models"
	"github.com/gin-gonic/gin"
)


func (h *Handler) createPullRequest(c *gin.Context) {
    var input models.CreatePRRequest
    if err := c.BindJSON(&input); err != nil {
        errorResponse := models.ErrorResponse{}
        errorResponse.Error.Code = "INVALID_REQUEST"
        errorResponse.Error.Message = "Invalid request body"
        c.JSON(http.StatusBadRequest, errorResponse)
        return
    }

    pr, err := h.service.PullRequest.CreatePullRequest(input)
    if err != nil {
        errorResponse := models.ErrorResponse{}
        errorCode := "INTERNAL_ERROR"
        statusCode := http.StatusInternalServerError

        switch err {
        case models.ErrPRExists:
            errorCode = "PR_EXISTS"
            statusCode = http.StatusConflict
        case models.ErrNotFound:
            errorCode = "NOT_FOUND"
            statusCode = http.StatusNotFound
        }

        errorResponse.Error.Code = errorCode
        errorResponse.Error.Message = err.Error()
        c.JSON(statusCode, errorResponse)
        return
    }

    c.JSON(http.StatusCreated, gin.H{"pr": pr})
}

func (h *Handler) mergePullRequest(c *gin.Context) {
    var input models.MergePRRequest
    if err := c.BindJSON(&input); err != nil {
        errorResponse := models.ErrorResponse{}
        errorResponse.Error.Code = "INVALID_REQUEST"
        errorResponse.Error.Message = "Invalid request body"
        c.JSON(http.StatusBadRequest, errorResponse)
        return
    }

    pr, err := h.service.PullRequest.MergePullRequest(input.PullRequestID)
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

    c.JSON(http.StatusOK, gin.H{"pr": pr})
}

func (h *Handler) reassignReviewer(c *gin.Context) {
    var input models.ReassignReviewerRequest
    if err := c.BindJSON(&input); err != nil {
        errorResponse := models.ErrorResponse{}
        errorResponse.Error.Code = "INVALID_REQUEST"
        errorResponse.Error.Message = "Invalid request body"
        c.JSON(http.StatusBadRequest, errorResponse)
        return
    }

    result, err := h.service.PullRequest.ReassignReviewer(input)
    if err != nil {
        errorResponse := models.ErrorResponse{}
        errorCode := "INTERNAL_ERROR"
        statusCode := http.StatusInternalServerError

        switch err {
        case models.ErrNotFound:
            errorCode = "NOT_FOUND"
            statusCode = http.StatusNotFound
        case models.ErrPRMerged:
            errorCode = "PR_MERGED"
            statusCode = http.StatusConflict
        case models.ErrNotAssigned:
            errorCode = "NOT_ASSIGNED"
            statusCode = http.StatusConflict
        case models.ErrNoCandidate:
            errorCode = "NO_CANDIDATE"
            statusCode = http.StatusConflict
        }

        errorResponse.Error.Code = errorCode
        errorResponse.Error.Message = err.Error()
        c.JSON(statusCode, errorResponse)
        return
    }

    c.JSON(http.StatusOK, result)
}