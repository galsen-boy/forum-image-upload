package utils

import (
	"database/sql"
	"fmt"
)

func GetLikeByPostID(id_post int) []int {
	var likes []int
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		fmt.Println(err)
		db.Close()
		return likes
	}
	rows, err := db.Query("SELECT id_user FROM likes_post WHERE id_post=$1 AND isLike=$2", id_post, true)
	if err != nil {
		db.Close()
		return likes
	}
	defer rows.Close()

	for rows.Next() {
		var p int
		rows.Scan(&p)
		likes = append(likes, p)
	}
	db.Close()
	// fmt.Println("nbre ike =", likes)
	return likes
}

// recuperation des dislike a un post donnés
func GetDislikeByPOstID(id_post int) []int {
	var likes []int
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		fmt.Println(err)
		db.Close()
		return likes
	}
	rows, err := db.Query("SELECT id_user FROM likes_post WHERE id_post=$1 AND isLike=$2", id_post, false)
	if err != nil {
		return likes
	}
	defer rows.Close()
	for rows.Next() {
		var p int
		rows.Scan(&p)
		likes = append(likes, p)
	}
	db.Close()
	// fmt.Println("nbre de dislike = ", likes)
	return likes
}
