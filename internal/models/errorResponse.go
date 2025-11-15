package models

import "errors"

var (
	ErrTeamExists   = errors.New("team name already exists")
	ErrPRExists     = errors.New("PR id already exists")
	ErrPRMerged     = errors.New("cannot reassign on merged PR")
	ErrNotAssigned  = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate  = errors.New("no active replacement candidate in team")
	ErrNotFound     = errors.New("resource not found")
)
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code" binding:"required"`
		Message string `json:"message" binding:"required"`
	} `json:"error"`
}