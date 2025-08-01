package types

import (
	"time"
)

type UserStore interface {
	CreateUser(*User) error
	GetUserByEmail(string) (*User, error)
	GetUserByUsername(string) (*User, error)
	GetUserById(int) (*User, error)
	GetUserFriends(int) ([]*User, error)
	CreateFriendRequest(int, int) error
	DeleteFriend(int, int) error
	RespondFriendRequest(int, int, string) error
}

type RegisterUserPayload struct {
	UserName string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	ID            int    `json:"id"`
	UserName      string `json:"username"`
	Email         string `json:"email"`
	Password      []byte `json:"-"`
	ProfilePicUrl string `json:"profile_pic"`

	CreatedAt time.Time `json:"created_at"`
}

type PostStore interface {
	GetPostById(int) (*Post, error)
	GetPostsFromUser(string) ([]*Post, error)
	CreatePost(Post) (*Post, error)
	EditPost(*Post) error
	DeletePost(int) error
	GetLastTenPosts() ([]*Post, error)
	GetPostsFromFriends(int) ([]*Post, error)
}

type Post struct {
	ID       int    `json:"id"`
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	Title    string `json:"title"`
	Text     string `json:"text"`
	ImgURL   string `json:"img_url"`
	Public   bool   `json:"public"`

	CreatedAt time.Time `json:"created_at"`
}

type CommentStore interface {
	GetCommentsFromUser(int) ([]*Comment, error)
	GetCommentsFromPost(int) ([]*Comment, error)
	CreateComment(*Comment) (*Comment, error)
	DeleteComment(int) error
}

type Comment struct {
	ID       int
	UserID   int
	PostID   int
	Username string
	Text     string

	CreatedAt time.Time `json:"created_at"`
}

type ctxKey string

const (
	CtxKeyUserID   ctxKey = "user_id"
	CtxKeyUsername ctxKey = "username"
)

type Data map[string]any

const AvatarURL string = "https://e7.pngegg.com/pngimages/84/165/png-clipart-united-states-avatar-organization-information-user-avatar-service-computer-wallpaper-thumbnail.png"
