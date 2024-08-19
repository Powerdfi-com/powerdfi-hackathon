package response

import (
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
)

type WalletResponse struct {
	AccountId string  `json:"accountId"`
	Address   string  `json:"address"`
	Balance   float64 `json:"balance"`
}

type UserResponse struct {
	Id         string    `json:"id"`
	Address    string    `json:"address"`
	FirstName  *string   `json:"firstName,omitempty"`
	LastName   *string   `json:"lastName,omitempty"`
	Email      *string   `json:"email,omitempty"`
	UserName   *string   `json:"username,omitempty"`
	IsActive   bool      `json:"isActive"`
	IsVerified bool      `json:"isVerified"`
	Avatar     string    `json:"avatar"`
	Roles      []string  `json:"roles"`
	Bio        string    `json:"bio"`
	Twitter    string    `json:"twitter"`
	Discord    string    `json:"discord"`
	Website    string    `json:"website"`
	CreatedAt  time.Time `json:"createdAt"`
}

func UserResponseFromModel(user models.User) UserResponse {
	response := UserResponse{
		Id:         user.Id,
		Address:    user.Address,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Avatar:     user.Avatar,
		Website:    user.Website,
		Twitter:    user.Twitter,
		Discord:    user.Discord,
		Bio:        user.Bio,
		UserName:   user.Username,
		// Roles:      user.GetRoles(),
		CreatedAt: time.Time{},
	}
	response.CreatedAt = user.CreatedAt.UTC()

	return response
}
