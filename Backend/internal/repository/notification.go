package repository

import "github.com/Powerdfi-com/Backend/internal/models"

type NotificationRepository interface {
	Create(notification models.Notification) (models.Notification, error)
	GetForUser(userId string, filter models.Filter) ([]models.Notification, int, error)
	MarkAsRead(userId, notificationId string) error
	MarkAllAsRead(userId string) error
	CountForUser(userId string) (int64, error)
	GetUserPrefs(userId string) (models.NotificationPrefs, error)
	UpdateUserPrefs(userId string, prefs models.NotificationPrefs) (models.NotificationPrefs, error)

	CreateForAdmin(notification models.AdminNotification) (models.AdminNotification, error)
	GetForAdmin(adminId string, filter models.Filter) ([]models.AdminNotification, int, error)
	MarkAllAsReadAdmin(adminId string) error
	CountForAdmin(adminId string) (int64, error)
	GetAdminUserPrefs(userId string) (models.AdminNotificationPrefs, error)
	UpdateAdminPrefs(userId string, prefs models.AdminNotificationPrefs) (models.AdminNotificationPrefs, error)
}
