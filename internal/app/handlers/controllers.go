package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fakovacic/aig/internal/app/actions"
	"github.com/fakovacic/aig/internal/app/config"
	"github.com/fakovacic/aig/internal/app/endpoints"
	pg "github.com/fakovacic/aig/internal/app/pages"
	"github.com/fakovacic/aig/internal/errors"
	"github.com/fakovacic/aig/internal/handlers"
	"github.com/fakovacic/aig/internal/log"
	"github.com/fakovacic/aig/internal/meta"
	"github.com/fakovacic/aig/internal/models"
	"github.com/fakovacic/aig/internal/response"
	"github.com/fakovacic/aig/internal/services/activities"
	"github.com/fakovacic/aig/internal/services/controllers"
	"github.com/fakovacic/aig/internal/services/stats"
	"github.com/fakovacic/aig/internal/types/apps"
	"github.com/fakovacic/aig/internal/types/content"
	"github.com/julienschmidt/httprouter"
)

var baseRedirect = endpoints.Dashoard

// New return new handler
func New(cfg config.Config) handlers.Handler {
	return Handler{
		cs:  controllers.New(cfg.Store()),
		aa:  activities.New(cfg.Store()),
		ss:  stats.New(cfg.Store()),
		cfg: cfg,
	}
}

// Handler handler struct
type Handler struct {
	aa  *activities.Service
	cs  *controllers.Service
	ss  *stats.Service
	cfg config.Config
}

// AppType implement handler interface
func (m Handler) AppType() apps.AppType {
	return apps.WEB
}

// ClearNotifications implement handler interface
func (m Handler) ClearNotifications() {
	m.cfg.ClearNotifications()
}

// Action execute service action
func (m Handler) Action(r *http.Request, par httprouter.Params, action handlers.Action) (res *response.Response, err error) {
	metaModel := meta.New()
	metaModel.SortFields = models.ControllerSortFields
	metaModel.FilterFields = models.ControllerFilterFields
	err = metaModel.ParseRequest(r)
	if err != nil {
		return &response.Response{
			Data:        err,
			ContentType: content.HTML,
		}, err
	}

	ctx := context.WithValue(r.Context(), "action", action)

	switch action {
	case actions.Dashoard:
		return m.Dashoard(ctx, metaModel)
	case actions.Controller:
		model, err := m.cs.ParseFormRequest(r, par)
		if err != nil {
			return &response.Response{
				Data:        err,
				ContentType: content.HTML,
			}, err
		}

		return m.Controller(ctx, model)
	case actions.ControllerConfigAdd:
		model, err := m.cs.ParseFormRequest(r, par)
		if err != nil {
			return &response.Response{
				Data:        err,
				ContentType: content.HTML,
			}, err
		}

		return m.ControllerConfigAdd(ctx, model)
	case actions.ControllerConfigRemove:
		model, err := m.cs.ParseFormRequest(r, par)
		if err != nil {
			return &response.Response{
				Data:        err,
				ContentType: content.HTML,
			}, err
		}

		return m.ControllerConfigRemove(ctx, model)
	case actions.Register:
		model, err := m.cs.ParseBodyRequest(r, par)
		if err != nil {
			return &response.Response{
				Data:        err,
				ContentType: content.JSON,
			}, err
		}

		return m.Register(ctx, model)
	default:
		m.cfg.AddNotification("error", "not implemented action")
		return &response.Response{
			Redirect:    endpoints.Dashoard,
			ContentType: content.HTML,
		}, nil
	}
}

type dashboardPage struct {
	pg.Page
	Meta *meta.Model
	Data []*models.Controller
}

// Dashoard action
func (m *Handler) Dashoard(ctx context.Context, metaModel *meta.Model) (res *response.Response, err error) {
	list, err := m.cs.List(ctx, metaModel)
	if err != nil {
		return &response.Response{
			ContentType: content.HTML,
		}, err
	}
	m.aa.Create(ctx)
	return &response.Response{
		Data: dashboardPage{
			pg.Page{
				Title:         "Dashboard",
				Notifications: m.cfg.GetNotifications(),
			},
			metaModel,
			list,
		},
		ContentType: content.HTML,
	}, nil
}

// Controller action
func (m *Handler) Controller(ctx context.Context, model *models.Controller) (res *response.Response, err error) {
	ctx = context.WithValue(ctx, "external_id", model.ID)
	_, err = m.cs.Update(ctx, model)
	if err != nil {
		m.cfg.AddNotification("error", fmt.Sprintf("update error:%v", err))
		return &response.Response{
			Redirect:    endpoints.Dashoard,
			ContentType: content.HTML,
		}, nil
	}

	m.aa.Create(ctx)
	m.cfg.AddNotification("success", "config changed")
	return &response.Response{
		Redirect:    endpoints.Dashoard,
		ContentType: content.HTML,
	}, nil
}

// ControllerConfigAdd action
func (m *Handler) ControllerConfigAdd(ctx context.Context, model *models.Controller) (res *response.Response, err error) {
	ctx = context.WithValue(ctx, "external_id", model.ID)

	model.Config = append(model.Config, models.Config{})
	_, err = m.cs.Update(ctx, model)
	if err != nil {
		m.cfg.AddNotification("error", fmt.Sprintf("update error:%v", err))
		return &response.Response{
			Redirect:    endpoints.Dashoard,
			ContentType: content.HTML,
		}, nil
	}

	m.aa.Create(ctx)
	m.cfg.AddNotification("success", "add controller config")
	return &response.Response{
		Redirect:    endpoints.Dashoard,
		ContentType: content.HTML,
	}, nil
}

// ControllerConfigRemove action
func (m *Handler) ControllerConfigRemove(ctx context.Context, model *models.Controller) (res *response.Response, err error) {
	ctx = context.WithValue(ctx, "external_id", model.ID)
	_, err = m.cs.Update(ctx, model)
	if err != nil {
		m.cfg.AddNotification("error", fmt.Sprintf("update error:%v", err))
		return &response.Response{
			Redirect:    endpoints.Dashoard,
			ContentType: content.HTML,
		}, nil
	}

	m.aa.Create(ctx)
	m.cfg.AddNotification("success", "remove controller config")
	return &response.Response{
		Redirect:    endpoints.Dashoard,
		ContentType: content.HTML,
	}, nil
}

// Register action
func (m *Handler) Register(ctx context.Context, model *models.Controller) (res *response.Response, err error) {
	ctx = context.WithValue(ctx, "external_id", model.ID)
	m.aa.Create(ctx)

	_, err = m.cs.Read(ctx, "id", model.ID)
	if err != nil {
		val, ok := err.(errors.Error)
		if !ok {
			return &response.Response{
				Data:        err,
				ContentType: content.JSON,
			}, nil
		}
		if val.Status != http.StatusNotFound {
			return &response.Response{
				Data:        err,
				ContentType: content.JSON,
			}, nil
		}
		// create
		model, err = m.cs.Create(ctx, model)
		if err != nil {
			log.Error(err)
			return &response.Response{
				Data:        err,
				ContentType: content.JSON,
			}, nil
		}

		return &response.Response{
			Data:        model,
			ContentType: content.JSON,
		}, nil
	}

	// update
	model, err = m.cs.Update(ctx, model)
	if err != nil {
		return &response.Response{
			Data:        err,
			ContentType: content.JSON,
		}, nil
	}

	return &response.Response{
		Data:        model,
		ContentType: content.JSON,
	}, nil
}
