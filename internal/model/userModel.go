package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID `json:"-" bson:"_id"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Country   string             `json:"country" bson:"country" validate:"required"`
	Nickname  string             `json:"nickname" bson:"nickname" validate:"required"`
	LastName  string             `json:"lastName" bson:"lastName" validate:"required"`
	FirstName string             `json:"firstName" bson:"firstName" validate:"required"`
	Password  string             `json:"password,omitempty" bson:"password" validate:"required"`
}

type ResponseError struct {
	Description string `json:"description"`
}

type Filter struct {
	Email     string `schema:"email"`
	Country   string `schema:"country"`
	Nickname  string `schema:"nickname"`
	LastName  string `schema:"lastName"`
	FirstName string `schema:"firstName"`
}
