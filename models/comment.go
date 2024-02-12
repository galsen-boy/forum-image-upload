package models

type Comment struct {
	Id              int
	Content_comment string
	Id_post         int
	Id_user         int
}

type CommentWithUser struct {
	Content    string
	UserID     int
	Name_User  string
	Id_comment int
	Like       []int
	Dislike    []int
}
