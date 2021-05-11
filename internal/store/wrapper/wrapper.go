package wrapper

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/fakovacic/aig/internal/errors"
	"github.com/fakovacic/aig/internal/meta"
	"github.com/fakovacic/aig/internal/store"
	"github.com/go-playground/validator"
	"go.elastic.co/apm"
)

// New create new model store
func New(m store.ModelStore) *Store {
	return &Store{
		m:     m,
		model: m.GetModel(),
		valid: validator.New(),
	}
}

// Store store wraper with predefined functions
type Store struct {
	m     store.ModelStore
	model store.Model
	valid *validator.Validate
}

// Init create, truncate, drop model table
func (ms Store) Init(ctx context.Context) {
	ms.m.Init(ctx)
}

// Model return model
func (ms Store) Model(ctx context.Context) store.Model {
	return ms.model
}

// Create add row in database
func (ms Store) Create(ctx context.Context, model store.Model) (store.Model, error) {
	err := model.BCreate()
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return nil, errors.InternalWrap(err, "before create:")
	}

	err = ms.Validate(model)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return nil, errors.BadRequestWrap(err, "validate:")
	}

	model, err = ms.m.Create(ctx, model)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	err = model.ACreate()
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return nil, errors.InternalWrap(err, "after create:")
	}
	return model, nil
}

// ReadBy read first row from database by given field and value
func (ms Store) ReadBy(ctx context.Context, field string, value interface{}) (store.Model, error) {
	ms.model = ms.model.Clear()
	err := ms.model.BRead()
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return nil, errors.InternalWrap(err, "before read:")
	}

	ms.model, err = ms.m.ReadBy(ctx, field, value)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	err = ms.model.ARead()
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return nil, errors.InternalWrap(err, "after read:")
	}
	return ms.model, nil
}

// Update update row in database
func (ms Store) Update(ctx context.Context, model store.Model, fields []string) (store.Model, error) {
	err := model.BUpdate()
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return nil, errors.InternalWrap(err, "before read:")
	}

	err = ms.Validate(model)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return nil, errors.BadRequestWrap(err, "validate:")
	}

	model, err = ms.m.Update(ctx, model, fields)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	err = model.AUpdate()
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return nil, errors.InternalWrap(err, "after read:")
	}
	return model, nil
}

// Delete deletes row from database
func (ms Store) Delete(ctx context.Context, model store.Model) error {
	err := model.BDelete()
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return errors.InternalWrap(err, "before delete:")
	}

	err = ms.m.Delete(ctx, model)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return err
	}

	err = model.ADelete()
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return errors.InternalWrap(err, "after delete:")
	}
	return nil
}

// List return list of models
func (ms Store) List(ctx context.Context, m *meta.Model) ([]json.RawMessage, error) {
	list, err := ms.m.List(ctx, m)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	return list, nil
}

// ListByID return list of models by id
func (ms Store) ListByID(ctx context.Context, idList ...store.IDf) (map[string]map[string]json.RawMessage, error) {
	list, err := ms.m.ListByID(ctx, idList...)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	return list, nil
}

// ListIDs return list of by ids
func (ms Store) ListIDs(ctx context.Context, idList ...store.IDs) ([]store.IDf, error) {
	list, err := ms.m.ListIDs(ctx, idList...)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	return list, nil
}

// Validate validate struct before send to database layer
func (ms Store) Validate(model store.Model) error {
	err := ms.valid.Struct(model)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		var invalidFields []string
		for _, err := range err.(validator.ValidationErrors) {
			invalidFields = append(invalidFields, err.Field())
		}

		return errors.BadRequest("fields not valid:%s", strings.Join(invalidFields, ", "))
	}

	return nil
}
