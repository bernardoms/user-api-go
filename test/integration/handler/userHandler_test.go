package handler

import (
	"bytes"
	"github.com/bernardoms/user-api/config"
	"github.com/bernardoms/user-api/internal/handler"
	"github.com/bernardoms/user-api/internal/logger"
	"github.com/bernardoms/user-api/internal/repository"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllUsersSuccessNoFilter(t *testing.T) {

	c := config.NewMongoConfig()
	repository.New(c)

	mongo := &repository.Mongo{Collection: repository.GetUserCollection(c)}
	sns := handler.NewSNS(config.NewSnsConfig(), logger.ConfigureLogger())

	h := handler.UserHandler{Repository: mongo, NotifyHandler: sns, Logger: logger.ConfigureLogger()}

	r, _ := http.NewRequest("GET", "/v1/users", nil)
	w := httptest.NewRecorder()

	h.GetAllUsers(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[{\"email\":\"test2@test.com\",\"country\":\"UK\",\"nickname\":\"testnick2\",\"lastName\":\"lastName2\",\"firstName\":\"firstName2\"},{\"email\":\"test1@test.com\",\"country\":\"UK\",\"nickname\":\"testnick1\",\"lastName\":\"lastName\",\"firstName\":\"firstName\"}]", w.Body.String())
}

func TestGetAllUsersSuccessFilters(t *testing.T) {

	c := config.NewMongoConfig()
	repository.New(c)

	mongo := &repository.Mongo{Collection: repository.GetUserCollection(c)}
	sns := handler.NewSNS(config.NewSnsConfig(), logger.ConfigureLogger())

	h := handler.UserHandler{Repository: mongo, NotifyHandler: sns, Logger: logger.ConfigureLogger()}

	r, _ := http.NewRequest("GET", "/v1/users?nickname=testnick1&firstName=firstName&lastName=lastName&country=UK&email=test1@test.com", nil)

	w := httptest.NewRecorder()

	h.GetAllUsers(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[{\"email\":\"test1@test.com\",\"country\":\"UK\",\"nickname\":\"testnick1\",\"lastName\":\"lastName\",\"firstName\":\"firstName\"}]", w.Body.String())
}

func TestGetUserByNickNameSuccess(t *testing.T) {

	c := config.NewMongoConfig()
	repository.New(c)

	mongo := &repository.Mongo{Collection: repository.GetUserCollection(c)}
	sns := handler.NewSNS(config.NewSnsConfig(), logger.ConfigureLogger())

	h := handler.UserHandler{Repository: mongo, NotifyHandler: sns, Logger: logger.ConfigureLogger()}

	r, _ := http.NewRequest("GET", "/v1/users", nil)

	vars := map[string]string{
		"nickname": "testnick1",
	}

	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	h.GetUserByNickname(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"email\":\"test1@test.com\",\"country\":\"UK\",\"nickname\":\"testnick1\",\"lastName\":\"lastName\",\"firstName\":\"firstName\"}", w.Body.String())
}

func TestGetUserByNickNameNotFound(t *testing.T) {

	c := config.NewMongoConfig()
	repository.New(c)

	mongo := &repository.Mongo{Collection: repository.GetUserCollection(c)}
	sns := handler.NewSNS(config.NewSnsConfig(), logger.ConfigureLogger())

	h := handler.UserHandler{Repository: mongo, NotifyHandler: sns, Logger: logger.ConfigureLogger()}

	r, _ := http.NewRequest("GET", "/v1/users", nil)

	vars := map[string]string{
		"nickname": "testnickname",
	}

	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	h.GetUserByNickname(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{\"description\":\"user with nickname testnickname not found!\"}", w.Body.String())
}

func TestDeleteUserByNickName(t *testing.T) {

	c := config.NewMongoConfig()
	repository.New(c)

	mongo := &repository.Mongo{Collection: repository.GetUserCollection(c)}
	sns := handler.NewSNS(config.NewSnsConfig(), logger.ConfigureLogger())

	h := handler.UserHandler{Repository: mongo, NotifyHandler: sns, Logger: logger.ConfigureLogger()}

	r, _ := http.NewRequest("DELETE", "/v1/users", nil)

	vars := map[string]string{
		"nickname": "testnickname",
	}

	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	h.DeleteUserByNickname(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestSaveUserSuccess(t *testing.T) {

	c := config.NewMongoConfig()
	repository.New(c)

	mongo := &repository.Mongo{Collection: repository.GetUserCollection(c)}
	sns := handler.NewSNS(config.NewSnsConfig(), logger.ConfigureLogger())

	h := handler.UserHandler{Repository: mongo, NotifyHandler: sns, Logger: logger.ConfigureLogger()}

	jsonStr := []byte(`{"email":"test4@test.com", "country" : "BR", "lastName" : "lastName", "firstName":"firstName", "password":"password", "nickname": "testnickname3"}`)

	r, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	h.SaveUser(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "", w.Body.String())
	assert.Equal(t, "v1/users/testnickname3", w.Header().Get("Location"))

	//clean database
	mongo.Delete("testnickname3")
}

func TestSaveValidationError(t *testing.T) {
	c := config.NewMongoConfig()
	repository.New(c)

	mongo := &repository.Mongo{Collection: repository.GetUserCollection(c)}
	sns := handler.NewSNS(config.NewSnsConfig(), logger.ConfigureLogger())

	h := handler.UserHandler{Repository: mongo, NotifyHandler: sns, Logger: logger.ConfigureLogger()}

	jsonStr := []byte(`{}`)

	r, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()

	h.SaveUser(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"description\":\"Key: 'User.Email' Error:Field validation for 'Email' failed on the 'required' tag\\nKey: 'User.Country' Error:Field validation for 'Country' failed on the 'required' tag\\nKey: 'User.Nickname' Error:Field validation for 'Nickname' failed on the 'required' tag\\nKey: 'User.LastName' Error:Field validation for 'LastName' failed on the 'required' tag\\nKey: 'User.FirstName' Error:Field validation for 'FirstName' failed on the 'required' tag\\nKey: 'User.Password' Error:Field validation for 'Password' failed on the 'required' tag\"}", w.Body.String())
}

//func TestSaveUserNickNameAlreadyExists(t *testing.T) {
//
//	c := config.NewMongoConfig()
//	repository.New(c)
//
//	mongo := &repository.Mongo{Collection: repository.GetUserCollection(c)}
//	sns := handler.NewSNS(config.NewSnsConfig(), logger.ConfigureLogger())
//
//	h := handler.UserHandler{Repository: mongo, NotifyHandler: sns, Logger: logger.ConfigureLogger()}
//
//	jsonStr := []byte(`{"email":"test@test.com", "country" : "UK", "lastName" : "lastName", "firstName":"firstName", "password":"password", "nickname": "testnick1"}`)
//
//	r, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonStr))
//
//	w := httptest.NewRecorder()
//
//	h.SaveUser(w, r)
//
//	assert.Equal(t, http.StatusConflict, w.Code)
//	assert.Equal(t,  "{\"description\":\"user with nick name testnick1 already exist!\"}", w.Body.String())
//}

func TestUpdateUserSuccess(t *testing.T) {

	c := config.NewMongoConfig()
	repository.New(c)

	mongo := &repository.Mongo{Collection: repository.GetUserCollection(c)}
	sns := handler.NewSNS(config.NewSnsConfig(), logger.ConfigureLogger())

	h := handler.UserHandler{Repository: mongo, NotifyHandler: sns, Logger: logger.ConfigureLogger()}

	jsonStr := []byte(`{"email":"test1@test.com", "country" : "UK", "lastName" : "lastName", "firstName":"firstName", "password":"password2", "nickname": "testnickname1"}`)

	r, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonStr))

	vars := map[string]string{
		"nickname": "testnickname1",
	}

	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	h.UpdateUser(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "", w.Body.String())
}
