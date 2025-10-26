package domain

type Client struct {
	ID               uint
	ClientID         string
	ClientSecretHash string
	RedirectURIs     string
	Grants           string
	IsConfidential   bool
}
