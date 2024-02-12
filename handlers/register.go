package handlers

import (
	"database/sql"
	"fmt"
	"forum/models"
	"forum/utils"
	"net/http"
	"strings"
)

func Register(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/register" {

		utils.Handle404Error(w)
		return
	}
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		//
		// utils.Handle500Error(w)
		return
	}
	defer db.Close()
	register_handler(w, r)
}

func register_handler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		utils.RenderTemplate(w, "register", nil)
		// http.ServeFile(w, r, "./templates/register.html")
		return
	}
	var sessionData models.SessionData

	if r.Method == "POST" {
		if r.FormValue("register") == "Sign Up" {
			fmt.Println(r.FormValue("register"))
			name_user := r.PostFormValue("username")
			mail_user := r.PostFormValue("email")
			passwd_user := r.PostFormValue("passwd")
			passwd_again := r.PostFormValue("passwd_again")
			name := name_user
			pass := passwd_user

			if len(pass) < 4 {
				fmt.Println("❌ Password Min character 4")
				err := map[string]interface{}{"Error": "❌ Password Min character 4"}
				utils.RenderTemplate(w, "register", err)
				return
			}
			if name_user == "" || mail_user == "" || passwd_again == "" || passwd_user == "" || len(strings.Fields(name)) == 0 {

				fmt.Println("❌ All input are required")
				err := map[string]interface{}{"Error": "❌ All input are required"}
				utils.RenderTemplate(w, "register", err)
				return
			}
			if !utils.IsvalidEmail(mail_user) {
				fmt.Println("❌ Email not valid")
				err := map[string]interface{}{"Error": "❌ Email not valid"}
				utils.RenderTemplate(w, "register", err)
				return
			}

			if passwd_user != passwd_again {
				fmt.Println("Login again incorrecte")
				err := map[string]interface{}{"Error": "❌ Password don't match"}
				utils.RenderTemplate(w, "register", err)
				return

			}
			passwd_user = string(utils.Crytper(passwd_user))

			/* Créer une nouvelle connexion à la base de données ici (si ce n'est pas déjà fait).*/
			db, err := sql.Open("sqlite3", "./db/forum.db")
			if err != nil {
				fmt.Println("Erreur se produite")
				utils.Handle500Error(w)
				return
			}
			defer db.Close() // Assurez-vous que db.Close() est différé ici.

			user.Name_user = (name_user)
			user.Mail_user = (mail_user)
			user.Passwd = (passwd_user)

			if !(GetEmail(db, user.Mail_user)) {

				fmt.Println("This mail exist !!!")
				err := map[string]interface{}{"Error": "❌ This mail is already used"}
				utils.RenderTemplate(w, "register", err)
				return

			}
			req, err := db.Prepare("INSERT INTO users (name_user, mail_user, password_user) VALUES (?, ?, ? )")
			if err != nil {

				utils.Handle500Error(w)
				return
			}

			_, err = req.Exec(user.Name_user, user.Mail_user, user.Passwd)
			if err != nil {
				fmt.Println(err.Error())

				utils.Handle500Error(w)
				return
			}
			fmt.Println("Success")
			err = db.QueryRow("SELECT id FROM users WHERE name_user = ? AND mail_user = ?", user.Name_user, user.Mail_user).Scan(&user.Id)
			if err != nil {
				fmt.Println("err")
			}
			sessionData = models.SessionData{UserID: user.Id, Username: user.Name_user, IsLoggedIn: true}
			createSession(w, r, sessionData, db)
		}
	} else {
		utils.Handle405Error(w)
	}
	// err := map[string]interface{}{"Success": "Success ✅ Please login !"}

	// utils.RenderTemplate(w, "register", err)
	http.Redirect(w, r, "/", http.StatusSeeOther)

	// register.Execute(w, nil)
}

func GetEmail(db *sql.DB, mail string) bool {
	sqlStatement := "SELECT id FROM users WHERE mail_user = ?"
	row := db.QueryRow(sqlStatement, mail)
	var userID int
	err := row.Scan(&userID)
	return err != nil
}
