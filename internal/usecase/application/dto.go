package application

type (
	CreateInput struct {
		Code        string
		Name        string
		Description string
		Metadata    map[string]interface{}
	}

	UpdateInput struct {
		FullName string
		Phone    string
	}
)
