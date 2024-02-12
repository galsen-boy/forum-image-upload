package handlers

import (
	"database/sql"
	"forum/utils"
	"net/http"
	"strconv"
)

func LikeComment(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/like_comment" {

		utils.Handle404Error(w)
		return
	}
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		utils.Handle500Error(w)
		return
	}
	defer db.Close()
	sessionData, isLoggedIn := getSessionData(w, r, db)

	if isLoggedIn && isValidSession(sessionData) {
		if r.Method == "POST" {
			if r.FormValue("like_comment") == "thumb_up" {
				id_comment, _ := strconv.Atoi(r.FormValue("idcomment"))
				var islike bool

				err = db.QueryRow("SELECT isLike FROM likes_comment WHERE id_user = ? AND id_comment = ? ", sessionData.UserID, id_comment).Scan(&islike)
				if err != nil {
					islike = true
					_, err = db.Exec("INSERT INTO likes_comment (id_comment, id_user, isLike) VALUES (?, ? , ?);", id_comment, sessionData.UserID, islike)
					if err != nil {
						utils.Handle500Error(w)
						return
					}
				} else {
					if !islike {
						_, err = db.Exec("UPDATE likes_comment SET isLike = true WHERE id_user = $1 AND id_comment = $2", sessionData.UserID, id_comment)
						if err != nil {
							utils.Handle500Error(w)
							return
						}
					} else {
						_, _ = db.Exec("DELETE FROM likes_comment WHERE id_user = ? AND id_comment = ? AND isLike = true ", sessionData.UserID, id_comment)
					}
					// _, err = db.Exec("UPDATE like_post SET isLike = false WHERE id_user = $1 AND id_post = $2", sessionData.UserID, id_post)

				}

				commentID, _ := strconv.Atoi(r.FormValue("idcomment"))
				var pID int
				err = db.QueryRow("SELECT id_post FROM comment WHERE id = ?", commentID).Scan(&pID)
				idPost := strconv.Itoa(pID)
				http.Redirect(w, r, "/#p"+idPost, http.StatusSeeOther)

			} else if r.FormValue("dislike_comment") == "thumb_down" {
				id_comment, _ := strconv.Atoi(r.FormValue("idcomment"))

				var islike bool

				err = db.QueryRow("SELECT isLike FROM likes_comment WHERE id_user = ? AND id_comment = ? ", sessionData.UserID, id_comment).Scan(&islike)
				if err != nil {
					islike = false
					_, err = db.Exec("INSERT INTO likes_comment (id_comment, id_user, isLike) VALUES (?, ? , ?);", id_comment, sessionData.UserID, islike)
					if err != nil {
						utils.Handle500Error(w)
						return
					}
				} else {
					if islike {
						_, err = db.Exec("UPDATE likes_comment SET isLike = false WHERE id_user = $1 AND id_comment = $2", sessionData.UserID, id_comment)
						if err != nil {
							utils.Handle500Error(w)
							return
						}
					} else {
						_, _ = db.Exec("DELETE FROM likes_comment WHERE id_user = ? AND id_comment = ? AND isLike = false ", sessionData.UserID, id_comment)
						if err != nil {
							utils.Handle500Error(w)
							return
						}
					}
				}

				commentID, _ := strconv.Atoi(r.FormValue("idcomment"))
				var pID int
				err = db.QueryRow("SELECT id_post FROM comment WHERE id = ?", commentID).Scan(&pID)
				idPost := strconv.Itoa(pID)
				http.Redirect(w, r, "/#p"+idPost, http.StatusSeeOther)
			}

		}else {
			utils.Handle405Error(w)
		}
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
