package main

import (
	"fmt"
	"github.com/bernardoms/user-api/config"
	_ "github.com/bernardoms/user-api/docs"
	"github.com/bernardoms/user-api/internal/handler"
	"github.com/bernardoms/user-api/internal/logger"
	"github.com/bernardoms/user-api/internal/repository"
	"github.com/gorilla/mux"
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgorilla/v1"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
)

// @title User Swagger API
// @version 1.0
// @description Swagger API for Golang Project User microservice api.
// @termsOfService http://swagger.io/terms/
// @BasePath /v1
func main() {
	app, errNewRelic := newrelic.NewApplication(
		newrelic.NewConfig(os.Getenv("NEWRELIC_APP"), os.Getenv("NEWRELIC_LICENSE")),
	)

	if errNewRelic != nil {
		log.Print("Error in new relic agent")
	}

	logging := logger.ConfigureLogger()

	snsHandler := handler.NewSNS(config.NewSnsConfig(), logging)

	userHandler := handler.UserHandler{
		Repository:    initUserMongoCollection(),
		NotifyHandler: snsHandler,
		Logger:        logging}

	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/v1/users", userHandler.GetAllUsers).Methods("GET")
	r.HandleFunc("/v1/users", userHandler.SaveUser).Methods("POST")
	r.HandleFunc("/v1/users/{nickname}", userHandler.GetUserByNickname).Methods("GET")
	r.HandleFunc("/v1/users/{nickname}", userHandler.DeleteUserByNickname).Methods("DELETE")
	r.HandleFunc("/v1/users/{nickname}", userHandler.UpdateUser).Methods("PUT")

	nrgorilla.InstrumentRoutes(r, app)

	fmt.Printf("running server on %d", 8080)

	err := http.ListenAndServe(":8080", r)

	if err != nil {
		fmt.Printf("error to open port %s with error %s", "8080", err)
	}

}

func initUserMongoCollection() repository.Mongo {
	c := config.NewMongoConfig()
	repository.New(c)
	mongo := repository.Mongo{Collection: repository.GetUserCollection(c)}
	return mongo
}
