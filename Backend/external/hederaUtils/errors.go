package hederaUtils

import "errors"

var (
	ErrIncompleteTreasuryDetails = errors.New("incomplete treasury details")
	ErrInvalidTokenType          = errors.New("invalid token type")
	ErrInvalidTokenID            = errors.New("invalid treasury token ID")
)
