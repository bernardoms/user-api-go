package mock

import (
	"github.com/bernardoms/user-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type NotifyMock struct {
	mock.Mock
}

func (n *NotifyMock) Publish(user *model.User) error {
	args := n.Called(user)
	return args.Error(0)
}
