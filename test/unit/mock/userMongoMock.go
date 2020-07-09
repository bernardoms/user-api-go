package mock

import (
	"github.com/bernardoms/user-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MongoMock struct {
	mock.Mock
}

func (m *MongoMock) FindAll() ([]*model.User, error) {
	args := m.Called()
	return args.Get(0).([]*model.User), args.Error(1)
}

func (m *MongoMock) FindByNickname(nickname string) (*model.User, error) {
	args := m.Called(nickname)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MongoMock) Save(user *model.User) (*model.User, error) {
	args := m.Called(user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MongoMock) UpdateByNickname(nickname string, user *model.User) (int64, error) {
	args := m.Called(nickname, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MongoMock) FindAllByFilter(filter *model.Filter) ([]model.User, error) {
	args := m.Called(filter)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MongoMock) Delete(nickname string) error {
	args := m.Called(nickname)
	return args.Error(0)
}
