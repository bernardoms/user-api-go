package integration

import (
	"github.com/bernardoms/user-api/config"
	"github.com/bernardoms/user-api/internal/model"
	"github.com/bernardoms/user-api/internal/repository"
	"os"
	"testing"
)

var (
	i *repository.Mongo
)

func TestMain(m *testing.M) {
	c := config.NewMongoConfig()
	repository.New(c)
	i = &repository.Mongo{Collection: repository.GetUserCollection(c)}

	user := model.User{
		Email:     "test1@test.com",
		FirstName: "firstName",
		LastName:  "lastName",
		Nickname:  "testnick1",
		Country:   "UK",
	}

	_, _ = i.Save(&user)

	user2 := model.User{
		Email:     "test2@test.com",
		FirstName: "firstName2",
		LastName:  "lastName2",
		Nickname:  "testnick2",
		Country:   "UK",
	}

	_, _ = i.Save(&user2)

	os.Exit(m.Run())
}
