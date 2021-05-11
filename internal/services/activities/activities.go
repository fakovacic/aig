package activities

import (
	"context"
	"fmt"

	"github.com/fakovacic/aig/internal/meta"
	"github.com/fakovacic/aig/internal/models"
	"github.com/fakovacic/aig/internal/store"
	"github.com/fakovacic/aig/internal/store/client"
	"github.com/fakovacic/aig/internal/store/wrapper"
	"go.elastic.co/apm"
)

// New create service struct
func New(s *store.Store) *Service {
	return &Service{
		us: client.New(s, &models.Activity{}),
	}
}

// Service struct
type Service struct {
	us *wrapper.Store
}

// Create create action
func (s *Service) Create(ctx context.Context) {
	activ := &models.Activity{}

	act := ctx.Value("action")
	if act != nil {
		activ.Action = fmt.Sprintf("%v", act)
	}

	exID := ctx.Value("external_id")
	if exID != nil {
		activ.ExteralID = fmt.Sprintf("%v", exID)
	}
	_, err := s.us.Create(ctx, activ)
	if err != nil {
		fmt.Println(err)
		apm.CaptureError(ctx, err).Send()
	}
}

// List action
func (s *Service) List(ctx context.Context) ([]*models.Activity, error) {
	m := &meta.Model{
		Page:  1,
		Limit: 100,
	}

	modelsList, err := s.us.List(ctx, m)
	if err != nil {
		return nil, err
	}

	if len(modelsList) == 0 {
		return nil, nil
	}

	list := make([]*models.Activity, len(modelsList))
	for i := range modelsList {
		m, err := models.RawActivity(modelsList[i])
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			return list, err
		}
		list[i] = m
	}

	return list, nil
}
