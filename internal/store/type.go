package store

import (
	"strings"

	"github.com/fakovacic/aig/internal/errors"
)

// DBType database type
type DBType string

const (
	// ElasticSearch const
	ElasticSearch DBType = "ELASTICSEARCH"
)

// Parse source from text
func (t *DBType) Parse(s string) error {
	s = strings.Trim(s, "\"")
	switch s {
	case "ELASTICSEARCH":
		*t = ElasticSearch
	case "":
	default:
		return errors.BadRequest("invalid db type '%s'", s)
	}

	return nil
}

func (t DBType) String() string {
	return string(t)
}

// UnmarshalJSON implements json.Marshaler
func (t *DBType) UnmarshalJSON(b []byte) error {
	return t.Parse(string(b))
}

// IDf use for list by id
type IDf struct {
	Table string
	ID    string
}

// IDs use for list to get ids
type IDs struct {
	Table string
	Field string
	Value interface{}
}
