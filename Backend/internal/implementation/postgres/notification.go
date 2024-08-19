package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type notificationImplementation struct {
	Db *sql.DB
}

func NewNotificationImplementation(db *sql.DB) repository.NotificationRepository {
	return notificationImplementation{Db: db}
}

func (n notificationImplementation) Create(notification models.Notification) (models.Notification, error) {
	// check if notification type is allowed by user
	stmt := `
	SELECT
		CASE
			WHEN 
			    -- check if the notification type is on the user's whitelist 
			    $2 = ANY(whitelist) 
			THEN TRUE
			ELSE FALSE
		END
	FROM
		notification_prefs
	WHERE 
	    user_id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var status bool
	err := n.Db.QueryRowContext(
		ctx,
		stmt,
		notification.UserId,
		notification.Type,
	).Scan(&status)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			break

		default:
			return models.Notification{}, err
		}
	}

	// return early if notification type isn't set
	if !status {
		return models.Notification{}, repository.ErrDisabledNotification
	}
	// move on to create the notification
	if notification.Id == "" {
		notification.Id = uuid.NewString()
	}

	stmt = `
	INSERT INTO user_notifications
	    (
	     id, 
	     type, 
	     user_id,
	     data,
	     viewed
	     )
	VALUES ($1, $2, $3, $4, $5)
	
	RETURNING 
	    id, 
	    created_at
	`

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	createdNotification := notification
	err = n.Db.QueryRowContext(
		ctx,
		stmt,
		notification.Id,
		notification.Type,
		notification.UserId,
		notification.Data,
		notification.Viewed,
	).Scan(
		&createdNotification.Id,
		&createdNotification.CreatedAt,
	)

	if err != nil {
		return models.Notification{}, err
	}

	return createdNotification, nil
}

// GetForUser retrieves user unseen notifications and sets them as seen.
func (n notificationImplementation) GetForUser(userId string, filter models.Filter) ([]models.Notification, int, error) {
	// statement to retrieve user notifications
	var totalCount int
	stmt := `
	SELECT
	    id, 
	    type,
	    user_id,
	    data,
	    viewed,
	    created_at
	FROM user_notifications
	WHERE user_id = $1
		LIMIT $2 OFFSET $3;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	notifications := []models.Notification{}

	rows, err := n.Db.QueryContext(ctx, stmt, userId, filter.Limit, filter.Offset())
	if err != nil {
		return []models.Notification{}, totalCount, err
	}
	defer rows.Close()

	totalCountQuery := `
	SELECT COUNT(*)
   FROM user_notifications
	WHERE user_id = $1
 `

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = n.Db.QueryRowContext(ctx, totalCountQuery, userId).Scan(&totalCount)
	if err != nil {
		return nil, totalCount, err
	}
	for rows.Next() {
		notification := models.Notification{}

		rows.Scan(
			&notification.Id,
			&notification.Type,
			&notification.UserId,
			&notification.Data,
			&notification.Viewed,

			&notification.CreatedAt,
		)

		notifications = append(notifications, notification)
	}

	if rows.Err() != nil {
		return []models.Notification{}, totalCount, err
	}

	// set fetched notifications as seen

	return notifications, totalCount, nil
}

func (n notificationImplementation) MarkAllAsRead(userId string) error {
	stmt := `
	UPDATE user_notifications
	SET viewed = true
	WHERE user_id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := n.Db.ExecContext(ctx, stmt, userId)
	if err != nil {
		return err
	}
	return nil
}
func (n notificationImplementation) MarkAsRead(userId, notificationId string) error {
	stmt := `
	UPDATE user_notifications
	SET viewed = true
	WHERE id=$1 AND user_id = $2;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := n.Db.ExecContext(ctx, stmt, notificationId, userId)
	if err != nil {
		return err
	}
	return nil
}

// CountForUser returns the number of unseen notifications for a user.
func (n notificationImplementation) CountForUser(userId string) (int64, error) {
	stmt := `
	SELECT COUNT(*)
	FROM user_notifications
	WHERE user_id = $1 AND viewed = FALSE
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var count int64
	err := n.Db.QueryRowContext(ctx, stmt, userId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (n notificationImplementation) GetUserPrefs(userId string) (models.NotificationPrefs, error) {
	stmt := `
	SELECT whitelist FROM notification_prefs
	WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	whitelist := []int64{}
	err := n.Db.QueryRowContext(ctx, stmt, userId).Scan(pq.Array(&whitelist))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.NotificationPrefs{}, nil

		default:
			return models.NotificationPrefs{}, err
		}
	}

	// map whitelist to notification prefs struct
	prefs := n.notificationWhitelistToPrefs(whitelist)
	return prefs, nil
}

