package response

import (
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
)

type AdminResponse struct {
	Id        string    `json:"id"`
	Email     string    `json:"email,omitempty"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"createdAt"`
}

func AdminResponseFromModel(admin models.Admin) AdminResponse {
	response := AdminResponse{
		Id:        admin.Id,
		Email:     admin.Email,
		Roles:     admin.GetRoles(),
		CreatedAt: time.Time{},
	}
	response.CreatedAt = admin.CreatedAt.UTC()

	return response
}
