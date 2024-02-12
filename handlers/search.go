package handlers

import (
	"database/sql"
	"fmt"
	"forum/models"
	"forum/utils"
	"log"
	"net/http"
	"strings"
)

var toAddPost bool = false

func Search(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		value := r.PostFormValue("search")
		fetchPostSearch(w, r, strings.ToLower(value))

	} else {
		utils.Handle405Error(w)
	}
}

func fetchPostSearch(w http.ResponseWriter, r *http.Request, valueToSearch string) {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {

		utils.Handle500Error(w)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title_post, content_post, media_post, date_post, id_user FROM posts")
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

				if strings.Contains(strings.ToLower(name_categ), valueToSearch) {
					toAddPost = true
				}
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
			if strings.Contains(strings.ToLower(name_User), valueToSearch) {
				toAddPost = true
			}
		}

		if err != nil {
			utils.Handle500Error(w)

			return
		}

		content_comment, err := db.Prepare("SELECT content_comment, id_user FROM comment WHERE id_post = ?")
		if err != nil {

			utils.Handle400Error(w)
			return
		}
		commentRows, err := content_comment.Query(post.ID)

		var comments []models.CommentWithUser

		for commentRows.Next() {
			var contentComment string
			var id_comment, idUser int

			if err := commentRows.Scan(&contentComment, &idUser); err != nil {
				log.Fatal(err)
			}

			if strings.Contains(strings.ToLower(contentComment), valueToSearch) {
				toAddPost = true
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
				Content:   contentComment,
				UserID:    idUser,
				Name_User: name_User,
				Like:      utils.GetLikeByCommentID(id_comment, db),
				Dislike:   utils.GetDislikeByCommentID(id_comment, db),
			}

			comments = append([]models.CommentWithUser{commentWithUser}, comments...)

			// comments = append(comments, commentWithUser)
		}

		if err != nil {
			utils.Handle500Error(w)

			return
		}

		post.Date = utils.FormatDate(post.Date)
		post.Comments = comments
		err = db.QueryRow("SELECT mail_user FROM users WHERE id = ?", post.UserID).Scan(&post.MailUser)
		if err != nil {
			utils.Handle500Error(w)
			return
		}
		if strings.Contains(strings.ToLower(post.Title), valueToSearch) || strings.Contains(strings.ToLower(post.Content), valueToSearch) || toAddPost {
			Posts = append([]models.Post{post}, Posts...)
			toAddPost = false
		}
	}

	if err := rows.Err(); err != nil {

		utils.Handle500Error(w)
		return
	}
	sessionData, _ := getSessionData(w, r, db)

	InfoUser := Info{
		IdUserConnect:   sessionData.UserID,
		NameUserConnect: sessionData.Username,
		State:           true,
	}
	Exec := map[string]interface{}{"Posts": Posts, "Infos": InfoUser}
	utils.RenderTemplate(w, "index", Exec)

	// return Posts
}
