package app

import (
	usersHttp "github.com/ClearingHouse/internal/users/delivery/http"
	usersRepository "github.com/ClearingHouse/internal/users/repository"
	usersUsecase "github.com/ClearingHouse/internal/users/usecase"
)

func (a *App) MapHandlers() error {
	usersGroup := a.gin.Group("/users")

	usersRepository := usersRepository.NewUsersRepository(a.postgresDB)

	usersUsecase := usersUsecase.NewUsersUsecase(usersRepository)

	usersHandlers := usersHttp.NewUsersHandler(usersUsecase)

	usersHttp.MapUsersRoutes(usersGroup, usersHandlers)

	return nil
}
