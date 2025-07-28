package app

import (
	OrganizationHttp "github.com/ClearingHouse/internal/organizations/delivery/http"
	OrganizationRepository "github.com/ClearingHouse/internal/organizations/repository"
	OrganizationUsecase "github.com/ClearingHouse/internal/organizations/usecase"

	ResourceHttp "github.com/ClearingHouse/internal/resources/delivery/http"
	ResourceRepository "github.com/ClearingHouse/internal/resources/repository"
	ResourceUsecase "github.com/ClearingHouse/internal/resources/usecase"

	QuotaHttp "github.com/ClearingHouse/internal/quota/delivery/http"
	QuotaRepository "github.com/ClearingHouse/internal/quota/repository"
	QuotaUsecase "github.com/ClearingHouse/internal/quota/usecase"
)

func (a *App) MapHandlers() error {
	// usersGroup := a.gin.Group("/users")
	organizationsGroup := a.gin.Group("/organizations")
	resourcesGroup := a.gin.Group("/resources")
	quotaGroup := a.gin.Group("/quota")

	orgRepo := OrganizationRepository.NewOrganizationRepository(a.postgresDB)
	resourceRepo, resourcePoolRepo, resourceTypeRepo := ResourceRepository.NewResourceRepository(a.postgresDB)
	quotaRepo := QuotaRepository.NewQuotaRepository(a.postgresDB)

	orgUsecase := OrganizationUsecase.NewOrganizationUsecase(orgRepo)
	resourceUsecase := ResourceUsecase.NewResourceUsecase(resourceRepo, resourcePoolRepo, resourceTypeRepo)
	quotaUsecase := QuotaUsecase.NewQuotaUsecase(quotaRepo, resourcePoolRepo)

	orgHandler := OrganizationHttp.NewOrganizationHandler(orgUsecase)
	resourceHandler := ResourceHttp.NewResourceHandler(resourceUsecase)
	quotaHandler := QuotaHttp.NewQuotaHandler(quotaUsecase)

	OrganizationHttp.MapOrganizationRoutes(organizationsGroup, orgHandler)
	ResourceHttp.MapResourceRoutes(resourcesGroup, resourceHandler)
	QuotaHttp.MapQuotaRoutes(quotaGroup, quotaHandler)

	// usersRepository := usersRepository.NewUsersRepository(a.postgresDB)

	// usersUsecase := usersUsecase.NewUsersUsecase(usersRepository)

	// usersHandlers := usersHttp.NewUsersHandler(usersUsecase)

	// usersHttp.MapUsersRoutes(usersGroup, usersHandlers)

	return nil
}
