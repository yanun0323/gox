package usecase

import (
	"context"
	"errors"
)

//go:generate domaingen -v -target=../../usecase/member.go
type MemberUseCase interface {
	Start(ctx context.Context, req *UpdatePhoneReq) (*UpdatePhoneResp, error)
}

var (
	ErrNotFound         = errors.New("not found")
	ErrPermissionDenied = errors.New("permission denied")
)

type UpdatePhoneReq struct {
	Phone       string
	AreaCode    string
	CaptchaCode string
	CreateTime  int64
	UpdateTime  int64
	CreateAt    int64
	UpdateAt    int64
}

type UpdatePhoneResp struct {
	Phone       string
	AreaCode    string
	CaptchaCode string
	CreateTime  string
	UpdateTime  string
	CreateAt    string
	UpdateAt    string
}
