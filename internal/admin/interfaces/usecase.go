package interfaces

type AdminUsecase interface {
	AssignSuperAdmin(email string) error
	RevokeSuperAdmin(email string) error
}
