package models

type Post struct {
	ID        int
	Title     string
	Content   string
	Media     string
	Date      string
	UserID    int
	NameUser  string
	MailUser  string
	Categorie []string
	Comments  []CommentWithUser
	Like      []int
	Dislike   []int
}
