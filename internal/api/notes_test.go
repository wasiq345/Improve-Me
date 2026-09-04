package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"myapp/internal/auth"
	"myapp/internal/database"

	"github.com/google/uuid"
)

type fakeNoteStore struct {
	createNoteFunc     func(ctx context.Context, arg database.CreateNoteParams) (database.Note, error)
	deleteNoteFunc     func(ctx context.Context, noteId uuid.UUID) error
	getAllNotesFunc    func(ctx context.Context, userId uuid.UUID) ([]database.Note, error)
	readNoteFunc       func(ctx context.Context, noteId uuid.UUID) (database.Note, error)
	updateNoteFunc     func(ctx context.Context, arg database.UpdateNoteParams) (database.Note, error)
	getSortedNotesFunc func(ctx context.Context, userId uuid.UUID) ([]database.Note, error)
}

func (f *fakeNoteStore) CreateNote(ctx context.Context, arg database.CreateNoteParams) (database.Note, error) {
	return f.createNoteFunc(ctx, arg)
}
func (f *fakeNoteStore) DeleteNote(ctx context.Context, noteId uuid.UUID) error {
	return f.deleteNoteFunc(ctx, noteId)
}
func (f *fakeNoteStore) GetAllNotes(ctx context.Context, userId uuid.UUID) ([]database.Note, error) {
	return f.getAllNotesFunc(ctx, userId)
}
func (f *fakeNoteStore) ReadNote(ctx context.Context, noteId uuid.UUID) (database.Note, error) {
	return f.readNoteFunc(ctx, noteId)
}
func (f *fakeNoteStore) UpdateNote(ctx context.Context, arg database.UpdateNoteParams) (database.Note, error) {
	return f.updateNoteFunc(ctx, arg)
}

func (f *fakeNoteStore) GetSortedNotes(ctx context.Context, userId uuid.UUID) ([]database.Note, error) {
	return f.getSortedNotesFunc(ctx, userId)
}

func newAuthedRequest(t *testing.T, method, path string, userID uuid.UUID, body any) *http.Request {
	t.Helper()

	token, err := auth.MakeJWT(userID, os.Getenv("JWT_SECRET"), time.Hour)
	if err != nil {
		t.Fatalf("test setup: failed to create JWT: %v", err)
	}

	var req *http.Request
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("test setup: failed to marshal request body: %v", err)
		}
		req = httptest.NewRequest(method, path, bytes.NewReader(raw))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func withNoteID(req *http.Request, noteID uuid.UUID) *http.Request {
	req.SetPathValue("note_id", noteID.String())
	return req
}

func TestGetNotes_Success(t *testing.T) {
	userID := uuid.New()
	noteID := uuid.New()

	store := &fakeNoteStore{
		getAllNotesFunc: func(ctx context.Context, userId uuid.UUID) ([]database.Note, error) {
			if userId != userID {
				t.Errorf("expected notes to be fetched for user %v, got %v", userID, userId)
			}
			return []database.Note{
				{NoteID: noteID, UserID: userID, DailyNote: "wrote some tests today"},
			}, nil
		},
	}
	config := &Config{Notes: store}

	req := newAuthedRequest(t, http.MethodGet, "/notes", userID, nil)
	rec := httptest.NewRecorder()

	config.GetNotes(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusOK, rec.Code, rec.Body.String())
	}

	var got []Note
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	if len(got) != 1 || got[0].NoteId != noteID {
		t.Errorf("expected one note with id %v, got %+v", noteID, got)
	}
}

