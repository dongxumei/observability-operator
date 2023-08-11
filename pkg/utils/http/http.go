package http

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"gopkg.in/yaml.v3"
)

// WriteJSONResponse writes some JSON as a HTTP response.
func WriteJSONResponse(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// We ignore errors here, because we cannot do anything about them.
	// Write will trigger sending Status code, so we cannot send a different status code afterwards.
	// Also this isn't internal error, but error communicating with client.
	_, _ = w.Write(data)
}

// WriteYAMLResponse writes some YAML as a HTTP response.
func WriteYAMLResponse(w http.ResponseWriter, v interface{}) {
	// There is not standardised content-type for YAML, text/plain ensures the
	// YAML is displayed in the browser instead of offered as a download
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	data, err := yaml.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// We ignore errors here, because we cannot do anything about them.
	// Write will trigger sending Status code, so we cannot send a different status code afterwards.
	// Also this isn't internal error, but error communicating with client.
	_, _ = w.Write(data)
}

// WriteTextResponse sends message as text/plain response with 200 status code.
func WriteTextResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/plain")

	// Ignore inactionable errors.
	_, _ = w.Write([]byte(message))
}

// WriteHTMLResponse sends message as text/html response with 200 status code.
func WriteHTMLResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/html")

	// Ignore inactionable errors.
	_, _ = w.Write([]byte(message))
}

// RenderHTTPResponse either responds with JSON or a rendered HTML page using the passed in template
// by checking the Accepts header.
func RenderHTTPResponse(w http.ResponseWriter, v interface{}, t *template.Template, r *http.Request) {
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		WriteJSONResponse(w, v)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.Execute(w, v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// StreamWriteYAMLResponse stream writes data as http response
func StreamWriteYAMLResponse(w http.ResponseWriter, iter chan interface{}, logger log.Logger) {
	w.Header().Set("Content-Type", "application/yaml")
	for v := range iter {
		data, err := yaml.Marshal(v)
		if err != nil {
			level.Error(logger).Log("msg", "yaml marshal failed", "err", err)
			continue
		}
		_, err = w.Write(data)
		if err != nil {
			level.Error(logger).Log("msg", "write http response failed", "err", err)
			return
		}
	}
}
