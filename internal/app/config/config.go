package config

import (
	"github.com/fakovacic/aig/internal/config"
	"github.com/fakovacic/aig/pkg/notifications"
)

// Config admin interface
type Config interface {
	config.Config
	ClearNotifications()
	AddNotification(notificationType, format string, args ...interface{})
	GetNotifications() *notifications.Notifications
}
