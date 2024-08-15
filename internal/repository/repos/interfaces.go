package repos

import (
	"app/pkg/domain/entity"
	"github.com/google/uuid"
)

type Repositories struct {
	Auth IAuthRepo
	User IUserRepo
}

type IAuthRepo interface {
	GetById(uuid.UUID) (*entity.UserAuth, error)
	Update(*entity.UserAuth) error
	Create(entity.UserAuth) error
	DeleteItem(uuid.UUID) error
}

type IUserRepo interface {
	FindById(uuid.UUID) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	CreatePersonal(newUser *entity.User) error
}
