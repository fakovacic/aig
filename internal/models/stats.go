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
	// StatsTable const
	StatsTable string = "stats"
	// StatsService const
	StatsService string = "stats"
)

// Stats model struct
type Stats struct {
	ID          string    `json:"id"`
	Controller  string    `json:"controller"`
	IP          string    `json:"ip"`
	Humidity    float64   `json:"humidity"`
	Temperature float64   `json:"temperature"`
	HeatIndex   float64   `json:"heatIndex"`
	Soil        float64   `json:"soil"`
	Water       bool      `json:"water"`
	Lights      bool      `json:"lights"`
	Vent        bool      `json:"vent"`
	Created     time.Time `json:"created"`
}

// StatsField fields
var StatsField = struct {
	ID          string
	IP          string
	Humidity    string
	Temperature string
	HeatIndex   string
	Soil        string
	Water       string
	Lights      string
	Vent        string
	Created     string
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
	Created:     "created",
}

// StatsSortFields sort fields
var StatsSortFields = map[string]string{
	StatsField.Created: "Created",
}

// StatsFilterFields filter fields
var StatsFilterFields = map[string]string{
	StatsField.ID: "ID",
}

// Table implement model table
func (d *Stats) Table() string {
	return StatsTable
}

// SearchField implement model table
func (d *Stats) SearchField() string {
	return StatsField.ID
}

// Mapping table mapping
func (d *Stats) Mapping() (string, error) {
	return mapping.Map(Stats{}).JSON()
}

// GetID implement model get id value
func (d *Stats) GetID() string {
	return d.ID
}

// Service implement model service
func (d *Stats) Service() string {
	return StatsService
}

// Clear implement clear
func (d *Stats) Clear() store.Model {
	return &Stats{}
}

// ClearRel clear relationships
func (d *Stats) ClearRel() {
}

// KeyVal return value for key val search
func (d *Stats) KeyVal() string {
	return d.ID
}

// BRead implement model before read method
func (d *Stats) BRead() error {
	return nil
}

// ARead implement model after read method
func (d *Stats) ARead() error {
	return nil
}

// BCreate implement model before creating method
func (d *Stats) BCreate() error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	d.Created = time.Now()
	return nil
}

// ACreate implement model after creating method
func (d *Stats) ACreate() error {
	return nil
}

// BUpdate implement model before update method
func (d *Stats) BUpdate() error {
	return nil
}

// AUpdate implement model after update method
func (d *Stats) AUpdate() error {
	return nil
}

// BDelete implement model before delete method
func (d *Stats) BDelete() error {
	return nil
}

// ADelete implement model after delete method
func (d *Stats) ADelete() error {
	return nil
}

// ToController adapter to controller
func (d *Stats) ToController() *Controller {
	return &Controller{
		ID:          d.Controller,
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

// RawStats unmarshal raw message
func RawStats(raw json.RawMessage) (*Stats, error) {
	m := &Stats{}
	if err := json.Unmarshal(raw, m); err != nil {
		return nil, errors.InternalWrap(err, "unmarshal document")
	}
	return m, nil
}
