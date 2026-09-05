package api

import (
	"database/sql"
	"encoding/json"
	"myapp/internal/auth"
	"myapp/internal/database"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type Note struct {
	NoteId      uuid.UUID `json:"note_id"`
	UserId      uuid.UUID `json:"user_id"`
	DailyNote   string    `json:"daily_note"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}

func isConsectiveDay(prev time.Time, curr time.Time) bool {
	prev = prev.UTC()
	curr = curr.UTC()

	prevDate := time.Date(
		prev.Year(), prev.Month(), prev.Day(), 0, 0, 0, 0, time.UTC,
	)

	currDate := time.Date(
		curr.Year(), curr.Month(), curr.Day(), 0, 0, 0, 0, time.UTC,
	)

	return currDate.Equal(prevDate.AddDate(0, 0, 1))
}

func isSameDay(x time.Time, y time.Time) bool {
	x = x.UTC()
	y = y.UTC()
	return x.Year() == y.Year() && x.YearDay() == y.YearDay()
}

func (config *Config) GetNotes(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	OriginalUserId, err := auth.ValidateJWT(token, os.Getenv("JWT_SECRET"))

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	DbNotes, err := config.Notes.GetAllNotes(r.Context(), OriginalUserId)

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

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	OriginalUserId, err := auth.ValidateJWT(token, os.Getenv("JWT_SECRET"))

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	JsonNote := Note{}
	id, err := uuid.Parse(r.PathValue("note_id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Note Id")
		return
	}

	note, err := config.Notes.ReadNote(r.Context(), id)
	if err == sql.ErrNoRows {
		RespondWithError(w, http.StatusNotFound, "Couldn't find Note")
		return
	}

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	if note.UserID != OriginalUserId {
		RespondWithError(w, http.StatusForbidden, "Not Allowed")
		return
	}

	JsonNote.CreatedAt = note.CreatedAt
	JsonNote.DailyNote = note.DailyNote
	JsonNote.LastUpdated = note.UpdatedAt
	JsonNote.NoteId = note.NoteID
	JsonNote.UserId = OriginalUserId

	RespondWithJson(w, http.StatusOK, JsonNote)
}

func (config *Config) AddNote(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	token, err := auth.GetBearerToken(r.Header)
	//println("Tokem isssue")
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}
	OriginalUserId, err := auth.ValidateJWT(token, os.Getenv("JWT_SECRET"))
	//println("Validate jwt issue")
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Not Authorized")
		//println(err.Error())
		return
	}
	//println("passes")
	note := Note{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&note); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Couldn't Read Note")
		return
	}

	DbNote, err := config.Notes.CreateNote(r.Context(), database.CreateNoteParams{
		UserID:    OriginalUserId,
		DailyNote: note.DailyNote,
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't Add Note")
		return
	}

	note.CreatedAt = DbNote.CreatedAt
	note.LastUpdated = DbNote.UpdatedAt
	note.NoteId = DbNote.NoteID
	note.DailyNote = DbNote.DailyNote
	note.UserId = DbNote.UserID

	notes, err := config.Notes.GetSortedNotes(r.Context(), OriginalUserId)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	if len(notes) == 1 {
		if err := config.Users.UpdateStreak(r.Context(), OriginalUserId); err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Server Error")
			return
		}
	} else {
		RecentNote := notes[1]
		consecutive := isConsectiveDay(RecentNote.CreatedAt, note.CreatedAt)
		if consecutive {
			if err := config.Users.UpdateStreak(r.Context(), OriginalUserId); err != nil {
				RespondWithError(w, http.StatusInternalServerError, "Server error")
				return
			}
		} else if !isSameDay(RecentNote.CreatedAt, note.CreatedAt) {
			if err := config.Users.ResetCurrentStreak(r.Context(), OriginalUserId); err != nil {
				RespondWithError(w, http.StatusInternalServerError, "Server Error")
				return
			}
		}
	}

	if err := config.Users.IncreaseNoteCount(r.Context(), OriginalUserId); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}
	RespondWithJson(w, http.StatusOK, note)
}

func (config *Config) DeleteNote(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	OriginalUserId, err := auth.ValidateJWT(token, os.Getenv("JWT_SECRET"))

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	id, err := uuid.Parse(r.PathValue("note_id"))

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Note id")
		return
	}

	note, err := config.Notes.ReadNote(r.Context(), id)

	if err == sql.ErrNoRows {
		RespondWithError(w, http.StatusNotFound, "Couldn't find Note")
		return
	}

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	if note.UserID != OriginalUserId {
		RespondWithError(w, http.StatusForbidden, "Not Allowed")
		return
	}

	if err = config.Notes.DeleteNote(r.Context(), id); err == sql.ErrNoRows {
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

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	OriginalUserId, err := auth.ValidateJWT(token, os.Getenv("JWT_SECRET"))

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	id, err := uuid.Parse(r.PathValue("note_id"))

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid note id")
		return
	}

	DBnote, err := config.Notes.ReadNote(r.Context(), id)

	if err == sql.ErrNoRows {
		RespondWithError(w, http.StatusNotFound, "Couldn't find Note")
		return
	}

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	if DBnote.UserID != OriginalUserId {
		RespondWithError(w, http.StatusForbidden, "Not Allowed")
		return
	}

	defer r.Body.Close()

	note := Note{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&note); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Couldn't Update the Note")
		return
	}

	DbNote, err := config.Notes.UpdateNote(r.Context(), database.UpdateNoteParams{
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
