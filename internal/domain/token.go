package domain

type Token struct {
	ID           uint
	UserID       uint
	ClientID     uint
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
}
