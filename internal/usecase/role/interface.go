package role

type Usecase interface {
	Create(name, description, application string) error
}
