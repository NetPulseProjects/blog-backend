package service

import (
	"app/internal/config"
	"app/internal/repository/repos"
)

type Deps struct {
	Repos  *repos.Repositories
	Config *config.Config
}

type Services struct {
	User User
	Auth Auth
}

func NewServices(deps Deps) *Services {
	authService := NewAuthService(deps.Repos.Auth, deps.Repos.User, deps.Config)
	userService := NewUserService(deps.Repos.User, authService)

	return &Services{
		User: *userService,
		Auth: *authService,
	}
}
