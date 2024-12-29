package contract

import "errors"

var (
	ErrDecodeJSON      = errors.New("decode json")
	ErrInvalidRequest  = errors.New("invalid request")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrMatchNotFound   = errors.New("match not found")
	MatchStatusInvalid = errors.New("match status invalid")
)

const (
	FailedDecodeJSON      = "Failed to decode json"
	InvalidRequest        = "Invalid request"
	Unauthorized          = "Unauthorized"
	MatchNotFound         = "Match not found"
	MatchStatusInvalidMsg = "Match status is not scheduled"
)
