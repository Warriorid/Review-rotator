package models

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code" binding:"required"`
		Message string `json:"message" binding:"required"`
	} `json:"error"`
}