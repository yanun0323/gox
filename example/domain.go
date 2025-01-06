package example

import (
	"context"
)

// // go:generate modelgen -destination=../output/entity/example.go -package=entity -name=ExampleEntity -struct -construct -tagged -replace -relative
type Example struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Key       string `gorm:"column:key"`
	Msg       string `gorm:"column:message"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime"`
	Extension *ExampleExtension
}

type ExampleExtension struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//go:generate domaingen -destination=../example_output/usecase/example.go -package=usecase -name=exampleUsecase
type ExampleUsecase interface {
	// Run
	Run()

	/* RunWith */
	RunWith(context.Context) error

	/*
	 RunWithString

	 hello
	*/
	RunWithString(ctx context.Context, s1, s2, s3 string) error

	// RunWithElement
	//
	// hello
	RunWithElement(context.Context, ExampleRequest) (*ExampleResponse /* response */, error /* error */)
}

type ExampleRequest struct {
	Key   string
	Value any
}

type ExampleResponse struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

//go:generate domaingen -replace -destination=../example_output/repository/example.go -package=repository
//go:generate domaingen -destination=same_folder_file.go -name=exampleRepo -package=example
//go:generate domaingen -destination=./output/output.go -name=exampleRepo -package=output
type ExampleRepository interface {
	EmbedInterface
	EmbedInterface3

	Create(context.Context, *Example) error
	Update(context.Context, *Example) error
	Delete(context.Context, int64) error
}

type EmbedInterface interface {
	EmbedInterface2

	Embed()
}

type EmbedInterface2 interface {
	Embed2()
}

type EmbedInterface3 interface {
	Embed3()
}
