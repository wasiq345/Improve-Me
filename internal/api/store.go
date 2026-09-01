package api

import (
	"context"
	"myapp/internal/database"

	"github.com/google/uuid"
)

type UserStore interface {
	RegisterUser(ctx context.Context, arg database.RegisterUserParams) (database.User, error)
	SearchUserByEmail(ctx context.Context, email string) (database.User, error)
	CreateRefreshToken(ctx context.Context, arg database.CreateRefreshTokenParams) (database.RefreshToken, error)
	GetRefreshToken(ctx context.Context, token string) (database.RefreshToken, error)
	UpdateRefreshToken(ctx context.Context, token string) error
}

type NoteStore interface {
	CreateNote(ctx context.Context, arg database.CreateNoteParams) (database.Note, error)
	DeleteNote(ctx context.Context, noteId uuid.UUID) error
	GetAllNotes(ctx context.Context, userId uuid.UUID) ([]database.Note, error)
	ReadNote(ctx context.Context, noteId uuid.UUID) (database.Note, error)
	UpdateNote(ctx context.Context, arg database.UpdateNoteParams) (database.Note, error)
}
