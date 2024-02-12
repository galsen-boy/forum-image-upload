package utils

import (
	"database/sql"
)

func GetLikeByCommentID(id_comment int, db *sql.DB) []int {
	var likes []int
	rows, err := db.Query("SELECT id_user FROM likes_comment WHERE id_comment=? AND isLike=?", id_comment, true)
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
	// fmt.Println("nbre ike =", likes)
	return likes
}

func GetDislikeByCommentID(id_comment int, db *sql.DB) []int {
	var likes []int
	rows, err := db.Query("SELECT id_user FROM likes_comment WHERE id_comment=? AND isLike=?", id_comment, false)
	if err != nil {
		return likes
	}
	defer rows.Close()
	for rows.Next() {
		var p int
		rows.Scan(&p)
		likes = append(likes, p)
	}
	// fmt.Println("nbre de dislike = ", likes)
	return likes
}
