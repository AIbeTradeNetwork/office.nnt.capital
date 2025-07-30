package domain

import "time"

type AuthType string

const (
	AuthTypeUndefined AuthType = ""
	AuthTypePassword  AuthType = "password"
	AuthTypeTelegram  AuthType = "telegram"
)

type UserAuth struct {
	UserUID   string
	Type      AuthType
	Token     string
	CreatedAt time.Time
}

type TelegramAuth struct {
	ID        string
	FirstName string
	LastName  string
	Username  string
	PhotoUrl  string
	AuthDate  string
	Hash      string
	RefUid    string
}

// LoginReq model of authorisation request
type LoginReq struct {
	Login    string
	Password string
}

// RegisterReq model of registration request
type RegisterReq struct {
	Email      string
	Nickname   string
	Password   string
	RePassword string
	RefUid     string
	Locale     string
}

type RefreshReq struct {
	RefreshToken string
}

type AuthRes struct {
	AuthToken    string
	RefreshToken string
}
