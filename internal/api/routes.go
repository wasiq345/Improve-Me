package api

import (
	"net/http"
)

func (config *Config) RegisterRoutes(mux *http.ServeMux) {
	//fileserver := http.FileServer(http.Dir("./frontend"))
	//mux.Handle("/Dashboard/", http.StripPrefix("/Dashboard", fileserver))
	mux.HandleFunc("GET /Profile/{username}", config.UserProfile)
	mux.HandleFunc("GET /Profile/{username}/ReadNote/{note_id}", config.ReadNote)
	mux.HandleFunc("DELETE /Profile/{username}/DeleteNote/{note_id}", config.DeleteNote)
	mux.HandleFunc("PUT /Profile/{username}/UpdateNote/{note_id}", config.UpdateNote)
	mux.HandleFunc("POST /Profile/{username}/CreateNote", config.AddNote)
	mux.HandleFunc("GET /Profile/{username}/GetNotes", config.GetNotes)
	mux.HandleFunc("POST /Dashboard/RegisterUser", config.RegisterUser)
	mux.HandleFunc("POST /Dashboard/LoginUser", config.LoginUser)
	mux.HandleFunc("POST /api/revoke", config.Revoke)
	mux.HandleFunc("POST /api/refresh", config.Refresh)
}
