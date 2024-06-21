package example

import (
	"context"
)

type Example struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Key       string `gorm:"column:key"`
	Msg       string `gorm:"column:message"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime"`
}

//go:generate domaingen -destination=../output/usecase/example.go -package=usecase -name=exampleUsecase
type ExampleUsecase interface {
	Run()
	RunWith(context.Context) error
	RunWithString(ctx context.Context, s1, s2, s3 string) error
	RunWithElement(context.Context, ExampleRequest) (*ExampleResponse, error)
}

type ExampleRequest struct {
	Key   string
	Value any
}

type ExampleResponse struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

//go:generate domaingen -replace -destination=../output/repository/example.go -package=repository
type ExampleRepository interface {
	Create(context.Context, *Example) error
	Update(context.Context, *Example) error
	Delete(context.Context, int64) error
}