func TestGetNotes_NoToken(t *testing.T) {
	config := &Config{Notes: &fakeNoteStore{}}
	req := httptest.NewRequest(http.MethodGet, "/notes", nil)
	rec := httptest.NewRecorder()

	config.GetNotes(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d with no token, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestGetNotes_DatabaseError(t *testing.T) {
	userID := uuid.New()
	store := &fakeNoteStore{
		getAllNotesFunc: func(ctx context.Context, userId uuid.UUID) ([]database.Note, error) {
			return nil, errors.New("connection refused")
		},
	}
	config := &Config{Notes: store}

	req := newAuthedRequest(t, http.MethodGet, "/notes", userID, nil)
	rec := httptest.NewRecorder()

	config.GetNotes(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestReadNote_Success(t *testing.T) {
	userID := uuid.New()
	noteID := uuid.New()

	store := &fakeNoteStore{
		readNoteFunc: func(ctx context.Context, id uuid.UUID) (database.Note, error) {
			return database.Note{NoteID: id, UserID: userID, DailyNote: "hello"}, nil
		},
	}
	config := &Config{Notes: store}

	req := newAuthedRequest(t, http.MethodGet, "/notes/"+noteID.String(), userID, nil)
	req = withNoteID(req, noteID)
	rec := httptest.NewRecorder()

	config.ReadNote(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusOK, rec.Code, rec.Body.String())
	}
}

func TestReadNote_NotFound(t *testing.T) {
	userID := uuid.New()
	noteID := uuid.New()

	store := &fakeNoteStore{
		readNoteFunc: func(ctx context.Context, id uuid.UUID) (database.Note, error) {
			return database.Note{}, sql.ErrNoRows
		},
	}
	config := &Config{Notes: store}

	req := newAuthedRequest(t, http.MethodGet, "/notes/"+noteID.String(), userID, nil)
	req = withNoteID(req, noteID)
	rec := httptest.NewRecorder()

	config.ReadNote(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestReadNote_ForbiddenWhenNoteBelongsToSomeoneElse(t *testing.T) {
	requestingUser := uuid.New()
	noteOwner := uuid.New()
	noteID := uuid.New()

	store := &fakeNoteStore{
		readNoteFunc: func(ctx context.Context, id uuid.UUID) (database.Note, error) {
			return database.Note{NoteID: id, UserID: noteOwner, DailyNote: "not yours"}, nil
		},
	}
	config := &Config{Notes: store}

	req := newAuthedRequest(t, http.MethodGet, "/notes/"+noteID.String(), requestingUser, nil)
	req = withNoteID(req, noteID)
	rec := httptest.NewRecorder()

	config.ReadNote(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected status %d when reading someone else's note, got %d", http.StatusForbidden, rec.Code)
	}
}

func TestReadNote_InvalidNoteID(t *testing.T) {
	userID := uuid.New()
	config := &Config{Notes: &fakeNoteStore{}} // store should never be touched

	req := newAuthedRequest(t, http.MethodGet, "/notes/not-a-uuid", userID, nil)
	req.SetPathValue("note_id", "not-a-uuid")
	rec := httptest.NewRecorder()

	config.ReadNote(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for a malformed note id, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestAddNote_Success(t *testing.T) {
	userID := uuid.New()
	noteID := uuid.New()

	store := &fakeNoteStore{
		createNoteFunc: func(ctx context.Context, arg database.CreateNoteParams) (database.Note, error) {
			if arg.UserID != userID {
				t.Errorf("expected note to be created for user %v, got %v", userID, arg.UserID)
			}
			if arg.DailyNote != "gym then leetcode" {
				t.Errorf("expected daily_note %q, got %q", "gym then leetcode", arg.DailyNote)
			}
			return database.Note{
				NoteID:    noteID,
				UserID:    arg.UserID,
				DailyNote: arg.DailyNote,
			}, nil
		},
	}
	config := &Config{Notes: store}

	req := newAuthedRequest(t, http.MethodPost, "/notes", userID, Note{
		DailyNote: "gym then leetcode",
	})
	rec := httptest.NewRecorder()

	config.AddNote(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusOK, rec.Code, rec.Body.String())
	}

	var got Note
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	if got.NoteId != noteID {
		t.Errorf("expected note id %v, got %v", noteID, got.NoteId)
	}
}

func TestAddNote_InvalidJSON(t *testing.T) {
	userID := uuid.New()
	config := &Config{Notes: &fakeNoteStore{}}

	token, err := auth.MakeJWT(userID, os.Getenv("JWT_SECRET"), time.Hour)
	if err != nil {
		t.Fatalf("test setup: failed to create JWT: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader([]byte("not json")))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	config.AddNote(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for malformed JSON, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestAddNote_DatabaseError(t *testing.T) {
	userID := uuid.New()
	store := &fakeNoteStore{
		createNoteFunc: func(ctx context.Context, arg database.CreateNoteParams) (database.Note, error) {
			return database.Note{}, errors.New("insert failed")
		},
	}
	config := &Config{Notes: store}

	req := newAuthedRequest(t, http.MethodPost, "/notes", userID, Note{DailyNote: "test"})
	rec := httptest.NewRecorder()

	config.AddNote(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestDeleteNote_Success(t *testing.T) {
	userID := uuid.New()
	noteID := uuid.New()

	store := &fakeNoteStore{
		readNoteFunc: func(ctx context.Context, id uuid.UUID) (database.Note, error) {
			return database.Note{NoteID: id, UserID: userID}, nil
		},
		deleteNoteFunc: func(ctx context.Context, id uuid.UUID) error {
			if id != noteID {
				t.Errorf("expected to delete note %v, got %v", noteID, id)
			}
			return nil
		},
	}
	config := &Config{Notes: store}

	req := newAuthedRequest(t, http.MethodDelete, "/notes/"+noteID.String(), userID, nil)
	req = withNoteID(req, noteID)
	rec := httptest.NewRecorder()

	config.DeleteNote(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d (body: %s)", http.StatusOK, rec.Code, rec.Body.String())
	}
}

func TestDeleteNote_ForbiddenWhenNoteBelongsToSomeoneElse(t *testing.T) {
	requestingUser := uuid.New()
	noteOwner := uuid.New()
	noteID := uuid.New()

	store := &fakeNoteStore{
		readNoteFunc: func(ctx context.Context, id uuid.UUID) (database.Note, error) {
			return database.Note{NoteID: id, UserID: noteOwner}, nil
		},
	}
	config := &Config{Notes: store}

	req := newAuthedRequest(t, http.MethodDelete, "/notes/"+noteID.String(), requestingUser, nil)
	req = withNoteID(req, noteID)
	rec := httptest.NewRecorder()

	config.DeleteNote(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
}

func TestDeleteNote_NotFound(t *testing.T) {
	userID := uuid.New()
	noteID := uuid.New()

	store := &fakeNoteStore{
		readNoteFunc: func(ctx context.Context, id uuid.UUID) (database.Note, error) {
			return database.Note{}, sql.ErrNoRows
		},
	}
	config := &Config{Notes: store}

	req := newAuthedRequest(t, http.MethodDelete, "/notes/"+noteID.String(), userID, nil)
	req = withNoteID(req, noteID)
	rec := httptest.NewRecorder()

	config.DeleteNote(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestUpdateNote_Success(t *testing.T) {
	userID := uuid.New()
	noteID := uuid.New()

	store := &fakeNoteStore{
		readNoteFunc: func(ctx context.Context, id uuid.UUID) (database.Note, error) {
			return database.Note{NoteID: id, UserID: userID, DailyNote: "old note"}, nil
		},
		updateNoteFunc: func(ctx context.Context, arg database.UpdateNoteParams) (database.Note, error) {
			if arg.NoteID != noteID {
				t.Errorf("expected to update note %v, got %v", noteID, arg.NoteID)
			}
			if arg.DailyNote != "updated note" {
				t.Errorf("expected updated text %q, got %q", "updated note", arg.DailyNote)
			}
			return database.Note{NoteID: noteID, UserID: userID, DailyNote: arg.DailyNote}, nil
		},
	}
	config := &Config{Notes: store}

	req := newAuthedRequest(t, http.MethodPut, "/notes/"+noteID.String(), userID, Note{
		DailyNote: "updated note",
	})
	req = withNoteID(req, noteID)
	rec := httptest.NewRecorder()

	config.UpdateNote(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusOK, rec.Code, rec.Body.String())
	}

	var got Note
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	if got.DailyNote != "updated note" {
		t.Errorf("expected response daily_note %q, got %q", "updated note", got.DailyNote)
	}
}

func TestUpdateNote_ForbiddenWhenNoteBelongsToSomeoneElse(t *testing.T) {
	requestingUser := uuid.New()
	noteOwner := uuid.New()
	noteID := uuid.New()

	store := &fakeNoteStore{
		readNoteFunc: func(ctx context.Context, id uuid.UUID) (database.Note, error) {
			return database.Note{NoteID: id, UserID: noteOwner}, nil
		},
	}
	config := &Config{Notes: store}

	req := newAuthedRequest(t, http.MethodPut, "/notes/"+noteID.String(), requestingUser, Note{
		DailyNote: "trying to edit someone else's note",
	})
	req = withNoteID(req, noteID)
	rec := httptest.NewRecorder()

	config.UpdateNote(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
}

func TestUpdateNote_NotFound(t *testing.T) {
	userID := uuid.New()
	noteID := uuid.New()

	store := &fakeNoteStore{
		readNoteFunc: func(ctx context.Context, id uuid.UUID) (database.Note, error) {
			return database.Note{}, sql.ErrNoRows
		},
	}
	config := &Config{Notes: store}

	req := newAuthedRequest(t, http.MethodPut, "/notes/"+noteID.String(), userID, Note{
		DailyNote: "doesn't matter, note doesn't exist",
	})
	req = withNoteID(req, noteID)
	rec := httptest.NewRecorder()

	config.UpdateNote(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
