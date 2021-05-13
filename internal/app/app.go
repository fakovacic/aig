package app

import (
	"context"

	rr "github.com/fakovacic/aig/internal/app/routes"
	"github.com/fakovacic/aig/internal/app/tasks"
	"github.com/fakovacic/aig/internal/models"
	"github.com/fakovacic/aig/internal/routes"
	"github.com/fakovacic/aig/internal/store"
	"github.com/fakovacic/aig/internal/store/create"
	"github.com/fakovacic/aig/pkg/notifications"
	"go.elastic.co/apm/module/apmhttprouter"
)

// New create new admin config
func New() *Config {
	return &Config{
		Storer: store.New().Register(
			store.ElasticSearch,
		),
		Notifications: &notifications.Notifications{},
	}
}

// Config web config
type Config struct {
	Notifications *notifications.Notifications
	Storer        *store.Store
}

// Init config
func (c *Config) Init() {
	c.Storer.SS.
		Register(&models.Controller{}, store.ElasticSearch).
		Register(&models.Stats{}, store.ElasticSearch).
		Register(&models.Activity{}, store.ElasticSearch)

	create.Tables(c.Storer,
		&models.Controller{},
		&models.Stats{},
		&models.Activity{},
	)

	c.BackgroundTasks()
}

// HandleError handle error
func (c *Config) HandleError(ctx context.Context, err error) {
	c.Notifications.Error(err.Error())
}

// Store return store
func (c *Config) Store() *store.Store {
	return c.Storer
}

// AddNotification add notification
func (c *Config) AddNotification(notificationType, format string, args ...interface{}) {
	switch notificationType {
	case "error":
		c.Notifications.Errorf(format, args...)
	case "success":
		c.Notifications.Successf(format, args...)
	case "info":
		c.Notifications.Infof(format, args...)
	}
}

// GetNotifications return notifications
func (c *Config) GetNotifications() *notifications.Notifications {
	return c.Notifications
}

// ClearNotifications remove notifications
func (c *Config) ClearNotifications() {
	c.Notifications.Clear()
}

// Routes register routes
func (c *Config) Routes(router *apmhttprouter.Router) {
	var pr routes.RoutesMap
	pr = append(pr,
		rr.Routes(c),
	)

	routes.Register(router, pr)
}

// BackgroundTasks start tasks
func (c *Config) BackgroundTasks() {
	tasks.Ping(c)
	tasks.Config(c)
}
