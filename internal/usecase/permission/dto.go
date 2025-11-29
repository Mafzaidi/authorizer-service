package permission

type (
	CreateInput struct {
		AppID       string
		Code        string
		Description string
	}

	UpdateInput struct {
		FullName string
		Phone    string
	}

	PermissionsInput struct {
		Code        string
		Description string
	}

	SyncInput struct {
		AppCode     string
		Permissions []*PermissionsInput
		Version     int
	}
)
