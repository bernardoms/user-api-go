package handler

import (
	"github.com/bernardoms/user-api/config"
	"github.com/bernardoms/user-api/internal/handler"
	"github.com/bernardoms/user-api/internal/logger"
	"github.com/bernardoms/user-api/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPublishUserSuccess(t *testing.T) {

	sns := handler.NewSNS(config.NewSnsConfig(), logger.ConfigureLogger())

	user := new(model.User)
	user.Email = "test@test.com"
	user.Country = "UK"
	user.LastName = "lastName"
	user.FirstName = "firstName"
	user.Password = "password"
	user.Nickname = "testnickname"

	err := sns.Publish(user)

	assert.Equal(t, nil, err, "exception on publish to sns")
}

func TestPublishUserFail(t *testing.T) {

	sns := handler.NewSNS(config.NewSnsConfig(), logger.ConfigureLogger())
	sns.Topic = "failed-topic"

	user := new(model.User)
	user.Email = "test@test.com"
	user.Country = "UK"
	user.LastName = "lastName"
	user.FirstName = "firstName"
	user.Password = "password"
	user.Nickname = "testnickname"

	err := sns.Publish(user)

	assert.Error(t, err)
}
