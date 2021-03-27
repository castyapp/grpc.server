package oauth

type User interface {
	GetUserId() string
	GetFullname() string
	GetAvatar() string
	GetEmailAddress() string
}
