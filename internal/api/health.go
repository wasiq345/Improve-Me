package api

import "net/http"

func (config *Config) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{'"status": "ok"}`))
}
