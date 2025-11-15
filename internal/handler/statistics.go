package handler

import (
	"net/http"
	"review-rotator/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getStats(c *gin.Context) {
    stats, err := h.service.Statistics.GetReviewStats()
    if err != nil {
        errorResponse := models.ErrorResponse{}
        errorResponse.Error.Code = "INTERNAL_ERROR"
        errorResponse.Error.Message = err.Error()
        c.JSON(http.StatusInternalServerError, errorResponse)
        return
    }

    c.JSON(http.StatusOK, stats)
}