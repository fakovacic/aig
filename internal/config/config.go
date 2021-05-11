package config

import (
	"github.com/fakovacic/aig/internal/store"
)

// Config interface
type Config interface {
	Store() *store.Store
}
