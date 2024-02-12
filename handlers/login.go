package handlers

import (
	"database/sql"
	"fmt"
	"forum/models"
	"forum/utils"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/login" {

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
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		Login_handler(w, r, db, &user)
	}

	// http.Redirect(w, r, "/", http.StatusFound)
}

type UserConnec struct {
	Id       int
	Liked    bool
	DisLiked bool
}

func Login_handler(w http.ResponseWriter, r *http.Request, db *sql.DB, user *models.User) {

	switch r.Method {
	case "GET":

		utils.RenderTemplate(w, "login", nil)
		// http.ServeFile(w, r, "./templates/login.html")
		return
	case "POST":
		if r.FormValue("login") == "Log In" {
			mail_user := (r.PostFormValue("email"))
			if !(utils.IsvalidEmail(mail_user)) {
				err := map[string]interface{}{"Error": "Mail inccorect"}
				utils.RenderTemplate(w, "login", err)
				return
			}
			user.Mail_user = mail_user
			passwd_user := r.PostFormValue("passwd")
			user.Passwd = passwd_user
			if user.Mail_user == "" || user.Passwd == "" {

				fmt.Println("Password or Email Empty")
				err := map[string]interface{}{"Error": "Password or Email Empty"}
				utils.RenderTemplate(w, "login", err)
				return
			}

			if !utils.IsvalidEmail(mail_user) {
				fmt.Println("❌ Email not valid")
				err := map[string]interface{}{"Error": "❌ Email not valid"}
				utils.RenderTemplate(w, "register", err)
				return
			}
			// Créez une nouvelle connexion à la base de données ici (si ce n'est pas déjà fait).
			req, err := db.Prepare("SELECT mail_user FROM users WHERE mail_user = ?")
			if err != nil {

				utils.Handle500Error(w)
				return
			}
			defer req.Close()
			mail, err := req.Query(user.Mail_user)
			if err != nil {
				fmt.Println("E-mail et ou mot passe incorecte")
				err := map[string]interface{}{"Error": "E-mail or Password incorrect"}
				utils.RenderTemplate(w, "login", err)
				return
			}
			defer mail.Close()
			var mailUser string
			for mail.Next() {
				if err := mail.Scan(&mailUser); err != nil {

					utils.Handle500Error(w)
					return
				}
			}
			stmt, err := db.Prepare("SELECT id, name_user, password_user FROM users WHERE mail_user = ?")
			if err != nil {

				utils.Handle500Error(w)
				return
			}
			defer stmt.Close()
			rows, err := stmt.Query(mailUser)
			if err != nil {

				utils.Handle500Error(w)
				return
			}
			var (
				id         int
				pass, name string
			)
			for rows.Next() {
				if err := rows.Scan(&id, &name, &pass); err != nil {

					utils.Handle500Error(w)
					return
				}
			}
			user.Id = id
			trouve := utils.Decrytper(user.Passwd, []byte(pass))
			if trouve {
				fmt.Println("Authentification Réussie")

				sessionData := models.SessionData{UserID: user.Id, Username: name, IsLoggedIn: true}
				// clearSession(w, r, db, id)
				_, _ = db.Exec("DELETE FROM sessions WHERE id_user = ?", id)

				createSession(w, r, sessionData, db)
				// utils.RenderTemplate(w, "index", nil)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			} else {
				fmt.Println("E-mail et ou mot passe incorecte")
				err := map[string]interface{}{"Error": "Login or Password incorrect"}
				utils.RenderTemplate(w, "login", err)
				// http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
		}
		// http.ServeFile(w, r, "./templates/index.html")
		// return
	default:
		utils.Handle405Error(w)
	}
	// utils.RenderTemplate(w, "login", nil)
}
