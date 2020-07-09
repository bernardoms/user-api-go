package repository

import (
	"github.com/bernardoms/user-api/internal/model"
)

type UserRepository interface {
	UpdateByNickname(nickname string, user *model.User) (int64, error)
	Save(user *model.User) (*model.User, error)
	FindByNickname(nickname string) (*model.User, error)
	FindAll() ([]*model.User, error)
	FindAllByFilter(filter *model.Filter) ([]model.User, error)
	Delete(nickname string) error
}
