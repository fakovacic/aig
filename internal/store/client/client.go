package client

import (
	"github.com/fakovacic/aig/internal/log"
	"github.com/fakovacic/aig/internal/store"
	"github.com/fakovacic/aig/internal/store/elasticsearch"
	"github.com/fakovacic/aig/internal/store/wrapper"
)

// New create new model store
func New(s *store.Store, t store.Model) *wrapper.Store {
	db := s.SS.Get(t.Clear())
	if db == "" {
		log.Info("not found db store for:", t.Table())
		return nil
	}

	switch db {
	case store.ElasticSearch:
		return wrapper.New(
			elasticsearch.New(s, t.Clear()),
		)
	default:
		return nil
	}
}