func (n notificationImplementation) UpdateUserPrefs(userId string, prefs models.NotificationPrefs) (models.NotificationPrefs, error) {
	stmt := `
	-- create non-existent prefs for user
	INSERT INTO notification_prefs (user_id, whitelist)
	VALUES ($1, $2)
	
	-- or update an existing preference
	ON CONFLICT (user_id) 
	    DO UPDATE SET whitelist = $2

	RETURNING whitelist
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	whitelist := []int64{}
	err := n.Db.QueryRowContext(
		ctx,
		stmt,
		userId,
		pq.Array(n.notificationPrefsToWhitelist(prefs)),
	).Scan(pq.Array(&whitelist))

	if err != nil {
		return models.NotificationPrefs{}, err
	}

	newPrefs := n.notificationWhitelistToPrefs(whitelist)
	return newPrefs, nil
}

func (n notificationImplementation) notificationWhitelistToPrefs(whitelist []int64) models.NotificationPrefs {
	prefs := models.NotificationPrefs{}
	for _, notification := range whitelist {
		switch notification {
		case int64(models.NOTIFICATION_TYPE_SALE):
			prefs.Sale = true

		case int64(models.NOTIFICATION_TYPE_APPROVE):
			prefs.Verified = true

		case int64(models.NOTIFICATION_TYPE_REJECT):
			prefs.Rejected = true

		case int64(models.NOTIFICATION_TYPE_LOGIN):
			prefs.Login = true

		}
	}
	return prefs
}

func (n notificationImplementation) notificationPrefsToWhitelist(prefs models.NotificationPrefs) []int {
	whitelist := []int{}

	if prefs.Sale {
		whitelist = append(whitelist, int(models.NOTIFICATION_TYPE_SALE))
	}

	if prefs.Rejected {
		whitelist = append(whitelist, int(models.NOTIFICATION_TYPE_REJECT))
	}

	if prefs.Verified {
		whitelist = append(whitelist, int(models.NOTIFICATION_TYPE_APPROVE))
	}

	if prefs.Login {
		whitelist = append(whitelist, int(models.NOTIFICATION_TYPE_LOGIN))
	}

	return whitelist
}

func (n notificationImplementation) CreateForAdmin(notification models.AdminNotification) (models.AdminNotification, error) {
	// check if notification type is allowed by user
	stmt := `
	SELECT
		CASE
			WHEN 
			    -- check if the notification type is on the user's whitelist 
			    $2 = ANY(whitelist) 
			THEN TRUE
			ELSE FALSE
		END
	FROM
		admin_notification_prefs
	WHERE 
	    admin_id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var status bool
	err := n.Db.QueryRowContext(
		ctx,
		stmt,
		notification.AdminId,
		notification.Type,
	).Scan(&status)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			break

		default:
			return models.AdminNotification{}, err
		}
	}

	// return early if notification type isn't set
	if !status {
		return models.AdminNotification{}, repository.ErrDisabledNotification
	}
	// move on to create the notification
	if notification.Id == "" {
		notification.Id = uuid.NewString()
	}

	stmt = `
	INSERT INTO admin_notifications
	    (
	     id, 
	     type, 
	     admin_id,
	     data,
	     viewed
	     )
	VALUES ($1, $2, $3, $4, $5)
	
	RETURNING 
	    id, 
	    created_at
	`

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	createdNotification := notification
	err = n.Db.QueryRowContext(
		ctx,
		stmt,
		notification.Id,
		notification.Type,
		notification.AdminId,
		notification.Data,
		notification.Viewed,
	).Scan(
		&createdNotification.Id,
		&createdNotification.CreatedAt,
	)

	if err != nil {
		return models.AdminNotification{}, err
	}

	return createdNotification, nil
}

