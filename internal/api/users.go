package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"myapp/internal/auth"
	"myapp/internal/database"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Id           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

type UserResponse struct {
	Id    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

func ValidateCredentials(Email string, Password string) error {
	if Email == "" || Password == "" {
		return errors.New("Email and Password Cant be empty")
	}

	if len(Password) < 6 {
		return errors.New("Password Should be 6 Characters Long")
	}
	return nil
}

func (config *Config) RegisterUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	UserReq := UserRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&UserReq); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	if err := ValidateCredentials(UserReq.Email, UserReq.Password); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	PassHash, err := auth.GenerateHash(UserReq.Password)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}
	DBuser, err := config.DB.RegisterUser(r.Context(), database.RegisterUserParams{
		Email:        UserReq.Email,
		PasswordHash: PassHash,
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	UserResp := UserResponse{}

	UserResp.Id = DBuser.ID
	UserResp.Email = DBuser.Email

	RespondWithJson(w, http.StatusOK, UserResp)
}

func (config *Config) LoginUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	UserReq := UserRequest{}

	if err := json.NewDecoder(r.Body).Decode(&UserReq); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	if err := ValidateCredentials(UserReq.Email, UserReq.Password); err != nil {
		RespondWithJson(w, http.StatusBadRequest, err.Error())
		return
	}

	DbUser, err := config.DB.SearchUserByEmail(r.Context(), UserReq.Email)

	if err == sql.ErrNoRows {
		RespondWithError(w, http.StatusUnauthorized, "Invalid Email or Password")
		return
	}
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	PassHash := DbUser.PasswordHash
	IsCorrect, err := auth.CompareHash(UserReq.Password, PassHash)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	if !IsCorrect {
		RespondWithError(w, http.StatusUnauthorized, "Invalid Email or Password")
		return
	}

	expiresIn := time.Hour

	token, err := auth.MakeJWT(DbUser.ID, os.Getenv("JWT_SECRET"), expiresIn)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "server error")
		return
	}

	RefreshToken := auth.MakeRefreshToken()
	_, err = config.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     RefreshToken,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		RevokedAt: sql.NullTime{},
		Userid:    DbUser.ID,
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	UserLogResp := UserLoginResponse{
		Id:           DbUser.ID,
		Email:        DbUser.Email,
		CreatedAt:    DbUser.CreatedAt,
		AccessToken:  token,
		RefreshToken: RefreshToken,
	}

	RespondWithJson(w, http.StatusOK, UserLogResp)
}
