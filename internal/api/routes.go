package api

import (
	"net/http"
)

func (config *Config) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/Dashboard", config.DashBoard)
	mux.HandleFunc("GET /Dashboard/ReadNote/{note_id}", config.ReadNote)
	mux.HandleFunc("DELETE /Dashboard/DeleteNote/{note_id}", config.DeleteNote)
	mux.HandleFunc("PUT /Dashboard/UpdateNote/{note_id}", config.UpdateNote)
	mux.HandleFunc("POST /Dashboard/CreateNote", config.AddNote)
	mux.HandleFunc("GET /Dashboard/GetNotes", config.GetNotes)
	mux.HandleFunc("POST /Dashboard/ResgiterUser", config.RegisterUser)
	mux.HandleFunc("POST /Dashboard/LoginUser", config.LoginUser)
}
