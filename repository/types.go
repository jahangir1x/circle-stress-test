package repository

type AccessTokenInfo struct {
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

type TokenNotFoundError struct {
}

func (e *TokenNotFoundError) Error() string {
	return "Token not found"
}
