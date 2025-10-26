package app

import (
	"time"

	OrganizationHttp "github.com/ClearingHouse/internal/organizations/delivery/http"
	OrganizationRepository "github.com/ClearingHouse/internal/organizations/repository"
	OrganizationUsecase "github.com/ClearingHouse/internal/organizations/usecase"
	"github.com/gin-gonic/gin"

	ResourceHttp "github.com/ClearingHouse/internal/resources/delivery/http"
	ResourceRepository "github.com/ClearingHouse/internal/resources/repository"
	ResourceUsecase "github.com/ClearingHouse/internal/resources/usecase"

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

	QuotaHttp "github.com/ClearingHouse/internal/quota/delivery/http"
	QuotaRepository "github.com/ClearingHouse/internal/quota/repository"
	QuotaUsecase "github.com/ClearingHouse/internal/quota/usecase"

	TicketHttp "github.com/ClearingHouse/internal/tickets/delivery/http"
	TicketRepository "github.com/ClearingHouse/internal/tickets/repository"
	TicketUsecase "github.com/ClearingHouse/internal/tickets/usecase"

	PrivateNamespaceHttp "github.com/ClearingHouse/internal/private_namespaces/delivery/http"
	PrivateNamespaceRepository "github.com/ClearingHouse/internal/private_namespaces/repository"
	PrivateNamespaceUsecase "github.com/ClearingHouse/internal/private_namespaces/usecase"
)

func (a *App) MapHandlers() error {
	// usersGroup := a.gin.Group("/users")
	a.gin.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to ClearingHouse API",
			"status":  "OK",
			"version": "v1",
			"time":    time.Now(),
		})
	})
	organizationsGroup := a.gin.Group("/organizations")
	resourcesGroup := a.gin.Group("/resources")
	quotaGroup := a.gin.Group("/quota")
	projectsGroup := a.gin.Group("/projects")
	namespacesGroup := a.gin.Group("/namespaces")
	authGroup := a.gin.Group("/auth")
	userGroup := a.gin.Group("/users")
	ticketGroup := a.gin.Group("/tickets")
	privNamespaceGroup := namespacesGroup.Group("/private")

	orgRepo := OrganizationRepository.NewOrganizationRepository(a.postgresDB)
	resourcePoolRepo, resourceRepo, resourceTypeRepo := ResourceRepository.NewResourceRepository(a.postgresDB)
	quotaRepo := QuotaRepository.NewQuotaRepository(a.postgresDB)
	projRepo := ProjectRepository.NewProjectRepository(a.postgresDB)
	namespaceRepo := NamespaceRepository.NewNamespaceRepository(a.postgresDB)
	userRepo := UserRepository.NewUsersRepository(a.postgresDB)
	ticketRepo := TicketRepository.NewTicketRepository(a.postgresDB)
	privNamespaceRepo := PrivateNamespaceRepository.NewPrivateNamespaceRepository(a.postgresDB)

	orgUsecase := OrganizationUsecase.NewOrganizationUsecase(orgRepo, userRepo)
	resourceUsecase := ResourceUsecase.NewResourceUsecase(resourcePoolRepo, resourceRepo, resourceTypeRepo)
	quotaUsecase := QuotaUsecase.NewQuotaUsecase(quotaRepo, resourceRepo, namespaceRepo, orgRepo, projRepo, userRepo)
	projUsecase := ProjectUsecase.NewProjectUsecase(projRepo, orgRepo, userRepo)
	namespaceUsecase := NamespaceUsecase.NewNamespaceUsecase(namespaceRepo, userRepo, projRepo, quotaRepo)
	userUsecase := UserUsecase.NewUsersUsecase(userRepo)
	authUsecase := AuthUsecase.NewAuthUsecase(userRepo)
	ticketUsecase := TicketUsecase.NewTicketUsecase(namespaceRepo, ticketRepo, quotaRepo, userRepo)
	privNamespaceUsecase := PrivateNamespaceUsecase.NewPrivateNamespaceUsecase(privNamespaceRepo, namespaceRepo, orgRepo, userRepo, quotaRepo, resourceRepo)

	orgHandler := OrganizationHttp.NewOrganizationHandler(orgUsecase)
	resourceHandler := ResourceHttp.NewResourceHandler(resourceUsecase)
	quotaHandler := QuotaHttp.NewQuotaHandler(quotaUsecase)
	projHandler := ProjectHttp.NewProjectHandler(projUsecase)
	namespaceHandler := NamespaceHttp.NewNamespaceHandler(namespaceUsecase)
	userHandler := UserHttp.NewUsersHandler(userUsecase)
	authHandler := AuthHttp.NewAuthHandler(authUsecase)
	ticketHandler := TicketHttp.NewTicketHandler(ticketUsecase)
	privNamespaceHandler := PrivateNamespaceHttp.NewPrivateNamespaceHandler(privNamespaceUsecase, namespaceUsecase)

	OrganizationHttp.MapOrganizationRoutes(organizationsGroup, orgHandler)
	ResourceHttp.MapResourceRoutes(resourcesGroup, resourceHandler)
	QuotaHttp.MapQuotaRoutes(quotaGroup, quotaHandler)
	ProjectHttp.MapProjectRoutes(projectsGroup, projHandler)
	NamespaceHttp.MapNamespaceRoutes(namespacesGroup, namespaceHandler)
	UserHttp.MapUsersRoutes(userGroup, userHandler)
	AuthHttp.MapAuthRoutes(authGroup, authHandler)
	TicketHttp.MapTicketRoutes(ticketGroup, ticketHandler)
	PrivateNamespaceHttp.MapPrivateNamespaceRoutes(privNamespaceGroup, privNamespaceHandler)

	return nil
}
