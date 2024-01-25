package payload

type RegisterReq struct {
	Email        string `json:"email" binding:"required"`
	Password     string `json:"password" binding:"required"`
	SecurityCode string `json:"security_code" binding:"required"`
	SecurityType string `json:"security_type" binding:"required"`
}

type RegisterRes struct {
	AccessToken string   `json:"access_token"`
	TokenType   string   `json:"token_type"`
	ExpiresIn   int64    `json:"expires_in"`
	Security    []string `json:"security"`
}

type RegisterSendEmailCodeReq struct {
	Email string `json:"email" binding:"required"`
}
