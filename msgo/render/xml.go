package render

import (
	"encoding/xml"
	"net/http"
)

type XML struct {
	Data any
}

func (x *XML) Render(w http.ResponseWriter) error {
	x.WritContentType(w)
	err := xml.NewEncoder(w).Encode(x.Data)
	return err
}

func (x *XML) WritContentType(w http.ResponseWriter) {
	writeContentType(w, "application/xml; charset=utf-8")
}
