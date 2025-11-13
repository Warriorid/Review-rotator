package models

import "errors"

var (
	ErrTeamExists   = errors.New("TEAM_EXISTS")
	ErrPRExists     = errors.New("PR_EXISTS")
	ErrPRMerged     = errors.New("PR_MERGED")
	ErrNotAssigned  = errors.New("NOT_ASSIGNED")
	ErrNoCandidate  = errors.New("NO_CANDIDATE")
	ErrNotFound     = errors.New("NOT_FOUND")
)
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code" binding:"required"`
		Message string `json:"message" binding:"required"`
	} `json:"error"`
}