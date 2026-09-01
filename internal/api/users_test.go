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

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-tests")
	os.Exit(m.Run())
}

type fakeUserStore struct {
	registerUserFunc       func(ctx context.Context, arg database.RegisterUserParams) (database.User, error)
	searchUserByEmailFunc  func(ctx context.Context, email string) (database.User, error)
	createRefreshTokenFunc func(ctx context.Context, arg database.CreateRefreshTokenParams) (database.RefreshToken, error)
	getRefreshTokenFunc    func(ctx context.Context, token string) (database.RefreshToken, error)
	updateRefreshTokenFunc func(ctx context.Context, token string) error
}

func (f *fakeUserStore) RegisterUser(ctx context.Context, arg database.RegisterUserParams) (database.User, error) {
	return f.registerUserFunc(ctx, arg)
}

func (f *fakeUserStore) SearchUserByEmail(ctx context.Context, email string) (database.User, error) {
	return f.searchUserByEmailFunc(ctx, email)
}

func (f *fakeUserStore) CreateRefreshToken(ctx context.Context, arg database.CreateRefreshTokenParams) (database.RefreshToken, error) {
	return f.createRefreshTokenFunc(ctx, arg)
}
func (f *fakeUserStore) GetRefreshToken(ctx context.Context, token string) (database.RefreshToken, error) {
	return f.getRefreshTokenFunc(ctx, token)
}
func (f *fakeUserStore) UpdateRefreshToken(ctx context.Context, token string) error {
	return f.updateRefreshTokenFunc(ctx, token)
}

func newJSONRequest(t *testing.T, method, path string, body any) *http.Request {
	t.Helper()

	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestValidateCredentials(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{"valid credentials", "user@example.com", "password123", false},
		{"empty email", "", "password123", true},
		{"empty password", "user@example.com", "", true},
		{"both empty", "", "", true},
		{"password too short", "user@example.com", "abc12", true},
		{"password exactly 6 characters is allowed", "user@example.com", "abc123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCredentials(tt.email, tt.password)

			if tt.wantErr && err == nil {
				t.Errorf("expected an error for email=%q password=%q, got nil", tt.email, tt.password)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error for email=%q password=%q, got: %v", tt.email, tt.password, err)
			}
		})
	}
}

func TestRegisterUser_Success(t *testing.T) {
	wantID := uuid.New()
	wantEmail := "new@example.com"

	store := &fakeUserStore{
		registerUserFunc: func(ctx context.Context, arg database.RegisterUserParams) (database.User, error) {
			if arg.Email != wantEmail {
				t.Errorf("expected email %q to reach the DB, got %q", wantEmail, arg.Email)
			}
			if arg.PasswordHash == "" || arg.PasswordHash == "supersecret" {
				t.Errorf("expected the password to be hashed before reaching the DB, got %q", arg.PasswordHash)
			}
			return database.User{ID: wantID, Email: wantEmail}, nil
		},
	}
	cfg := &Config{Users: store}

	req := newJSONRequest(t, http.MethodPost, "/api/users", UserRequest{
		Email:    wantEmail,
		Password: "supersecret",
	})
	w := httptest.NewRecorder()

	cfg.RegisterUser(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusCreated, w.Code, w.Body.String())
	}

	var got UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	if got.Email != wantEmail {
		t.Errorf("expected response email %q, got %q", wantEmail, got.Email)
	}
	if got.Id != wantID {
		t.Errorf("expected response id %v, got %v", wantID, got.Id)
	}
}

func TestRegisterUser_InvalidJSON(t *testing.T) {
	cfg := &Config{Users: &fakeUserStore{}} // DB should never be touched
	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader([]byte("not valid json")))
	w := httptest.NewRecorder()

	cfg.RegisterUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for malformed JSON, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRegisterUser_InvalidCredentials(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
	}{
		{"empty email", "", "password123"},
		{"empty password", "user@example.com", ""},
		{"password too short", "user@example.com", "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Users: &fakeUserStore{}} // DB should never be touched
			req := newJSONRequest(t, http.MethodPost, "/api/users", UserRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			w := httptest.NewRecorder()

			cfg.RegisterUser(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
			}
		})
	}
}

