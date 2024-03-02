package entity

import "gox/internal/domain/repository"

type UpdatePhoneReq struct {
	Phone       string
	AreaCode    string
	CaptchaCode string
	CreateTime  int64
	UpdateTime  int64
	CreateAt    int64
	UpdateAt    int64
	repository.UpdatePhoneEntity
}
