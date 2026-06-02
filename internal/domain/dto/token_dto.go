package dto

type TokenRequest struct {
	GrantType    string `json:"grant_type" binding:"required"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
	RefreshToken string `json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type AuthLoginResponse struct {
	AccessToken  string       `json:"access_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int          `json:"expires_in"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type IntrospectRequest struct {
	Token string `json:"token" binding:"required"`
}

type IntrospectResponse struct {
	Active      bool              `json:"active"`
	User        map[string]string `json:"user,omitempty"`
	Roles       []string          `json:"roles,omitempty"`
	Permissions []string          `json:"permissions,omitempty"`
}

type PermissionsResponse struct {
	Permissions []string `json:"permissions"`
}

type AuthorizeRequest struct {
	UserID     string `json:"user_id" binding:"required"`
	Permission string `json:"permission" binding:"required"`
}

type AuthorizeResponse struct {
	Authorized bool   `json:"authorized"`
	UserID     string `json:"user_id"`
	Permission string `json:"permission"`
}
