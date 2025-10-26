package domain

type User struct {
	ID           uint
	Username     string
	Email        string
	PasswordHash string
	FullName     string
	IsActive     bool
}
