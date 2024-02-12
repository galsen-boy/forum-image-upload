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

func Filter(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/filter" {

		utils.Handle404Error(w)
		return
	}
	if r.Method == "POST" {
		var categories []string
		var Posts []models.Post
		technologie := handleCategory(r.PostFormValue("techno"))
		sante := handleCategory(r.PostFormValue("sante"))
		sport := handleCategory(r.PostFormValue("sport"))
		music := handleCategory(r.PostFormValue("music"))
		news := handleCategory(r.PostFormValue("news"))
		other := handleCategory(r.PostFormValue("other"))
		// min_like := handleCategory(r.PostFormValue("liked_post_min"))
		// max_like := handleCategory(r.PostFormValue("liked_post_max"))

		categ := []string{technologie, sante, sport, music, news, other}
		for _, value := range categ {
			if value != "NULL" {
				categories = append(categories, value)
			}
		}
		if len(categories) == 0 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		posts := fetchPostFilter(w, r)
		for _, post := range posts {
			for _, category := range categories {
				category = strings.ToLower(category)
				if strings.Contains(strings.ToLower(strings.Join(post.Categorie, "")), category) {
					Posts = append(Posts, post)
				}
			}
		}
		db, err := sql.Open("sqlite3", "./db/forum.db")
		if err != nil {

			utils.Handle500Error(w)
			return
		}
		defer db.Close()
		sessionData, _ := getSessionData(w, r, db)

		InfoUser := Info{
			NameUserConnect: sessionData.Username,
			State:           true,
		}
		UserConnected.Id = sessionData.UserID
		Exec := map[string]interface{}{"Posts": Posts, "Infos": InfoUser}
		utils.RenderTemplate(w, "index", Exec)
		return
	} else {
		utils.Handle405Error(w)
		return
	}
}

func fetchPostFilter(w http.ResponseWriter, r *http.Request) []models.Post {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {

		utils.Handle500Error(w)
		return nil
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title_post, content_post, media_post, date_post, id_user FROM posts")
	if err != nil {

		utils.Handle500Error(w)
		return nil
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
			return nil
		}

		db_id_category, err := db.Prepare("SELECT id_category FROM belong where id_post = ?")
		if err != nil {

			utils.Handle400Error(w)
			return nil
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
				return nil
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
			return nil
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

			return nil
		}

		content_comment, err := db.Prepare("SELECT content_comment, id_user FROM comment WHERE id_post = ?")
		if err != nil {

			utils.Handle400Error(w)
			return nil
		}
		commentRows, err := content_comment.Query(post.ID)

		var comments []models.CommentWithUser

		for commentRows.Next() {
			var contentComment string
			var id_comment, idUser int

			if err := commentRows.Scan(&contentComment, &idUser); err != nil {
				log.Fatal(err)
			}

			name_user, err := db.Prepare("SELECT name_user FROM users WHERE id = ?")
			if err != nil {

				utils.Handle400Error(w)
				return nil
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

				return nil
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

			return nil
		}

		post.Date = utils.FormatDate(post.Date)
		post.Comments = comments
		err = db.QueryRow("SELECT mail_user FROM users WHERE id = ?", post.UserID).Scan(&post.MailUser)
		if err != nil {
			utils.Handle500Error(w)
		}

		Posts = append([]models.Post{post}, Posts...)
	}

	if err := rows.Err(); err != nil {

		utils.Handle500Error(w)
		return nil
	}

	// Exec := map[string]interface{}{"Posts": Posts}
	// utils.RenderTemplate(w, "index", Exec)
	// http.Redirect(w, r, "/", http.StatusSeeOther)
	return Posts
}
