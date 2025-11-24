package user

type (
	RegisterInput struct {
		Username string
		FullName string
		Phone    string
		Email    string
		Password string
	}

	UpddateInput struct {
		FullName string
		Phone    string
	}
)
