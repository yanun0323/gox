package example

import "context"

type exampleRepo struct {
	// Replace by domaingen
	// TODO: Implement me
}

func NewExampleRepository() ExampleRepository {
	// Replace by domaingen
	// TODO: Implement me
	return &exampleRepo{}
}

func (r *exampleRepo) Create(context.Context, *Example) error {
	// Replace by domaingen
	// TODO: Implement me
}

func (r *exampleRepo) Update(context.Context, *Example) error {
	// Replace by domaingen
	// TODO: Implement me
}

func (r *exampleRepo) Delete(context.Context, int64) error {
	// Replace by domaingen
	// TODO: Implement me
}

func (r *exampleRepo) Embed() {
	// TODO: Implement me
}

func (r *exampleRepo) Embed2() {
	// TODO: Implement me
}
