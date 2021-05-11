package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fakovacic/aig/internal/errors"
	"github.com/fakovacic/aig/internal/meta"
	"github.com/fakovacic/aig/internal/models"
	"github.com/fakovacic/aig/internal/store"
	"github.com/fakovacic/aig/internal/store/client"
	"github.com/fakovacic/aig/internal/store/wrapper"
	"github.com/fakovacic/aig/internal/types/methods"
	"github.com/julienschmidt/httprouter"
	"go.elastic.co/apm"
)

// New create service struct
func New(s *store.Store) *Service {
	return &Service{
		us: client.New(s, &models.Controller{}),
	}
}

// Service users
type Service struct {
	us *wrapper.Store
}

// ParseBodyRequest parse http request
func (s *Service) ParseBodyRequest(r *http.Request, par httprouter.Params) (*models.Controller, []string, error) {
	u := models.Controller{}

	if par != nil {
		id := par.ByName("id")
		if id != "" {
			u.ID = id
		}
	}

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
		models.ControllerField.ID,
		models.ControllerField.IP,
		models.ControllerField.Humidity,
		models.ControllerField.Temperature,
		models.ControllerField.HeatIndex,
		models.ControllerField.Soil,
		models.ControllerField.Water,
		models.ControllerField.Lights,
		models.ControllerField.Vent,
		models.ControllerField.Online,
	}

	for i := range fields {
		val, ok := reqBody[fields[i]]
		if !ok {
			continue
		}

		switch fields[i] {
		case models.ControllerField.ID:
			u.ID = val.(string)
		case models.ControllerField.IP:
			u.IP = val.(string)
			updatedFields = append(updatedFields, models.ControllerField.IP)
		case models.ControllerField.Humidity:
			u.Humidity = val.(float64)
			updatedFields = append(updatedFields, models.ControllerField.Humidity)
		case models.ControllerField.Temperature:
			u.Temperature = val.(float64)
			updatedFields = append(updatedFields, models.ControllerField.Temperature)
		case models.ControllerField.HeatIndex:
			u.HeatIndex = val.(float64)
			updatedFields = append(updatedFields, models.ControllerField.HeatIndex)
		case models.ControllerField.Soil:
			u.Soil = val.(float64)
			updatedFields = append(updatedFields, models.ControllerField.Soil)
		case models.ControllerField.Water:
			u.Water = val.(bool)
			updatedFields = append(updatedFields, models.ControllerField.Water)
		case models.ControllerField.Lights:
			u.Lights = val.(bool)
			updatedFields = append(updatedFields, models.ControllerField.Lights)
		case models.ControllerField.Vent:
			u.Vent = val.(bool)
			updatedFields = append(updatedFields, models.ControllerField.Vent)
		case models.ControllerField.Online:
			u.Online = val.(bool)
			updatedFields = append(updatedFields, models.ControllerField.Online)
		}
	}

	return &u, updatedFields, nil
}

// ControllerCommand read action
func (s *Service) ControllerCommand(ctx context.Context, id, command, power string) (*models.Controller, error) {
	model, err := s.us.ReadBy(ctx, "id", id)
	if err != nil {
		return nil, errors.BadRequestWrap(err, "read controller by id %s", id)
	}

	m := model.(*models.Controller)
	if m.IP == "" {
		return nil, errors.BadRequest("controller ip empty")
	}

	url := fmt.Sprintf("http://%s/%s/%s", m.IP, command, power)
	if command == "stats" {
		url = fmt.Sprintf("http://%s/stats", m.IP)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		methods.Get.String(),
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var mc models.Controller

	err = json.Unmarshal(b, &mc)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.BadRequest("status code %d", res.StatusCode)
	}

	return &mc, nil
}

// Read read action
func (s *Service) Read(ctx context.Context, field string, value interface{}) (*models.Controller, error) {
	model, err := s.us.ReadBy(ctx, field, value)
	if err != nil {
		return nil, err
	}
	return model.(*models.Controller), nil
}

// Create create action
func (s *Service) Create(ctx context.Context, m *models.Controller) (*models.Controller, error) {
	model, err := s.us.Create(ctx, m)
	if err != nil {
		return nil, err
	}
	return model.(*models.Controller), nil
}

// Update update action
func (s *Service) Update(ctx context.Context, m *models.Controller, fields []string) (*models.Controller, error) {
	if len(fields) == 0 {
		return m, nil
	}
	model, err := s.us.Update(ctx, m, fields)
	if err != nil {
		return nil, err
	}
	return model.(*models.Controller), nil
}

// Delete delete action
func (s *Service) Delete(ctx context.Context, m *models.Controller) error {
	return s.us.Delete(ctx, m)
}

// List action
func (s *Service) List(ctx context.Context, m *meta.Model) ([]*models.Controller, error) {
	modelsList, err := s.us.List(ctx, m)
	if err != nil {
		return nil, err
	}

	list := make([]*models.Controller, len(modelsList))
	for i := range modelsList {
		m, err := models.RawController(modelsList[i])
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			return nil, err
		}
		list[i] = m
	}

	return list, nil
}

// KeyVal action
func (s *Service) KeyVal(ctx context.Context, m *meta.Model) (map[string]string, error) {
	modelsList, err := s.us.List(ctx, m)
	if err != nil {
		return nil, err
	}

	list := make(map[string]string, len(modelsList))
	for i := range modelsList {
		m, err := models.RawController(modelsList[i])
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			return list, err
		}
		list[m.ID] = m.KeyVal()
	}
	return list, nil
}
