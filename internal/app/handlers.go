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

	ProjectHttp "github.com/ClearingHouse/internal/projects/delivery/http"
	ProjectRepository "github.com/ClearingHouse/internal/projects/repository"
	ProjectUsecase "github.com/ClearingHouse/internal/projects/usecase"

	NamespaceHttp "github.com/ClearingHouse/internal/namespaces/delivery/http"
	NamespaceRepository "github.com/ClearingHouse/internal/namespaces/repository"
	NamespaceUsecase "github.com/ClearingHouse/internal/namespaces/usecase"

	AuthHttp "github.com/ClearingHouse/internal/auth/delivery/http"
	AuthUsecase "github.com/ClearingHouse/internal/auth/usecase"

	UserHttp "github.com/ClearingHouse/internal/users/delivery/http"
	UserRepository "github.com/ClearingHouse/internal/users/repository"
	UserUsecase "github.com/ClearingHouse/internal/users/usecase"
)

func (a *App) MapHandlers() error {
	// usersGroup := a.gin.Group("/users")
	organizationsGroup := a.gin.Group("/organizations")
	resourcesGroup := a.gin.Group("/resources")
	quotaGroup := a.gin.Group("/quota")
	projectsGroup := a.gin.Group("/projects")
	namespacesGroup := a.gin.Group("/namespaces")
	authGroup := a.gin.Group("/auth")
	userGroup := a.gin.Group("/users")

	orgRepo := OrganizationRepository.NewOrganizationRepository(a.postgresDB)
	resourceRepo, resourcePoolRepo, resourceTypeRepo := ResourceRepository.NewResourceRepository(a.postgresDB)
	quotaRepo := QuotaRepository.NewQuotaRepository(a.postgresDB)
	projRepo := ProjectRepository.NewProjectRepository(a.postgresDB)
	namespaceRepo := NamespaceRepository.NewNamespaceRepository(a.postgresDB)
	userRepo := UserRepository.NewUsersRepository(a.postgresDB)

	orgUsecase := OrganizationUsecase.NewOrganizationUsecase(orgRepo, userRepo)
	resourceUsecase := ResourceUsecase.NewResourceUsecase(resourceRepo, resourcePoolRepo, resourceTypeRepo)
	quotaUsecase := QuotaUsecase.NewQuotaUsecase(quotaRepo, resourcePoolRepo, namespaceRepo, orgRepo, projRepo, userRepo)
	projUsecase := ProjectUsecase.NewProjectUsecase(projRepo, orgRepo, userRepo)
	namespaceUsecase := NamespaceUsecase.NewNamespaceUsecase(namespaceRepo, userRepo, projRepo)
	userUsecase := UserUsecase.NewUsersUsecase(userRepo)
	authUsecase := AuthUsecase.NewAuthUsecase(userRepo)

	orgHandler := OrganizationHttp.NewOrganizationHandler(orgUsecase)
	resourceHandler := ResourceHttp.NewResourceHandler(resourceUsecase)
	quotaHandler := QuotaHttp.NewQuotaHandler(quotaUsecase)
	projHandler := ProjectHttp.NewProjectHandler(projUsecase)
	namespaceHandler := NamespaceHttp.NewNamespaceHandler(namespaceUsecase)
	userHandler := UserHttp.NewUsersHandler(userUsecase)
	authHandler := AuthHttp.NewAuthHandler(authUsecase)

	OrganizationHttp.MapOrganizationRoutes(organizationsGroup, orgHandler)
	ResourceHttp.MapResourceRoutes(resourcesGroup, resourceHandler)
	QuotaHttp.MapQuotaRoutes(quotaGroup, quotaHandler)
	ProjectHttp.MapProjectRoutes(projectsGroup, projHandler)
	NamespaceHttp.MapNamespaceRoutes(namespacesGroup, namespaceHandler)
	UserHttp.MapUsersRoutes(userGroup, userHandler)
	AuthHttp.MapAuthRoutes(authGroup, authHandler)

	// usersRepository := usersRepository.NewUsersRepository(a.postgresDB)
	// usersUsecase := usersUsecase.NewUsersUsecase(usersRepository)

	// usersHandlers := usersHttp.NewUsersHandler(usersUsecase)

	// usersHttp.MapUsersRoutes(usersGroup, usersHandlers)

	return nil
}
