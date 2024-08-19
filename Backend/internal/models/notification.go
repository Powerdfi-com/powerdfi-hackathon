package models

import "time"

type NotificationType int8

const (
	NOTIFICATION_TYPE_SALE NotificationType = iota + 1
	NOTIFICATION_TYPE_REJECT
	NOTIFICATION_TYPE_APPROVE
	NOTIFICATION_TYPE_LOGIN
	NOTIFICATION_TYPE_CREATED
)

var notificationNameMap = map[NotificationType]string{
	NOTIFICATION_TYPE_SALE:    "sale",
	NOTIFICATION_TYPE_REJECT:  "reject",
	NOTIFICATION_TYPE_APPROVE: "approve",
	NOTIFICATION_TYPE_LOGIN:   "login",
	NOTIFICATION_TYPE_CREATED: "created",
}

type NotificationVerifiedData struct {
	AssetId   string `json:"assetId"`
	AssetName string `json:"assetName"`
}

type NotificationRejectData struct {
	NotificationVerifiedData
	Reason string `json:"reason"`
}
type NotificationSaleData struct {
	NotificationVerifiedData
	Price float64 `json:"price"`
}

type Notification struct {
	Id        string
	Type      NotificationType
	UserId    string
	Viewed    bool
	CreatedAt time.Time
	Data      string
}

func (n Notification) GetType() string {
	return notificationNameMap[n.Type]
}

type AdminNotification struct {
	Id        string
	Type      NotificationType
	AdminId   string
	Viewed    bool
	CreatedAt time.Time
	Data      string
}

func (n AdminNotification) GetType() string {
	return notificationNameMap[n.Type]
}

type NotificationPrefs struct {
	Sale     bool
	Verified bool
	Rejected bool
	Login    bool
}
type AdminNotificationPrefs struct {
	Created bool
	Login   bool
}
