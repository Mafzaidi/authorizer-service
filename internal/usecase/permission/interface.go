package permission

type Usecase interface {
	Create(input *CreateInput) error
	SyncPermissions(i *SyncInput) error
}
