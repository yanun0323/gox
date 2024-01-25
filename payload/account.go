package payload

//go:generate inspector -foo=bar -arg2 -arg3=123 -o ./internal/entity
type UpdateEmailReq struct {
	Email       string `json:"email" binding:"required"`
	CaptchaCode string `json:"code" binding:"required"`
}

//go:generate goentity -dir=123
type UpdatePhoneReq struct {
	Phone       string `json:"phone" binding:"required"`
	AreaCode    string `json:"area_code" binding:"required"`
	CaptchaCode string `json:"code" binding:"required"`
}

type UpdateNicknameReq struct {
	Nickname string `json:"nickname" binding:"required"`
}

type UpdateFiatCurrencyReq struct {
	FiatCurrency string `json:"fiat_currency" binding:"required,uppercase"`
}

type GetPublicProfileReq struct {
	UID string `form:"uid"`
}

type UpdateLanguageReq struct {
	Language string `json:"language"`
}
