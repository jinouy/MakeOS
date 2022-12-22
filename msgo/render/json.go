package render

import (
	"encoding/json"
	"net/http"
)

type JSON struct {
	Data any
}

func (j *JSON) Render(w http.ResponseWriter) error {
	j.WritContentType(w)
	jsonData, err := json.Marshal(j.Data)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonData)
	return err
}

func (j *JSON) WritContentType(w http.ResponseWriter) {
	writeContentType(w, "application/json; charset=utf-8")
}
