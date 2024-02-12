package handlers

import (
	"database/sql"
	"fmt"
	"forum/models"
	"forum/utils"
	"log"
	"net/http"
)

var UserConnected UserConnec

var IsConnected bool

func Index(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {

		utils.Handle500Error(w)
		return
	}
	defer db.Close()
	sessionData, isLoggedIn := getSessionData(w, r, db)
	if isLoggedIn && isValidSession(sessionData) {
		// fmt.Fprintf(w, "Bienvenue, %s !", sessionData.Username)
		IsConnected = true

	} else {
		IsConnected = false
	}
	if r.URL.Path != "/" {

		utils.Handle404Error(w)
		return
	}
	fetchAndRenderPosts(w, "index", db, sessionData)
	// fetchAndRenderPosts(w, "index_no_user", sessionData)
}

func fetchAndRenderPosts(w http.ResponseWriter, template string, db *sql.DB, sessionData models.SessionData) {

	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {

		utils.Handle500Error(w)
		return
	}
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

		err = db.QueryRow("SELECT mail_user FROM users WHERE id = ?", post.UserID).Scan(&post.MailUser)
		if err != nil {
			utils.Handle500Error(w)
			return
		}
		Posts = append([]models.Post{post}, Posts...)

		// Posts = append(Posts, post)
	}

	if err := rows.Err(); err != nil {

		utils.Handle500Error(w)
		return
	}

	InfoUser := Info{
		IdUserConnect:   sessionData.UserID,
		NameUserConnect: sessionData.Username,
		State:           IsConnected,
	}
	UserConnected.Id = sessionData.UserID
	Exec := map[string]interface{}{"Posts": Posts, "Infos": InfoUser}
	utils.RenderTemplate(w, template, Exec)
}

type Info struct {
	IdUserConnect   int
	NameUserConnect string
	State           bool
	Mail            string
}
