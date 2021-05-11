package stats

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/fakovacic/aig/internal/models"
	"github.com/fakovacic/aig/internal/store"
	"github.com/fakovacic/aig/internal/store/client"
	"github.com/fakovacic/aig/internal/store/wrapper"
	"github.com/julienschmidt/httprouter"
)

// New create service struct
func New(s *store.Store) *Service {
	return &Service{
		us: client.New(s, &models.Stats{}),
	}
}

// Service users
type Service struct {
	us *wrapper.Store
}

// ParseBodyRequest parse http request
func (s *Service) ParseBodyRequest(r *http.Request, par httprouter.Params) (*models.Stats, []string, error) {
	u := models.Stats{}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, nil, err
	}

	var reqBody map[string]interface{}
	err = json.Unmarshal(b, &reqBody)
	if err != nil {
		return nil, nil, err
	}

	var updatedFields []string
	fields := []string{
		models.StatsField.ID,
		models.StatsField.IP,
		models.StatsField.Humidity,
		models.StatsField.Temperature,
		models.StatsField.HeatIndex,
		models.StatsField.Soil,
		models.StatsField.Water,
		models.StatsField.Lights,
		models.StatsField.Vent,
	}

	for i := range fields {
		val, ok := reqBody[fields[i]]
		if !ok {
			continue
		}

		switch fields[i] {
		case models.StatsField.ID:
			u.Controller = val.(string)
		case models.StatsField.IP:
			u.IP = val.(string)
			updatedFields = append(updatedFields, models.ControllerField.IP)
		case models.StatsField.Humidity:
			u.Humidity = val.(float64)
			updatedFields = append(updatedFields, models.StatsField.Humidity)
		case models.StatsField.Temperature:
			u.Temperature = val.(float64)
			updatedFields = append(updatedFields, models.StatsField.Temperature)
		case models.StatsField.HeatIndex:
			u.HeatIndex = val.(float64)
			updatedFields = append(updatedFields, models.StatsField.HeatIndex)
		case models.StatsField.Soil:
			u.Soil = val.(float64)
			updatedFields = append(updatedFields, models.StatsField.Soil)
		case models.StatsField.Water:
			u.Water = val.(bool)
			updatedFields = append(updatedFields, models.StatsField.Water)
		case models.StatsField.Lights:
			u.Lights = val.(bool)
			updatedFields = append(updatedFields, models.StatsField.Lights)
		case models.StatsField.Vent:
			u.Vent = val.(bool)
			updatedFields = append(updatedFields, models.StatsField.Vent)
		}
	}

	return &u, updatedFields, nil
}

// Create create action
func (s *Service) Create(ctx context.Context, m *models.Stats) (*models.Stats, error) {
	model, err := s.us.Create(ctx, m)
	if err != nil {
		return nil, err
	}
	return model.(*models.Stats), nil
}
