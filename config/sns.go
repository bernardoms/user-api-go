package config

import "os"

type SnsConfig struct {
	Topic    string
	Region   string
	Endpoint string
}

func NewSnsConfig() *SnsConfig {
	return &SnsConfig{
		Topic:    os.Getenv("SNS_TOPIC"),
		Region:   os.Getenv("AWS_REGION"),
		Endpoint: os.Getenv("ENDPOINT"),
	}
}
