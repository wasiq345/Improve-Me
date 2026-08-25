package api

import (
	"database/sql"
	"encoding/json"
	"myapp/internal/database"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Note struct {
	NoteId      uuid.UUID `json:"note_id"`
	DailyNote   string    `json:"daily_note"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}

func (config *Config) DashBoard(w http.ResponseWriter, r *http.Request) {
	RespondWithJson(w, http.StatusOK, "Welcome To DashBoard")
}

func (config *Config) GetNotes(w http.ResponseWriter, r *http.Request) {

	DbNotes, err := config.DB.GetAllNotes(r.Context())

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	notes := make([]Note, len(DbNotes))

	for i, DbNote := range DbNotes {
		notes[i].NoteId = DbNote.NoteID
		notes[i].CreatedAt = DbNote.CreatedAt
		notes[i].LastUpdated = DbNote.UpdatedAt
		notes[i].DailyNote = DbNote.DailyNote
	}

	RespondWithJson(w, http.StatusOK, notes)
}

func (config *Config) ReadNote(w http.ResponseWriter, r *http.Request) {

	JsonNote := Note{}
	id, err := uuid.Parse(r.PathValue("note_id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Note Id")
		return
	}

	note, err := config.DB.ReadNote(r.Context(), id)
	if err == sql.ErrNoRows {
		RespondWithError(w, http.StatusNotFound, "Couldn't find Note")
		return
	}

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	JsonNote.CreatedAt = note.CreatedAt
	JsonNote.DailyNote = note.DailyNote
	JsonNote.LastUpdated = note.UpdatedAt
	JsonNote.NoteId = note.NoteID

	RespondWithJson(w, http.StatusOK, JsonNote)
}

func (config *Config) AddNote(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	note := Note{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&note); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Couldn't Read Note")
		return
	}

	DbNote, err := config.DB.CreateNote(r.Context(), note.DailyNote)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't Add Note")
		return
	}

	note.CreatedAt = DbNote.CreatedAt
	note.LastUpdated = DbNote.UpdatedAt
	note.NoteId = DbNote.NoteID
	note.DailyNote = DbNote.DailyNote

	RespondWithJson(w, http.StatusOK, note)
}

func (config *Config) DeleteNote(w http.ResponseWriter, r *http.Request) {

	id, err := uuid.Parse(r.PathValue("note_id"))

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Note id")
		return
	}

	if err = config.DB.DeleteNote(r.Context(), id); err == sql.ErrNoRows {
		RespondWithError(w, http.StatusNotFound, "Couldn't Find Note")
		return
	}
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	RespondWithJson(w, http.StatusOK, "Note Deleted Succesfully")
}

func (config *Config) UpdateNote(w http.ResponseWriter, r *http.Request) {

	id, err := uuid.Parse(r.PathValue("note_id"))

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid note id")
		return
	}
	defer r.Body.Close()

	note := Note{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&note); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Couldn't Update the Note")
		return
	}

	DbNote, err := config.DB.UpdateNote(r.Context(), database.UpdateNoteParams{
		NoteID:    id,
		DailyNote: note.DailyNote,
	})

	if err == sql.ErrNoRows {
		RespondWithError(w, http.StatusNotFound, "Couldn't find note")
		return
	}
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}
	note.CreatedAt = DbNote.CreatedAt
	note.LastUpdated = DbNote.UpdatedAt
	note.NoteId = DbNote.NoteID
	note.DailyNote = DbNote.DailyNote
	RespondWithJson(w, http.StatusOK, note)
}
