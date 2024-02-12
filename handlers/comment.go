package handlers

import (
	"database/sql"
	"fmt"
	"forum/models"
	"forum/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Comment(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		utils.Handle500Error(w)
		return
	}
	defer db.Close()
	sessionData, isLoggedIn := getSessionData(w, r, db)

	if r.URL.Path != "/comment" {

		utils.Handle404Error(w)
		return
	}

	if isLoggedIn && isValidSession(sessionData) {
		var p models.Comment
		Comments(w, r, &p, sessionData)
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func Comments(w http.ResponseWriter, r *http.Request, p *models.Comment, sessionData models.SessionData) {
	if r.Method == "POST" {
		if r.FormValue("comment") == "Comment" {
			idpost, _ := strconv.Atoi(r.FormValue("id"))
			content_comment := r.PostFormValue("content_comment")
			if len(strings.Fields(content_comment)) == 0 {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			CommentPost(sessionData.UserID, idpost, content_comment)
			db, err := sql.Open("sqlite3", "./db/forum.db")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer db.Close()
			posted, err := db.Prepare("SELECT id_user, id_post, content_comment FROM comment WHERE id_post = ? ")
			if err != nil {
				fmt.Println(err)
				fmt.Println("nono")
				return
			}
			defer posted.Close()

			rows, _ := posted.Query(idpost)
			var (
				IdUser         int
				IdPost         int
				ContentComment string
			)
			for rows.Next() {
				if err := rows.Scan(&IdUser, &IdPost, &ContentComment); err != nil {
					fmt.Println(err)
					fmt.Println("error")
					return
				}
			}
			p.Id_user = IdUser
			p.Id_post = IdPost
			p.Content_comment = ContentComment
			postID := r.FormValue("id")

			// InfoUser := Info{
			// 	IdUserConnect:   sessionData.UserID,
			// 	NameUserConnect: sessionData.Username,
			// 	State:           IsConnected,
			// }
			UserConnected.Id = sessionData.UserID
			// Exec := map[string]interface{}{"Infos": InfoUser}
			// utils.RenderTemplate(w, "header", Exec)

			http.Redirect(w, r, "/#p"+postID, http.StatusFound)
		}
	} else {
		utils.Handle405Error(w)
		return
	}
}

func CommentPost(id_user int, id_post int, content_comment string) bool {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		fmt.Println("Erreur se produite")

		log.Fatal(err)
	}
	defer db.Close() // Assurez-vous que db.Close() est différé ici.
	users, err := db.Prepare("INSERT INTO comment (id_user, id_post, content_comment) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Print("usertable creation err!: ")
		fmt.Println(err)
		return false
	}

	users.Exec(id_user, id_post, content_comment)
	db.Close()
	return true

}
