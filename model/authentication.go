package model

type RegisterRequest struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type AuthenticateResponse struct {
	Token string `json:"token,omitempty"`
}

type AuthError struct {
	Error string `json:"error,omitempty"`
}
