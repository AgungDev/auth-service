package http

type RegisterRequest struct {
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"john.doe@example.com"`
	Password string `json:"password" example:"secret123"`
	FullName string `json:"full_name" example:"John Doe"`
}

type LoginRequest struct {
	Username string `json:"username" example:"johndoe"`
	Password string `json:"password" example:"secret123"`
	ClientID string `json:"client_id" example:"gateway-app"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token" example:"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType    string       `json:"token_type" example:"Bearer"`
	ExpiresIn    int          `json:"expires_in" example:"900"`
	RefreshToken string       `json:"refresh_token" example:"dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4="`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID       string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"john.doe@example.com"`
	FullName string `json:"full_name" example:"John Doe"`
	Status   string `json:"status" example:"active"`
}

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"validation error"`
	Error   string `json:"error" example:"username is required"`
}

type PermissionRequest struct {
	Permission string `json:"permission" example:"users:read"`
}

type PermissionResponse struct {
	Allowed bool `json:"allowed" example:"true"`
}

type PermissionsResponse struct {
	Permissions []string `json:"permissions" example:"[\"users:read\", \"users:write\"]"`
}

type IntrospectRequest struct {
	Token string `json:"token" example:"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type IntrospectResponse struct {
	Active      bool              `json:"active" example:"true"`
	User        map[string]string `json:"user"`
	Roles       []string          `json:"roles" example:"[\"admin\"]"`
	Permissions []string          `json:"permissions" example:"[\"users:read\"]"`
}

type AuthorizeRequest struct {
	UserID     string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Permission string `json:"permission" example:"users:read"`
}

type AuthorizeResponse struct {
	Authorized bool   `json:"authorized" example:"true"`
	UserID     string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Permission string `json:"permission" example:"users:read"`
}

type MessageResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"logout success"`
}
