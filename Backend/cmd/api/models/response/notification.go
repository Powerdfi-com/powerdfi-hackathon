package response

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
)

type NotificationResponse struct {
	Id        string      `json:"id"`
	Type      string      `json:"type"`
	UserId    string      `json:"userId"`
	Data      interface{} `json:"data"`
	CreatedAt time.Time   `json:"createdAt"`
	Viewed    bool        `json:"viewed"`
}

type NotificationPrefsResponse struct {
	Sale     bool `json:"sale"`
	Verified bool `json:"verified"`
	Rejected bool `json:"rejected"`
	Login    bool `json:"login"`
}
type AdminNotificationResponse struct {
	Id        string      `json:"id"`
	Type      string      `json:"type"`
	AdminId   string      `json:"userId"`
	Data      interface{} `json:"data"`
	CreatedAt time.Time   `json:"createdAt"`
	Viewed    bool        `json:"viewed"`
}

type AdminNotificationPrefsResponse struct {
	Created bool `json:"created"`
	Login   bool `json:"login"`
}

func NotificationResponseFromModel(notification models.Notification) NotificationResponse {

	resp := NotificationResponse{
		Id:        notification.Id,
		Type:      notification.GetType(),
		UserId:    notification.UserId,
		CreatedAt: notification.CreatedAt,
	}

	_ = json.NewDecoder(strings.NewReader(notification.Data)).Decode(&resp.Data)

	resp.CreatedAt = notification.CreatedAt.UTC()
	return resp
}

func NotificationResponsePrefsFromModel(notification models.NotificationPrefs) NotificationPrefsResponse {
	return NotificationPrefsResponse{
		Sale:     notification.Sale,
		Verified: notification.Verified,
		Rejected: notification.Rejected,
		Login:    notification.Login,
	}
}

func AdminNotificationResponseFromModel(notification models.AdminNotification) AdminNotificationResponse {

	resp := AdminNotificationResponse{
		Id:        notification.Id,
		Type:      notification.GetType(),
		AdminId:   notification.AdminId,
		CreatedAt: notification.CreatedAt,
	}

	_ = json.NewDecoder(strings.NewReader(notification.Data)).Decode(&resp.Data)

	resp.CreatedAt = notification.CreatedAt.UTC()
	return resp
}

func AdminNotificationResponsePrefsFromModel(notification models.AdminNotificationPrefs) AdminNotificationPrefsResponse {
	return AdminNotificationPrefsResponse{
		Created: notification.Created,
		Login:   notification.Login,
	}
}
