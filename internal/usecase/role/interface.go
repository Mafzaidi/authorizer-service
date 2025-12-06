package role

type Usecase interface {
	Create(input *CreateInput) error
	GrantPerms(roleID string, perms []string) error
}
