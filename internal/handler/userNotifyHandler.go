package handler

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/bernardoms/user-api/config"
	"github.com/bernardoms/user-api/internal/logger"
	"github.com/bernardoms/user-api/internal/model"
)

type Sns struct {
	Topic    string
	Region   string
	Endpoint string
	Logger   *logger.Logger
}

func NewSNS(config *config.SnsConfig, logger *logger.Logger) *Sns {
	newSns := new(Sns)
	newSns.Topic = config.Topic
	newSns.Region = config.Region
	newSns.Endpoint = config.Endpoint
	newSns.Logger = logger
	return newSns
}

func (s Sns) Publish(user *model.User) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String(s.Endpoint),
		Region:   aws.String(s.Region)},
	))

	svc := sns.New(sess)
	message, err := json.Marshal(user)

	if err != nil {
		f := map[string]interface{}{"msg": err}
		s.Logger.LogWithFields(nil, "error", f)
	}

	_, err = svc.Publish(&sns.PublishInput{
		Message:  aws.String(string(message)),
		TopicArn: aws.String(s.Topic),
	})

	if err != nil {
		f := map[string]interface{}{"msg": err}
		s.Logger.LogWithFields(nil, "error", f)
	} else {
		f := map[string]interface{}{"msg": "notifying updated user " + string(message)}
		s.Logger.LogWithFields(nil, "info", f)
	}

	return err
}
