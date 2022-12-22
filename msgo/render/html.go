package render

import (
	"github.com/jinouy/msgo/internal/bytesconv"
	"html/template"
	"net/http"
)

type HTML struct {
	Data       any
	Name       string
	Template   *template.Template
	IsTemplate bool
}
type HTMLRender struct {
	Template *template.Template
}

func (h *HTML) Render(w http.ResponseWriter) error {
	h.WritContentType(w)
	if h.IsTemplate {
		err := h.Template.ExecuteTemplate(w, h.Name, h.Data)
		return err
	}
	_, err := w.Write(bytesconv.StringToBytes(h.Data.(string)))
	return err
}

func (h *HTML) WritContentType(w http.ResponseWriter) {
	writeContentType(w, "text/html; charset=utf-8")
}
