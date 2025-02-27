package usecase

import "github.com/ClearingHouse/internal/namespaces/interfaces"

type NamespacesUsecase struct {
	namespacesRepository interfaces.NamespacesRepository
}

func NewNamespacesUsecase(namespacesRepository interfaces.NamespacesRepository) interfaces.NamespacesUsecase {
	return &NamespacesUsecase{namespacesRepository: namespacesRepository}
}
