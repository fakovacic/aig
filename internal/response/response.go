package response

import (
	"github.com/fakovacic/aig/internal/types/content"
	"github.com/fakovacic/aig/pkg/templates"
)

// Response type
type Response struct {
	Data        interface{}
	Redirect    interface{}
	ContentType content.ContentType
	Template    *templates.Template
}
