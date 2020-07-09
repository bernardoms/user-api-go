package handler

import "github.com/bernardoms/user-api/internal/model"

type NotifyInterface interface {
	Publish(user *model.User) error
}
