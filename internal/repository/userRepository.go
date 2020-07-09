package repository

import (
	"context"
	"github.com/bernardoms/user-api/config"
	"github.com/bernardoms/user-api/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
)

type SessionCreator struct {
	Client *mongo.Client
}

type Mongo struct {
	Collection *mongo.Collection
}

var session *mongo.Client

func New(config *config.MongoConfig) {

	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURI))

	m := SessionCreator{
		Client: client,
	}

	ctx := context.TODO()

	if err != nil {
		log.Print("Error on creating mongo client ", err)
	}

	err = m.Client.Connect(ctx)

	session = client

	if err != nil {
		log.Print("Error on connecting to database ", err)
	}
}

func GetUserCollection(mongoConfig *config.MongoConfig) *mongo.Collection {
	c := session.Database(mongoConfig.Database).Collection("users")
	_, _ = c.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bsonx.Doc{{Key: "nickname", Value: bsonx.String("text")}},
		},
		{
			Keys: bsonx.Doc{{Key: "email", Value: bsonx.String("text")}},
		},
	})
	return c
}

func (m Mongo) FindAll() ([]*model.User, error) {
	var results []*model.User

	cur, err := m.Collection.Find(context.TODO(), bson.D{})

	if err == nil {

		for cur.Next(context.TODO()) {
			var elem model.User
			err := cur.Decode(&elem)
			if err != nil {
				return results, err
			}
			results = append(results, &elem)
		}
	}
	return results, err
}

func (m Mongo) Save(user *model.User) (*model.User, error) {

	_, err := m.Collection.InsertOne(context.TODO(), &user)

	return user, err
}

func (m Mongo) FindByNickname(nickname string) (*model.User, error) {
	var result *model.User

	cur, err := m.Collection.Find(context.TODO(), bson.M{"nickname": nickname})

	if err != nil {
		return &model.User{}, err
	}

	if cur.Next(context.TODO()) != false {
		err = cur.Decode(&result)
	}

	return result, err
}

func (m Mongo) UpdateByNickname(nickname string, user *model.User) (int64, error) {
	filter := bson.M{"nickname": nickname}

	update := bson.M{"$set": bson.M{"nickname": user.Nickname,
		"country":   user.Country,
		"firstname": user.FirstName,
		"lastname":  user.LastName,
		"password":  user.Password,
		"email":     user.Email}}

	r, err := m.Collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return 0, err
	}

	return r.MatchedCount, err
}

func (m Mongo) Delete(nickname string) error {
	filter := bson.M{"nickname": nickname}

	_, err := m.Collection.DeleteOne(context.TODO(), filter)

	return err
}

func (m Mongo) FindAllByFilter(userFilter *model.Filter) ([]model.User, error) {
	filter := bson.M{}
	var results = make([]model.User, 0)

	filter = mountFilter(filter, userFilter)

	cur, err := m.Collection.Find(context.TODO(), filter)

	if err == nil {

		for cur.Next(context.TODO()) {
			var elem model.User
			err := cur.Decode(&elem)
			if err != nil {
				log.Print(err)
			}
			elem.Password = ""
			results = append(results, elem)
		}
	}
	return results, err

}

func mountFilter(filter bson.M, userFilter *model.Filter) bson.M {
	if userFilter != nil {
		if userFilter.Nickname != "" {
			filter["nickname"] = userFilter.Nickname
		}
		if userFilter.LastName != "" {
			filter["lastName"] = userFilter.LastName
		}
		if userFilter.Email != "" {
			filter["email"] = userFilter.Email
		}
		if userFilter.Country != "" {
			filter["country"] = userFilter.Country
		}
	}
	return filter
}
