package handlers

import (
	"database/sql"
	"forum/utils"
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/logout" {

		utils.Handle404Error(w)
		return
	}
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	_, err = db.Exec("DELETE FROM sessions WHERE session_id = ?", cookie.Value)
	if err != nil {
		// Gérer l'erreur
		utils.Handle400Error(w)
		return
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/login", http.StatusFound)
}
