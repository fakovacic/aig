package models

import (
	"encoding/json"
	"time"

	"github.com/fakovacic/aig/internal/errors"
	"github.com/fakovacic/aig/internal/store"
	"github.com/fakovacic/aig/pkg/mapping"
	"github.com/google/uuid"
)

const (
	// ActivitiesTable const
	ActivitiesTable string = "activities"
	// ActivityService const
	ActivityService string = "activity"
)

// Activity model struct
type Activity struct {
	ID        string    `json:"id" validate:"required"`
	Action    string    `json:"action"`
	ExteralID string    `json:"external_id"`
	Created   time.Time `json:"created" mapping:"date"`
}

// ActivityField user fields
var ActivityField = struct {
	ID        string
	Action    string
	ExteralID string
	Created   string
}{
	ID:        "id",
	Action:    "action",
	ExteralID: "external_id",
	Created:   "created",
}

// ActivitySortFields sort fields
var ActivitySortFields = map[string]string{
	ActivityField.Created: "Created",
}

// ActivityFilterFields filter fields
var ActivityFilterFields = map[string]string{
	ActivityField.ExteralID: "Type",
}

// Table implement model table
func (d *Activity) Table() string {
	return ActivitiesTable
}

// SearchField implement model table
func (d *Activity) SearchField() string {
	return ActivityField.ExteralID
}

// Mapping table mapping
func (d *Activity) Mapping() (string, error) {
	return mapping.Map(Activity{}).JSON()
}

// GetID implement model get id value
func (d *Activity) GetID() string {
	return d.ID
}

// Service implement model service
func (d *Activity) Service() string {
	return ActivityService
}

// Clear implement clear
func (d *Activity) Clear() store.Model {
	return &Activity{}
}

// KeyVal return value for key val search
func (d *Activity) KeyVal() string {
	return d.ExteralID
}

// BRead implement model before read method
func (d *Activity) BRead() error {
	return nil
}

// ARead implement model after read method
func (d *Activity) ARead() error {
	return nil
}

// BCreate implement model before creating method
func (d *Activity) BCreate() error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	d.Created = time.Now()
	return nil
}

// ACreate implement model after creating method
func (d *Activity) ACreate() error {
	return nil
}

// BUpdate implement model before update method
func (d *Activity) BUpdate() error {
	return nil
}

// AUpdate implement model after update method
func (d *Activity) AUpdate() error {
	return nil
}

// BDelete implement model before delete method
func (d *Activity) BDelete() error {
	return nil
}

// ADelete implement model after delete method
func (d *Activity) ADelete() error {
	return nil
}

// RawActivity unmarshal raw message
func RawActivity(raw json.RawMessage) (*Activity, error) {
	m := &Activity{}
	if err := json.Unmarshal(raw, m); err != nil {
		return nil, errors.InternalWrap(err, "unmarshal document")
	}
	return m, nil
}
