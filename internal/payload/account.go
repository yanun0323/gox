package payload

//go:generate inspector -foo=bar -arg2 -arg3=123 -o ./internal/entity
type UpdateEmailReq struct {
	Email       string `json:"email" binding:"required"`
	CaptchaCode string `json:"code" binding:"required"`
}

//go:generate goentity -v -rename=UpdatePhoneEntity -replace=true -unix -entity="member.go" -use="member.go" -repo="member.go"
type UpdatePhoneReq struct {
	Phone       string `json:"phone" binding:"required"`
	AreaCode    string `json:"area_code" binding:"required"`
	CaptchaCode string `json:"code" binding:"required"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
	CreateAt    string `json:"create_at"`
	UpdateAt    string `json:"update_at"`
	Timer
}

type Timer struct {
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
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
