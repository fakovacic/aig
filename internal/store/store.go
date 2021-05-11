package store

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/fakovacic/aig/internal/meta"
	esConfig "github.com/fakovacic/aig/internal/store/elasticsearch/config"
	elastic "github.com/olivere/elastic/v7"
)

// Model service model
type Model interface {
	Table() string
	SearchField() string
	Mapping() (string, error)
	Clear() Model
	Service() string
	GetID() string
	ModelReader
	ModelCreator
	ModelUpdater
	ModelDeleter
}

// ModelReader service model
type ModelReader interface {
	BRead() error
	ARead() error
}

// ModelCreator service model
type ModelCreator interface {
	BCreate() error
	ACreate() error
}

// ModelUpdater service model
type ModelUpdater interface {
	BUpdate() error
	AUpdate() error
}

// ModelDeleter service model
type ModelDeleter interface {
	BDelete() error
	ADelete() error
}

// ModelStore models store struct
type ModelStore interface {
	GetModel() Model
	Init(context.Context)
	Create(context.Context, Model) (Model, error)
	ReadBy(context.Context, string, interface{}) (Model, error)
	Update(context.Context, Model, []string) (Model, error)
	Delete(context.Context, Model) error
	List(context.Context, *meta.Model) ([]json.RawMessage, error)
	ListByID(context.Context, ...IDf) (map[string]map[string]json.RawMessage, error)
	ListIDs(context.Context, ...IDs) ([]IDf, error)
}

// Store databases store
type Store struct {
	sync.Mutex
	es *elastic.Client
	SS *ServicesStore
}

// New return new store
func New() *Store {
	return &Store{
		SS: NewServicesStore(),
	}
}

// ES return elasticsearch client
func (s *Store) ES() *elastic.Client {
	return s.es
}

// Register register store database type
func (s *Store) Register(t ...DBType) *Store {
	s.Lock()
	defer s.Unlock()

	for i := range t {
		switch t[i] {
		case ElasticSearch:
			err := s.newElasticSearch()
			if err != nil {
				fmt.Printf("error elasticsearch:%v\n", err)
			}
		}
	}

	return s
}

// newElasticSearch return elasticsearch store session
func (s *Store) newElasticSearch() error {
	var err error

	s.es, err = esConfig.RefreshSession(s.es)
	if err != nil {
		return err
	}
	return nil
}

// ServicesStore database types for services
type ServicesStore struct {
	sync.Mutex
	Services map[string]DBType
}

// NewServicesStore create new service store
func NewServicesStore() *ServicesStore {
	return &ServicesStore{
		Services: make(map[string]DBType),
	}
}

// Register to map
func (s *ServicesStore) Register(k Model, v DBType) *ServicesStore {
	s.Lock()
	defer s.Unlock()
	s.Services[k.Service()] = v
	return s
}

// Get return service if exist
func (s *ServicesStore) Get(k Model) DBType {
	s.Lock()
	defer s.Unlock()
	val, ok := s.Services[k.Service()]
	if ok {
		return val
	}
	return ""
}
