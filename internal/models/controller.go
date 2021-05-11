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
	// ControllerTable const
	ControllerTable string = "controllers"
	// ControllerService const
	ControllerService string = "controller"
)

// Controller model struct
type Controller struct {
	ID          string    `json:"id"`
	IP          string    `json:"ip"`
	Humidity    float64   `json:"humidity"`
	Temperature float64   `json:"temperature"`
	HeatIndex   float64   `json:"heatIndex"`
	Soil        float64   `json:"soil"`
	Water       bool      `json:"water"`
	Lights      bool      `json:"lights"`
	Vent        bool      `json:"vent"`
	Online      bool      `json:"online"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

// ControllerField fields
var ControllerField = struct {
	ID          string
	IP          string
	Humidity    string
	Temperature string
	HeatIndex   string
	Soil        string
	Water       string
	Lights      string
	Vent        string
	Online      string
	Created     string
	Modified    string
}{
	ID:          "id",
	IP:          "ip",
	Humidity:    "humidity",
	Temperature: "temperature",
	HeatIndex:   "heatIndex",
	Soil:        "soil",
	Water:       "water",
	Lights:      "lights",
	Vent:        "vent",
	Online:      "online",
	Created:     "created",
	Modified:    "modified",
}

// ControllerAllFields contain all fields
var ControllerAllFields = []string{
	ControllerField.IP,
	ControllerField.Humidity,
	ControllerField.Temperature,
	ControllerField.HeatIndex,
	ControllerField.Soil,
	ControllerField.Water,
	ControllerField.Lights,
	ControllerField.Vent,
}

// ControllerSortFields sort fields
var ControllerSortFields = map[string]string{
	ControllerField.Created: "Created",
}

// ControllerFilterFields filter fields
var ControllerFilterFields = map[string]string{
	ControllerField.ID: "ID",
}

// Table implement model table
func (d *Controller) Table() string {
	return ControllerTable
}

// SearchField implement model table
func (d *Controller) SearchField() string {
	return ControllerField.ID
}

// Mapping table mapping
func (d *Controller) Mapping() (string, error) {
	return mapping.Map(Controller{}).JSON()
}

// GetID implement model get id value
func (d *Controller) GetID() string {
	return d.ID
}

// Service implement model service
func (d *Controller) Service() string {
	return ControllerService
}

// Clear implement clear
func (d *Controller) Clear() store.Model {
	return &Controller{}
}

// ClearRel clear relationships
func (d *Controller) ClearRel() {
}

// KeyVal return value for key val search
func (d *Controller) KeyVal() string {
	return d.ID
}

// BRead implement model before read method
func (d *Controller) BRead() error {
	return nil
}

// ARead implement model after read method
func (d *Controller) ARead() error {
	return nil
}

// BCreate implement model before creating method
func (d *Controller) BCreate() error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	d.Created = time.Now()
	return nil
}

// ACreate implement model after creating method
func (d *Controller) ACreate() error {
	return nil
}

// BUpdate implement model before update method
func (d *Controller) BUpdate() error {
	return nil
}

// AUpdate implement model after update method
func (d *Controller) AUpdate() error {
	return nil
}

// BDelete implement model before delete method
func (d *Controller) BDelete() error {
	return nil
}

// ADelete implement model after delete method
func (d *Controller) ADelete() error {
	return nil
}

// ToStats adapter to stats
func (d *Controller) ToStats() *Stats {
	return &Stats{
		Controller:  d.ID,
		IP:          d.IP,
		Humidity:    d.Humidity,
		Temperature: d.Temperature,
		HeatIndex:   d.HeatIndex,
		Soil:        d.Soil,
		Water:       d.Water,
		Lights:      d.Lights,
		Vent:        d.Vent,
	}
}

// RawController unmarshal raw message
func RawController(raw json.RawMessage) (*Controller, error) {
	m := &Controller{}
	if err := json.Unmarshal(raw, m); err != nil {
		return nil, errors.InternalWrap(err, "unmarshal document")
	}
	return m, nil
}
