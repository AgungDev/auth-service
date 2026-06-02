package dto

type ClientRequest struct {
	ClientID       string   `json:"client_id" binding:"required"`
	Name           string   `json:"name" binding:"required"`
	ClientSecret   string   `json:"client_secret" binding:"required"`
	RedirectURIs   []string `json:"redirect_uris" binding:"required,min=1,dive,required"`
	Grants         []string `json:"grants" binding:"required,min=1,dive,required"`
	IsConfidential bool     `json:"is_confidential"`
}

type ClientResponse struct {
	ID             string   `json:"id"`
	ClientID       string   `json:"client_id"`
	Name           string   `json:"name"`
	RedirectURIs   []string `json:"redirect_uris"`
	Grants         []string `json:"grants"`
	IsConfidential bool     `json:"is_confidential"`
}
