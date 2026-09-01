package api

import (
	"database/sql"
	"myapp/internal/auth"
	"net/http"
	"os"
	"time"
)

func (apicfg *Config) Revoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid Format")
		return
	}

	err = apicfg.Users.UpdateRefreshToken(r.Context(), refreshToken)

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Server Error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (apicfg *Config) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid format")
		return
	}

	Rtoken, err := apicfg.Users.GetRefreshToken(r.Context(), refreshToken)

	if err == sql.ErrNoRows {
		RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	if Rtoken.ExpiresAt.Before(time.Now()) {
		RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	if Rtoken.RevokedAt.Valid {
		RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "DataBase error")
	}
	jwt, err := auth.MakeJWT(Rtoken.Userid, os.Getenv("JWT_SECRET"), time.Hour)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	RespondWithJson(w, http.StatusOK, map[string]string{"token": jwt})
}
