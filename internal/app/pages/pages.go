package pages

import (
	"github.com/fakovacic/aig/pkg/notifications"
)

// Page define web pages base struct
type Page struct {
	Title         string
	Notifications *notifications.Notifications
}
