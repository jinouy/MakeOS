package render

import "net/http"

type Render interface {
	Render(w http.ResponseWriter, code int) error
	WritContentType(w http.ResponseWriter)
}

func writeContentType(w http.ResponseWriter, value string) {
	w.Header().Set("Content-Type", value)
}
