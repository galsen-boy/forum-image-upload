package handlers

import (
	"database/sql"
	"forum/utils"
	"net/http"
	"strconv"
)

func Liked(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/like" {

		utils.Handle404Error(w)
		return
	}
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		//
		// Handle500Error(w)
		return
	}
	defer db.Close()
	sessionData, isLoggedIn := getSessionData(w, r, db)
	postID := r.FormValue("idpost")
	if isLoggedIn && isValidSession(sessionData) {
		if r.Method == "POST" {
			if r.FormValue("like") == "thumb_up" {
				id_post, _ := strconv.Atoi(r.FormValue("idpost"))

				var islike bool

				err = db.QueryRow("SELECT isLike FROM likes_post WHERE id_user = ? AND id_post = ? ", sessionData.UserID, id_post).Scan(&islike)
				if err != nil {
					islike = true
					_, err = db.Exec("INSERT INTO likes_post (id_post, id_user, isLike) VALUES (?, ? , ?);", id_post, sessionData.UserID, islike)
					if err != nil {
						utils.Handle500Error(w)
						return
					}
				} else {
					if !islike {
						_, err = db.Exec("UPDATE likes_post SET isLike = true WHERE id_user = $1 AND id_post = $2", sessionData.UserID, id_post)
						if err != nil {
							utils.Handle500Error(w)
							return
						}
					} else {
						_, _ = db.Exec("DELETE FROM likes_post WHERE id_user = ? AND id_post = ? AND isLike = true ", sessionData.UserID, id_post)
					}
					// _, err = db.Exec("UPDATE like_post SET isLike = false WHERE id_user = $1 AND id_post = $2", sessionData.UserID, id_post)

				}

				http.Redirect(w, r, "/#p"+postID, http.StatusSeeOther)

			} else if r.FormValue("dislike") == "thumb_down" {
				id_post, _ := strconv.Atoi(r.FormValue("idpost"))
				var islike bool

				err = db.QueryRow("SELECT isLike FROM likes_post WHERE id_user = ? AND id_post = ? ", sessionData.UserID, id_post).Scan(&islike)
				if err != nil {
					islike = false
					_, err = db.Exec("INSERT INTO likes_post (id_post, id_user, isLike) VALUES (?, ? , ?);", id_post, sessionData.UserID, islike)
					if err != nil {
						utils.Handle500Error(w)
						return
					}
				} else {
					if islike {
						_, err = db.Exec("UPDATE likes_post SET isLike = false WHERE id_user = $1 AND id_post = $2", sessionData.UserID, id_post)
						if err != nil {
							utils.Handle500Error(w)
							return
						}
					} else {
						_, _ = db.Exec("DELETE FROM likes_post WHERE id_user = ? AND id_post = ? AND isLike = false ", sessionData.UserID, id_post)
						if err != nil {
							utils.Handle500Error(w)
							return
						}
					}
				}
				http.Redirect(w, r, "/#p"+postID, http.StatusSeeOther)

			}

		} else {
			utils.Handle405Error(w)
		}
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
