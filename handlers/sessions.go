package handlers

import (
	"database/sql"
	"encoding/json"
	"forum/models"
	"forum/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type SessionData struct {
	UserID      int
	Username    string
	IsLoggedIn  bool
	CreatedTime time.Time
}

func createSession(w http.ResponseWriter, r *http.Request, data models.SessionData, db *sql.DB) {
	// Calculez l'heure d'expiration (5 heures à partir de maintenant)
	expiration := time.Now().Add(5 * time.Hour)

	data.CreatedTime = time.Now()
	token := uuid.NewString()

	dataJSON, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = db.Exec("INSERT INTO sessions (session_id, data, expiration, id_user, name_user) VALUES (?, ?, ?, ?, ?)", token, string(dataJSON), expiration, data.UserID, data.Username)

	cookie := http.Cookie{
		Name:     "session",
		Value:    token,
		Expires:  expiration,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}

func getSessionData(w http.ResponseWriter, r *http.Request, db *sql.DB) (models.SessionData, bool) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return models.SessionData{}, false
	}

	var dataJSON string
	var expiration time.Time
	err = db.QueryRow("SELECT data, expiration FROM sessions WHERE session_id = ?", cookie.Value).Scan(&dataJSON, &expiration)
	if err != nil {
		// utils.Handle400Error(w)

		return models.SessionData{}, false
	}

	if time.Now().After(expiration) {
		_, err := db.Exec("DELETE FROM sessions WHERE session_id = ?", cookie.Value)
		if err != nil {
			utils.Handle400Error(w)
		}

		cookie.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, cookie)
		return models.SessionData{}, false
	}

	var sessionData models.SessionData
	err = json.Unmarshal([]byte(dataJSON), &sessionData)
	if err != nil {
		return models.SessionData{}, false
	}

	return sessionData, true
}

func isValidSession(sessionData models.SessionData) bool {
	return time.Now().Before(sessionData.CreatedTime.Add(5 * time.Hour))
}
