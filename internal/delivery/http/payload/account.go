package payload

//go:generate modelgen -destination=../../../entity/account.go -package=entity -name=UpdatePhone -function=source,target -relative
type UpdatePhoneReq struct {
	Phone       string `json:"phone" binding:"required"`
	AreaCode    string `json:"area_code" binding:"required"`
	CaptchaCode string `json:"code" binding:"required"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
	CreateAt    int64  `json:"create_at"`
	UpdateAt    int64  `json:"update_at"`

	PhoneInfo []*PhoneInfo
	*Pager
}

type UpdatePhoneResp struct {
	Phone       string `json:"phone" binding:"required"`
	AreaCode    string `json:"area_code" binding:"required"`
	CaptchaCode string `json:"code" binding:"required"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
	CreateAt    string `json:"create_at"`
	UpdateAt    string `json:"update_at"`
}

type PhoneInfo struct {
	User        string
	Country     string
	CountryCode string
}

type Pager struct {
	Index int64
	Size  int64
	Count int64
}
