package application

type Usecase interface {
	Create(input *CreateInput) error
}
