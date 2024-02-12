package handlers

import (
	"database/sql"
	"fmt"
	"forum/models"
	"forum/utils"
	"net/http"
	"strings"
	"time"
)

var user models.User

func CreatePOST(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		fmt.Println("Error when opening db on : Func CreatePOST")
		utils.Handle500Error(w)
		return
	}
	defer db.Close()
	sessionData, isLoggedIn := getSessionData(w, r, db)

	if isLoggedIn && isValidSession(sessionData) {
		create(w, r, db, &sessionData)
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func create(w http.ResponseWriter, r *http.Request, db *sql.DB, sessionData *models.SessionData) {
	if r.Method == "GET" {

		InfoUser := Info{
			NameUserConnect: sessionData.Username,
			State:           true,
		}
		UserConnected.Id = sessionData.UserID
		Exec := map[string]interface{}{"Infos": InfoUser}

		utils.RenderTemplate(w, "createPost", Exec)
		return
	} else if r.Method == "POST" {

		if r.FormValue("send") == "send" {
			var post models.Post
			title_post := r.PostFormValue("title_post")
			content_post := r.PostFormValue("content_post")
			if len(strings.Fields(title_post)) == 0 || len(strings.Fields(content_post)) == 0 {
				err := map[string]interface{}{"Error": " TITLE AND/OR DESCRIPTION EMPTY", "User": user}

				utils.RenderTemplate(w, "createPost", err)
				return
			}

			//Upload img

			var media_post string
			// _, _, err := r.FormFile("media_post")
			erro := uploadHandler(w, r, &media_post)
			if erro != "ok" && erro != "no file" {
				// media_post = header.Filename
				err := map[string]interface{}{"Error": "Error ", "User": user}
				utils.RenderTemplate(w, "createPost", err)

				return
			}

			// End upload
			heureActuelle := time.Now()
			format := "2006-01-02 15:04:05" //"YYYY-MM-DD HH:MM:SS"
			date_post := heureActuelle.Format(format)
			technologie := handleCategory(r.PostFormValue("techno"))
			sante := handleCategory(r.PostFormValue("sante"))
			sport := handleCategory(r.PostFormValue("sport"))
			music := handleCategory(r.PostFormValue("music"))
			news := handleCategory(r.PostFormValue("news"))
			other := handleCategory(r.PostFormValue("other"))

			req, err := db.Prepare("SELECT id FROM category WHERE name_category = ? OR name_category = ? OR name_category = ? OR name_category = ? OR name_category = ? OR name_category = ?")
			if err != nil {

				utils.Handle500Error(w)
				return
			}
			defer req.Close()

			rows, err := req.Query(technologie, sante, sport, music, news, other)
			if err != nil {

				utils.Handle500Error(w)
				return
			}
			defer rows.Close()

			// Parcourir les résultats
			tabID_Category := []int{}
			for rows.Next() {
				var id int
				if err := rows.Scan(&id); err != nil {

					utils.Handle500Error(w)
					return
				}
				tabID_Category = append(tabID_Category, id)
			}
			// fmt.Println(tabID_Category)
			if len(tabID_Category) == 0 {
				err := map[string]interface{}{"Error": " CATEGORY IS NOT CHECKED", "User": user}

				utils.RenderTemplate(w, "createPost", err)
				return
			}

			if err := rows.Err(); err != nil {

				utils.Handle500Error(w)
				return
			}

			// Mettre les donnés dans la structure post
			post.Title = title_post
			post.Content = content_post
			post.Media = media_post
			post.Date = date_post

			// Insertion for Posts
			sql, err := db.Prepare("INSERT INTO posts (title_post, content_post, media_post, date_post, id_user) VALUES (?, ?, ?, ?, ?)")
			if err != nil {

				utils.Handle500Error(w)
				return
			}

			// Execution de la requette INSERT
			user.Id = sessionData.UserID

			_, err = sql.Exec(post.Title, post.Content, post.Media, post.Date, user.Id)
			if err != nil {

				utils.Handle500Error(w)
				return
			}

			id_post, err := db.Prepare("SELECT id FROM posts WHERE title_post = ? AND date_post = ? ")
			if err != nil {

				utils.Handle500Error(w)
				return
			}
			id, err := id_post.Query(post.Title, post.Date)

			var id_Post int
			for id.Next() {
				if err := id.Scan(&id_Post); err != nil {

					utils.Handle500Error(w)
					return
				}
			}

			if err != nil {

				utils.Handle500Error(w)
				return
			}
			defer id.Close()

			// Insertion de l'id_post et de(s) id_category dans la table belong
			for i := 0; i < len(tabID_Category); i++ {
				// for i := range tabID_Category {
				req_b, err := db.Prepare("INSERT INTO belong (id_post, id_category) VALUES (?, ?)")
				if err != nil {

					utils.Handle500Error(w)
					return
				}
				_, err = req_b.Exec(id_Post, tabID_Category[i])
				if err != nil {
					fmt.Println("Problem dans l'insertion de la categorie post")

					utils.Handle500Error(w)
					return
				}

			}

		}
		// var errRet map[string]interface{}

		// errRet = map[string]interface{}{"Success": "✅ POST CREATED", "User": user}
	} else {
		utils.Handle405Error(w)
	}
	http.Redirect(w, r, "/", http.StatusFound)

	// utils.RenderTemplate(w, "createPost", errRet)

}

func handleCategory(category string) string {
	if category == "" {
		return "NULL"
	}
	return category
}
