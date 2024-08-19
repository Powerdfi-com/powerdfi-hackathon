package models

import "time"

type UserType int

const (
	USER_TYPE_USER UserType = iota
	USER_TYPE_BUISNESS
)

var roleMappings = map[UserType]string{
	USER_TYPE_USER:     "user",
	USER_TYPE_BUISNESS: "buisness",
}

type User struct {
	Id        string
	Address   string
	AccountID string

	Email               *string
	FirstName           *string
	LastName            *string
	Username            *string
	Bio                 string
	Website             string
	Twitter             string
	Discord             string
	Avatar              string
	PublicKey           *string
	EncryptedPrivateKey []byte
	UserType            UserType

	// kyc registered
	IsVerified bool
	IsActive   bool

	KYCRegisteredDate *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (u User) GetUserType() string {
	return roleMappings[u.UserType]
}