func TestRegisterUser_DatabaseError(t *testing.T) {
	store := &fakeUserStore{
		registerUserFunc: func(ctx context.Context, arg database.RegisterUserParams) (database.User, error) {
			return database.User{}, errors.New("duplicate email")
		},
	}
	cfg := &Config{Users: store}

	req := newJSONRequest(t, http.MethodPost, "/api/users", UserRequest{
		Email:    "taken@example.com",
		Password: "supersecret",
	})
	w := httptest.NewRecorder()

	cfg.RegisterUser(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d when the DB call fails, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestLoginUser_Success(t *testing.T) {
	const password = "correct-password"
	hash, err := auth.GenerateHash(password)
	if err != nil {
		t.Fatalf("test setup: failed to hash password: %v", err)
	}

	userID := uuid.New()

	store := &fakeUserStore{
		searchUserByEmailFunc: func(ctx context.Context, email string) (database.User, error) {
			return database.User{
				ID:           userID,
				Email:        email,
				PasswordHash: hash,
				CreatedAt:    time.Now(),
			}, nil
		},
		createRefreshTokenFunc: func(ctx context.Context, arg database.CreateRefreshTokenParams) (database.RefreshToken, error) {
			if arg.Userid != userID {
				t.Errorf("expected refresh token to be created for user %v, got %v", userID, arg.Userid)
			}
			return database.RefreshToken{Token: arg.Token}, nil
		},
	}
	cfg := &Config{Users: store}

	req := newJSONRequest(t, http.MethodPost, "/api/login", UserRequest{
		Email:    "user@example.com",
		Password: password,
	})
	w := httptest.NewRecorder()

	cfg.LoginUser(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusOK, w.Code, w.Body.String())
	}

	var got UserLoginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	if got.Id != userID {
		t.Errorf("expected response id %v, got %v", userID, got.Id)
	}
	if got.AccessToken == "" {
		t.Errorf("expected a non-empty access token")
	}
	if got.RefreshToken == "" {
		t.Errorf("expected a non-empty refresh token")
	}
}

func TestLoginUser_WrongPassword(t *testing.T) {
	hash, err := auth.GenerateHash("the-real-password")
	if err != nil {
		t.Fatalf("test setup: failed to hash password: %v", err)
	}

	store := &fakeUserStore{
		searchUserByEmailFunc: func(ctx context.Context, email string) (database.User, error) {
			return database.User{Email: email, PasswordHash: hash}, nil
		},
	}
	cfg := &Config{Users: store}

	req := newJSONRequest(t, http.MethodPost, "/api/login", UserRequest{
		Email:    "user@example.com",
		Password: "wrong-password",
	})
	w := httptest.NewRecorder()

	cfg.LoginUser(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d for wrong password, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestLoginUser_UserNotFound(t *testing.T) {
	store := &fakeUserStore{
		searchUserByEmailFunc: func(ctx context.Context, email string) (database.User, error) {
			return database.User{}, sql.ErrNoRows
		},
	}
	cfg := &Config{Users: store}

	req := newJSONRequest(t, http.MethodPost, "/api/login", UserRequest{
		Email:    "missing@example.com",
		Password: "whatever1",
	})
	w := httptest.NewRecorder()

	cfg.LoginUser(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d for unknown email, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestLoginUser_DatabaseError(t *testing.T) {
	store := &fakeUserStore{
		searchUserByEmailFunc: func(ctx context.Context, email string) (database.User, error) {
			return database.User{}, errors.New("connection refused")
		},
	}
	cfg := &Config{Users: store}

	req := newJSONRequest(t, http.MethodPost, "/api/login", UserRequest{
		Email:    "user@example.com",
		Password: "whatever1",
	})
	w := httptest.NewRecorder()

	cfg.LoginUser(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d when the lookup fails, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestLoginUser_InvalidJSON(t *testing.T) {
	cfg := &Config{Users: &fakeUserStore{}}
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader([]byte("not valid json")))
	w := httptest.NewRecorder()

	cfg.LoginUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for malformed JSON, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestLoginUser_InvalidCredentialsFormat(t *testing.T) {
	cfg := &Config{Users: &fakeUserStore{}}
	req := newJSONRequest(t, http.MethodPost, "/api/login", UserRequest{
		Email:    "",
		Password: "whatever1",
	})
	w := httptest.NewRecorder()

	cfg.LoginUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestLoginUser_RefreshTokenCreationFails(t *testing.T) {
	hash, err := auth.GenerateHash("correct-password")
	if err != nil {
		t.Fatalf("test setup: failed to hash password: %v", err)
	}

	store := &fakeUserStore{
		searchUserByEmailFunc: func(ctx context.Context, email string) (database.User, error) {
			return database.User{ID: uuid.New(), Email: email, PasswordHash: hash}, nil
		},
		createRefreshTokenFunc: func(ctx context.Context, arg database.CreateRefreshTokenParams) (database.RefreshToken, error) {
			return database.RefreshToken{}, errors.New("insert failed")
		},
	}
	cfg := &Config{Users: store}

	req := newJSONRequest(t, http.MethodPost, "/api/login", UserRequest{
		Email:    "user@example.com",
		Password: "correct-password",
	})
	w := httptest.NewRecorder()

	cfg.LoginUser(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d when refresh token creation fails, got %d", http.StatusInternalServerError, w.Code)
	}
}
