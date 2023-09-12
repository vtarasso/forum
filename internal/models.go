package internal

import "time"

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword string
	Token          *string
	Created        time.Time
}

type Data struct {
	UserID          int
	UserName        string
	Token           string
	IsAuthenticated bool
	CountLikes      int
	CountDislikes   int
	IntErr          int
}
type Snippet struct {
	ID         int
	UserName   string
	Title      string
	Content    string
	Catigories []string
	CatID      []int
	Likes      int
	Dislikes   int
	Created    time.Time
	Validator  Validator
}
type FormComment struct {
	PostID    int
	UserId    int
	UserName  string
	Content   string
	Validator Validator
}
type SnippetCom struct {
	ID          int
	PostID      int
	UserId      int
	UserName    string
	Content     string
	LikesCom    int
	DislikesCom int
	Created     time.Time
}
type userSignupForm struct {
	Name      string    `form:"name"`
	Email     string    `form:"email"`
	Password  string    `form:"password"`
	Validator Validator `form:"-"`
}
type userLoginForm struct {
	Email     string    `form:"email"`
	Password  string    `form:"password"`
	Validator Validator `form:"-"`
}
type ErrorStruct struct {
	Text   string
	Status int
}
