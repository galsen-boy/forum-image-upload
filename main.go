package main

import (
	"database/sql"
	"fmt"
	"forum/handlers"
	"forum/utils"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	clearScreen()
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		//
		// 	Handle500Error(w)
		return
	}
	utils.CreateTable(db)
	defer db.Close()

	//Parse css file
	css := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", css))
	img := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", img))

	// Parse js files
	// js := http.FileServer(http.Dir("js"))
	// http.Handle("/js/", http.StripPrefix("/js/", js))

	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/createPost", handlers.CreatePOST)
	http.HandleFunc("/profil", handlers.Profil)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/logout", handlers.Logout)
	http.HandleFunc("/comment", handlers.Comment)
	http.HandleFunc("/like", handlers.Liked)
	http.HandleFunc("/search", handlers.Search)
	http.HandleFunc("/filter", handlers.Filter)
	http.HandleFunc("/like_comment", handlers.LikeComment)

	// http.HandleFunc("/posteGestion", handlers.PosteGestion)
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
