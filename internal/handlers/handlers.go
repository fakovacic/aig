package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fakovacic/aig/internal/errors"
	"github.com/fakovacic/aig/internal/log"
	"github.com/fakovacic/aig/internal/response"
	"github.com/fakovacic/aig/internal/types/apps"
	"github.com/fakovacic/aig/internal/types/content"
	"go.elastic.co/apm"

	"github.com/fakovacic/aig/pkg/templates"
	"github.com/julienschmidt/httprouter"
)

// Action interface
type Action interface {
	String() string
}

// Handler execute web template
type Handler interface {
	Action(*http.Request, httprouter.Params, Action) (*response.Response, error)
	ClearNotifications()
	AppType() apps.AppType
}

// Handle handle will execute template for web
func Handle(template *templates.Template, mng Handler, action Action) httprouter.Handle {
	if template != nil {
		err := template.Parse()
		if err != nil {
			fmt.Println("error parsing template", err)
		}
	}
	return func(w http.ResponseWriter, r *http.Request, par httprouter.Params) {
		s := apm.TransactionFromContext(r.Context()).
			StartSpan("handler", action.String(), nil)
		defer s.End()

		log.Debugf("request- params:%v action:%s", par, action)

		res, err := mng.Action(r, par, action)
		log.Debugf("handler action - res:%v error: %v", res, err)
		if res.Template != nil {
			template = res.Template

			err := template.Parse()
			if err != nil {
				fmt.Println("error parsing template", err)
			}
		}

		if res.Redirect != nil {
			http.Redirect(w, r, fmt.Sprintf("%v", res.Redirect), http.StatusFound)
			return
		}

		switch res.ContentType {
		case content.HTML:
			if err != nil {
				fmt.Println("error executing action", err)
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			if template != nil {
				err = template.Execute(w, res.Data)
				if err != nil {
					fmt.Println("error executing template", err)
					http.Redirect(w, r, "/", http.StatusFound)
					return
				}
			}

			mng.ClearNotifications()
		case content.JSON:
			if err != nil {
				val, ok := err.(errors.Error)
				if ok {
					http.Error(w, val.Error(), val.GetStatus())
					return
				}

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			js, err := json.Marshal(res.Data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(js)
			return
		}
	}
}