// GetForUser retrieves user unseen notifications and sets them as seen.
func (n notificationImplementation) GetForAdmin(adminId string, filter models.Filter) ([]models.AdminNotification, int, error) {
	// statement to retrieve user notifications
	var totalCount int
	stmt := `
	SELECT
	    id, 
	    type,
	    admin_id,
	    data,
	    viewed,
	    created_at
	FROM admin_notifications
	WHERE admin_id = $1
		LIMIT $2 OFFSET $3;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	notifications := []models.AdminNotification{}

	rows, err := n.Db.QueryContext(ctx, stmt, adminId, filter.Limit, filter.Offset())
	if err != nil {
		return []models.AdminNotification{}, totalCount, err
	}
	defer rows.Close()

	totalCountQuery := `
	SELECT COUNT(*)
   FROM admin_notifications
	WHERE admin_id = $1
 `

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = n.Db.QueryRowContext(ctx, totalCountQuery, adminId).Scan(&totalCount)
	if err != nil {
		return nil, totalCount, err
	}
	for rows.Next() {
		notification := models.AdminNotification{}

		rows.Scan(
			&notification.Id,
			&notification.Type,
			&notification.AdminId,
			&notification.Data,
			&notification.Viewed,

			&notification.CreatedAt,
		)

		notifications = append(notifications, notification)
	}

	if rows.Err() != nil {
		return []models.AdminNotification{}, totalCount, err
	}

	// set fetched notifications as seen

	return notifications, totalCount, nil
}

func (n notificationImplementation) MarkAllAsReadAdmin(adminId string) error {
	stmt := `
	UPDATE admin_notifications
	SET viewed = true
	WHERE admin_id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := n.Db.ExecContext(ctx, stmt, adminId)
	if err != nil {
		return err
	}
	return nil
}
func (n notificationImplementation) CountForAdmin(adminId string) (int64, error) {
	stmt := `
	SELECT COUNT(*)
	FROM admin_notifications
	WHERE admin_id = $1 AND viewed = FALSE
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var count int64
	err := n.Db.QueryRowContext(ctx, stmt, adminId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (n notificationImplementation) GetAdminUserPrefs(userId string) (models.AdminNotificationPrefs, error) {
	stmt := `
	SELECT whitelist FROM admin_notification_prefs
	WHERE admin_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	whitelist := []int64{}
	err := n.Db.QueryRowContext(ctx, stmt, userId).Scan(pq.Array(&whitelist))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.AdminNotificationPrefs{}, nil

		default:
			return models.AdminNotificationPrefs{}, err
		}
	}

	// map whitelist to notification prefs struct
	prefs := n.adminNotificationWhitelistToPrefs(whitelist)
	return prefs, nil
}

func (n notificationImplementation) UpdateAdminPrefs(userId string, prefs models.AdminNotificationPrefs) (models.AdminNotificationPrefs, error) {
	stmt := `
	-- create non-existent prefs for user
	INSERT INTO admin_notification_prefs (admin_id, whitelist)
	VALUES ($1, $2)
	
	-- or update an existing preference
	ON CONFLICT (admin_id) 
	    DO UPDATE SET whitelist = $2

	RETURNING whitelist
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	whitelist := []int64{}
	err := n.Db.QueryRowContext(
		ctx,
		stmt,
		userId,
		pq.Array(n.adminNotificationPrefsToWhitelist(prefs)),
	).Scan(pq.Array(&whitelist))

	if err != nil {
		return models.AdminNotificationPrefs{}, err
	}

	newPrefs := n.adminNotificationWhitelistToPrefs(whitelist)
	return newPrefs, nil
}

func (n notificationImplementation) adminNotificationWhitelistToPrefs(whitelist []int64) models.AdminNotificationPrefs {
	prefs := models.AdminNotificationPrefs{}
	for _, notification := range whitelist {
		switch notification {
		case int64(models.NOTIFICATION_TYPE_CREATED):
			prefs.Created = true

		case int64(models.NOTIFICATION_TYPE_LOGIN):
			prefs.Login = true

		}
	}
	return prefs
}

func (n notificationImplementation) adminNotificationPrefsToWhitelist(prefs models.AdminNotificationPrefs) []int {
	whitelist := []int{}

	if prefs.Created {
		whitelist = append(whitelist, int(models.NOTIFICATION_TYPE_CREATED))
	}

	if prefs.Login {
		whitelist = append(whitelist, int(models.NOTIFICATION_TYPE_LOGIN))
	}

	return whitelist
}
