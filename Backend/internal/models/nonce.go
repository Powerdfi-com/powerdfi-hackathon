package models

const (
	NonceActionCreate            = "create"
	NonceActionUpdate            = "update"
	NonceInactiveAddressDuration = "15 minutes"
)

type Nonce struct {
	// Message represents an unsigned message
	Message string

	// Action can be either "create" or "update". It signifies if a new user should be created
	// or if an existing user should have their nonce updated.
	Action string
}
