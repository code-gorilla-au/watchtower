package notifications

import (
	"time"
	"watchtower/internal/database"
)

func toTime(ts int64) time.Time {
	return time.Unix(ts, 0).UTC()
}

func fromNotificationModel(model database.OrganisationNotification) Notification {
	return Notification{
		ID:             model.ID,
		OrganisationID: model.OrganisationID.Int64,
		ExternalID:     model.ExternalID,
		Status:         NotificationStatus(model.Status),
		Content:        model.Content,
		Type:           model.Type,
		CreatedAt:      toTime(model.CreatedAt),
		UpdatedAt:      toTime(model.UpdatedAt),
	}
}

func fromNotificationModels(models []database.OrganisationNotification) []Notification {
	notifications := make([]Notification, len(models))
	for i, model := range models {
		notifications[i] = fromNotificationModel(model)
	}
	return notifications
}
