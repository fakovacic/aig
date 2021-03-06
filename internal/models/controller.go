package models

import (
	"encoding/json"
	"strconv"
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
	Config      []Config  `json:"config"`
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

// Config controller struct
type Config struct {
	Stat    string `json:"stat"`
	From    string `json:"from"`
	To      string `json:"to"`
	Control string `json:"control"`
	Power   string `json:"power"`
}

// ControllerConfigField fields
var ControllerConfigField = struct {
	Stat    string
	From    string
	To      string
	Control string
	Power   string
}{
	Stat:    "stat",
	From:    "from",
	To:      "to",
	Control: "control",
	Power:   "power",
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

// RemoveConfig remove config key
func (d *Controller) RemoveConfig(i int) {
	d.Config[i] = d.Config[len(d.Config)-1]
	d.Config = d.Config[:len(d.Config)-1]
}

// ConfigPower return bool power
func (d *Controller) ConfigPower(s string) bool {
	return s == "on"
}

// ConfigPowerOp return oposite bool power
func (d *Controller) ConfigPowerOp(s string) bool {
	return s != "on"
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

// ValidCommands return command that need to execute
func (d *Controller) ValidCommands() map[string]string {
	m := make(map[string]string)

	if len(d.Config) == 0 {
		return m
	}

	for i := range d.Config {
		power := d.ConfigPowerOp(d.Config[i].Power)
		switch d.Config[i].Stat {
		case "time":
			now := time.Now()
			from, _ := time.Parse("15:04:05", d.Config[i].From)
			to, _ := time.Parse("15:04:05", d.Config[i].To)

			if from.After(now) && to.Before(now) {
				power = d.ConfigPower(d.Config[i].Power)
			}

		case "temperature":
			from, _ := strconv.ParseFloat(d.Config[i].From, 64)
			to, _ := strconv.ParseFloat(d.Config[i].To, 64)

			if d.Temperature > from && d.Temperature < to {
				power = d.ConfigPower(d.Config[i].Power)
			}
		case "humidity":
			from, _ := strconv.ParseFloat(d.Config[i].From, 64)
			to, _ := strconv.ParseFloat(d.Config[i].To, 64)

			if d.Humidity > from && d.Humidity < to {
				power = d.ConfigPower(d.Config[i].Power)
			}
		case "soil":
			from, _ := strconv.ParseFloat(d.Config[i].From, 64)
			to, _ := strconv.ParseFloat(d.Config[i].To, 64)

			if d.Soil > from && d.Soil < to {
				power = d.ConfigPower(d.Config[i].Power)
			}
		}

		switch d.Config[i].Control {
		case "water":
			if d.Water != power {
				m[d.Config[i].Control] = d.Config[i].Power
			}
		case "vent":
			if d.Vent != power {
				m[d.Config[i].Control] = d.Config[i].Power
			}
		case "lights":
			if d.Lights != power {
				m[d.Config[i].Control] = d.Config[i].Power
			}
			// case "heater":
			// 	if d.Heater != power {
			// 		m[d.Config[i].Control] = d.Config[i].Power
			// 	}
		}
	}

	return m
}

// RawController unmarshal raw message
func RawController(raw json.RawMessage) (*Controller, error) {
	m := &Controller{}
	if err := json.Unmarshal(raw, m); err != nil {
		return nil, errors.InternalWrap(err, "unmarshal document")
	}
	return m, nil
}
