package oauth

type User interface {
	GetFullname() string
	GetAvatar() string
	GetEmailAddress() string
}