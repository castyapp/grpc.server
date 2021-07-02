package oauth

type User interface {
	GetUserID() string
	GetFullname() string
	GetAvatar() string
	GetEmailAddress() string
}
