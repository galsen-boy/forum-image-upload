package handlers

import (
	"database/sql"
	"fmt"
	"forum/models"
	"forum/utils"
	"log"
	"net/http"
)

func Profil(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/profil" {

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

	if isLoggedIn && isValidSession(sessionData) {
		// InfoUser := Info{
		// 	NameUserConnect: sessionData.Username,
		// 	State:           IsConnected,
		// }
		UserConnected.Id = sessionData.UserID
		// Exec := map[string]interface{}{"Infos": InfoUser}
		fetchAndPostsProfil(w, "profil", db, sessionData)

		// utils.RenderTemplate(w, "profil", Exec)
		return
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func fetchAndPostsProfil(w http.ResponseWriter, template string, db *sql.DB, sessionData models.SessionData) {

	row, err := db.Prepare("SELECT * FROM posts WHERE id_user = ?")
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		utils.Handle400Error(w)
		return
	}
	rows, _ := row.Query(sessionData.UserID)

	defer rows.Close()

	var Posts []models.Post

	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Media, &post.Date, &post.UserID)
		post.Like = utils.GetLikeByPostID(post.ID)
		post.Dislike = utils.GetDislikeByPOstID(post.ID)
		if err != nil {
			fmt.Println(err)

			utils.Handle400Error(w)
			return
			// return
		}

		db_id_category, err := db.Prepare("SELECT id_category FROM belong where id_post = ?")
		if err != nil {

			utils.Handle400Error(w)
			return
		}
		id_category, _ := db_id_category.Query(post.ID)
		var id_categ int
		for id_category.Next() {

			if err := id_category.Scan(&id_categ); err != nil {
				log.Fatal(err)
			}
			db_name_category, err := db.Prepare("SELECT name_category FROM category where id = ?")
			if err != nil {

				utils.Handle400Error(w)
				return
			}
			name_category, _ := db_name_category.Query(id_categ)
			var name_categ string
			for name_category.Next() {

				if err := name_category.Scan(&name_categ); err != nil {
					log.Fatal(err)
				}
				post.Categorie = append(post.Categorie, name_categ)
			}
		}

		name_user, err := db.Prepare("SELECT name_user FROM users WHERE id = ?")
		if err != nil {

			utils.Handle400Error(w)
			return
		}
		name, err := name_user.Query(post.UserID)

		var name_User string
		for name.Next() {
			if err := name.Scan(&name_User); err != nil {
				log.Fatal(err)
			}
			post.NameUser = name_User
		}

		if err != nil {
			utils.Handle500Error(w)

			return
		}

		content_comment, err := db.Prepare("SELECT id, content_comment, id_user FROM comment WHERE id_post = ?")
		if err != nil {

			utils.Handle400Error(w)
			return
		}
		commentRows, err := content_comment.Query(post.ID)

		var comments []models.CommentWithUser

		for commentRows.Next() {
			var contentComment string
			var id_comment, idUser int

			if err := commentRows.Scan(&id_comment, &contentComment, &idUser); err != nil {
				log.Fatal(err)
			}

			name_user, err := db.Prepare("SELECT name_user FROM users WHERE id = ?")
			if err != nil {

				utils.Handle400Error(w)
				return
			}
			name, err := name_user.Query(idUser)

			var name_User string
			for name.Next() {
				if err := name.Scan(&name_User); err != nil {
					log.Fatal(err)
				}
			}

			if err != nil {
				utils.Handle500Error(w)

				return
			}

			commentWithUser := models.CommentWithUser{
				Content:    contentComment,
				UserID:     idUser,
				Name_User:  name_User,
				Id_comment: id_comment,
				Like:       utils.GetLikeByCommentID(id_comment, db),
				Dislike:    utils.GetDislikeByCommentID(id_comment, db),
			}

			comments = append([]models.CommentWithUser{commentWithUser}, comments...)

			// comments = append(comments, commentWithUser)
		}

		if err != nil {
			utils.Handle500Error(w)

			return
		}

		post.Date = utils.FormatDate(post.Date)
		post.Comments = comments // Attribuez la liste des commentaires à la publication
		Posts = append([]models.Post{post}, Posts...)

		// Posts = append(Posts, post)
	}

	if err := rows.Err(); err != nil {

		utils.Handle500Error(w)
		return
	}
	var mail_user string
	_ = db.QueryRow("SELECT mail_user FROM users WHERE id = ?", sessionData.UserID).Scan(&mail_user)
	InfoUser := Info{
		IdUserConnect:   sessionData.UserID,
		NameUserConnect: sessionData.Username,
		State:           true,
		Mail:            mail_user,
	}
	UserConnected.Id = sessionData.UserID
	Exec := map[string]interface{}{"Posts": Posts, "Infos": InfoUser}
	utils.RenderTemplate(w, template, Exec)
}
