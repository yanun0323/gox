package repository

import "context"

//go:generate esc-domain-gen -f=member.go
type MemberRepository interface {
	Start(ctx context.Context, req *UpdatePhoneReq) (*UpdatePhoneResp, error)
}
