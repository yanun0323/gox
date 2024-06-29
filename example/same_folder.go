package example

import "context"

type repo struct {
	// TODO: Implement me
}

func (repo *repo) Embed() {
	// TODO: Implement me
}

type exampleRepo struct {
	// TODO: Implement me
}

func NewExampleRepository() ExampleRepository {
	// TODO: Implement me
	return &exampleRepo{}
}

func (repo *exampleRepo) Update(context.Context, *Example) error {
	// TODO: Implement me
}

func (repo *exampleRepo) Delete(context.Context, int64) error {
	// TODO: Implement me
}

func (repo *exampleRepo) Embed() {
	// TODO: Implement me
}

func (repo *exampleRepo) Embed2() {
	// TODO: Implement me
}

func (repo *exampleRepo) Create(context.Context, *Example) error {
	// TODO: Implement me
}
