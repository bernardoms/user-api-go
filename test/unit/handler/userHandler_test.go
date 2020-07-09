package handler

import (
	"bytes"
	"errors"
	"github.com/bernardoms/user-api/internal/handler"
	"github.com/bernardoms/user-api/internal/logger"
	"github.com/bernardoms/user-api/internal/model"
	"github.com/bernardoms/user-api/test/unit/mock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllUsersSuccessNoFilter(t *testing.T) {

	mongoMock := mock.MongoMock{}

	snsMock := mock.NotifyMock{}

	h := handler.UserHandler{Repository: &mongoMock, NotifyHandler: &snsMock, Logger: logger.ConfigureLogger()}

	r, _ := http.NewRequest("GET", "/v1/users", nil)
	w := httptest.NewRecorder()

	id, _ := primitive.ObjectIDFromHex("5ea7208049e00ddb76994ede")
	id2, _ := primitive.ObjectIDFromHex("5ea7208049e00ddb76994eda")

	var users = []model.User{
		{Id: id, Nickname: "test1", Password: "password", LastName: "lastName", FirstName: "firstName", Country: "UK", Email: "test@test.com"},
		{Id: id2, Nickname: "test1", Password: "password2", LastName: "lastName2", FirstName: "firstName2", Country: "UK", Email: "test2@test.com"},
	}

	mongoMock.On("FindAllByFilter", &model.Filter{}).Return(users, nil)
	h.GetAllUsers(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[{\"email\":\"test@test.com\",\"country\":\"UK\",\"nickname\":\"test1\",\"lastName\":\"lastName\",\"firstName\":\"firstName\",\"password\":\"password\"},{\"email\":\"test2@test.com\",\"country\":\"UK\",\"nickname\":\"test1\",\"lastName\":\"lastName2\",\"firstName\":\"firstName2\",\"password\":\"password2\"}]", w.Body.String())
}

func TestGetAllUsersSuccessFilters(t *testing.T) {

	mongoMock := mock.MongoMock{}

	snsMock := mock.NotifyMock{}

	filter := new(model.Filter)

	filter.Nickname = "test1"
	filter.FirstName = "firstName"
	filter.LastName = "lastName"
	filter.Country = "UK"
	filter.Email = "test@test.com"

	h := handler.UserHandler{Repository: &mongoMock, NotifyHandler: &snsMock, Logger: logger.ConfigureLogger()}

	r, _ := http.NewRequest("GET", "/v1/users?nickname=test1&firstName=firstName&lastName=lastName&country=UK&email=test@test.com", nil)

	w := httptest.NewRecorder()

	id, _ := primitive.ObjectIDFromHex("5ea7208049e00ddb76994ede")

	var users = []model.User{
		{Id: id, Nickname: "test1", Password: "password", LastName: "lastName", FirstName: "firstName", Country: "UK", Email: "test@test.com"},
	}

	mongoMock.On("FindAllByFilter", filter).Return(users, nil)
	h.GetAllUsers(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[{\"email\":\"test@test.com\",\"country\":\"UK\",\"nickname\":\"test1\",\"lastName\":\"lastName\",\"firstName\":\"firstName\",\"password\":\"password\"}]", w.Body.String())
}

func TestGetAllUsersErrorOnMongo(t *testing.T) {

	mongoMock := mock.MongoMock{}

	snsMock := mock.NotifyMock{}

	filter := new(model.Filter)

	filter.Nickname = "test1"
	filter.FirstName = "firstName"
	filter.LastName = "lastName"
	filter.Country = "UK"
	filter.Email = "test@test.com"

	h := handler.UserHandler{Repository: &mongoMock, NotifyHandler: &snsMock, Logger: logger.ConfigureLogger()}

	r, _ := http.NewRequest("GET", "/v1/users?nickname=test1&firstName=firstName&lastName=lastName&country=UK&email=test@test.com", nil)

	w := httptest.NewRecorder()

	id, _ := primitive.ObjectIDFromHex("5ea7208049e00ddb76994ede")

	var users = []model.User{
		{Id: id, Nickname: "test1", Password: "password", LastName: "lastName", FirstName: "firstName", Country: "UK", Email: "test@test.com"},
	}

	mongoMock.On("FindAllByFilter", filter).Return(users, errors.New("error on mongo"))
	h.GetAllUsers(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"description\":\"error on mongo\"}", w.Body.String())
}

func TestGetUserByNickNameSuccess(t *testing.T) {

	mongoMock := mock.MongoMock{}

	snsMock := mock.NotifyMock{}

	user := new(model.User)

	user.Email = "test@test.com"
	user.Country = "UK"
	user.LastName = "lastName"
	user.FirstName = "firstName"
	user.Password = "password"

	h := handler.UserHandler{Repository: &mongoMock, NotifyHandler: &snsMock, Logger: logger.ConfigureLogger()}

	r, _ := http.NewRequest("GET", "/v1/users", nil)

	vars := map[string]string{
		"nickname": "testnickname",
	}

	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	mongoMock.On("FindByNickname", "testnickname").Return(user, nil)
	h.GetUserByNickname(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"email\":\"test@test.com\",\"country\":\"UK\",\"nickname\":\"\",\"lastName\":\"lastName\",\"firstName\":\"firstName\",\"password\":\"password\"}", w.Body.String())
}

func TestGetUserByNickNameErrorOnMongo(t *testing.T) {

	mongoMock := mock.MongoMock{}

	snsMock := mock.NotifyMock{}

	user := new(model.User)

	user.Email = "test@test.com"
	user.Country = "UK"
	user.LastName = "lastName"
	user.FirstName = "firstName"
	user.Password = "password"

	h := handler.UserHandler{Repository: &mongoMock, NotifyHandler: &snsMock, Logger: logger.ConfigureLogger()}

	r, _ := http.NewRequest("GET", "/v1/users", nil)

	vars := map[string]string{
		"nickname": "testnickname",
	}

	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	mongoMock.On("FindByNickname", "testnickname").Return(user, errors.New("error on mongo"))
	h.GetUserByNickname(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"description\":\"error on mongo\"}", w.Body.String())
}

func TestGetUserByNickNameNotFound(t *testing.T) {

	mongoMock := mock.MongoMock{}

	snsMock := mock.NotifyMock{}

	var user *model.User

	h := handler.UserHandler{Repository: &mongoMock, NotifyHandler: &snsMock, Logger: logger.ConfigureLogger()}

	r, _ := http.NewRequest("GET", "/v1/users", nil)

	vars := map[string]string{
		"nickname": "testnickname",
	}

	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	mongoMock.On("FindByNickname", "testnickname").Return(user, nil)
	h.GetUserByNickname(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{\"description\":\"user with nickname testnickname not found!\"}", w.Body.String())
}

func TestDeleteUserByNickName(t *testing.T) {

	mongoMock := mock.MongoMock{}

	snsMock := mock.NotifyMock{}

	h := handler.UserHandler{Repository: &mongoMock, NotifyHandler: &snsMock, Logger: logger.ConfigureLogger()}

	r, _ := http.NewRequest("DELETE", "/v1/users", nil)

	vars := map[string]string{
		"nickname": "testnickname",
	}

	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	mongoMock.On("Delete", "testnickname").Return(nil)
	h.DeleteUserByNickname(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestDeleteUserByNickNameErrorOnMongo(t *testing.T) {

	mongoMock := mock.MongoMock{}

	snsMock := mock.NotifyMock{}

	h := handler.UserHandler{Repository: &mongoMock, NotifyHandler: &snsMock, Logger: logger.ConfigureLogger()}

	r, _ := http.NewRequest("DELETE", "/v1/users", nil)

	vars := map[string]string{
		"nickname": "testnickname",
	}

	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	mongoMock.On("Delete", "testnickname").Return(errors.New("error on mongo"))
	h.DeleteUserByNickname(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"description\":\"error on mongo\"}", w.Body.String())
}

func TestSaveUserSuccess(t *testing.T) {

	mongoMock := mock.MongoMock{}

	snsMock := mock.NotifyMock{}

	var notFoundNick *model.User

	user := new(model.User)
	user.Email = "test@test.com"
	user.Country = "UK"
	user.LastName = "lastName"
	user.FirstName = "firstName"
	user.Password = "password"
	user.Nickname = "testnickname"

	h := handler.UserHandler{Repository: &mongoMock, NotifyHandler: &snsMock, Logger: logger.ConfigureLogger()}

	jsonStr := []byte(`{"email":"test@test.com", "country" : "UK", "lastName" : "lastName", "firstName":"firstName", "password":"password", "nickname": "testnickname"}`)

	r, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()

	mongoMock.On("FindByNickname", "testnickname").Return(notFoundNick, nil)
	mongoMock.On("Save", mock2.Anything).Return(user, nil)
	h.SaveUser(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "", w.Body.String())
	assert.Equal(t, "v1/users/testnickname", w.Header().Get("Location"))
}

func TestSaveValidationError(t *testing.T) {

	mongoMock := mock.MongoMock{}

	snsMock := mock.NotifyMock{}

	var notFoundNick *model.User

	user := new(model.User)
	user.Email = "test@test.com"
	user.Country = "UK"
	user.LastName = "lastName"
	user.FirstName = "firstName"
	user.Password = "password"
	user.Nickname = "testnickname"

	h := handler.UserHandler{Repository: &mongoMock, NotifyHandler: &snsMock, Logger: logger.ConfigureLogger()}

	jsonStr := []byte(`{}`)

	r, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()

	mongoMock.On("FindByNickname", "testnickname").Return(notFoundNick, nil)
	mongoMock.On("Save", mock2.Anything).Return(user, nil)
	h.SaveUser(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"description\":\"Key: 'User.Email' Error:Field validation for 'Email' failed on the 'required' tag\\nKey: 'User.Country' Error:Field validation for 'Country' failed on the 'required' tag\\nKey: 'User.Nickname' Error:Field validation for 'Nickname' failed on the 'required' tag\\nKey: 'User.LastName' Error:Field validation for 'LastName' failed on the 'required' tag\\nKey: 'User.FirstName' Error:Field validation for 'FirstName' failed on the 'required' tag\\nKey: 'User.Password' Error:Field validation for 'Password' failed on the 'required' tag\"}", w.Body.String())
}

func TestSaveUserErrorOnFindNickName(t *testing.T) {

	mongoMock := mock.MongoMock{}

	snsMock := mock.NotifyMock{}

	var notFoundNick *model.User

	user := new(model.User)
	user.Email = "test@test.com"
	user.Country = "UK"
	user.LastName = "lastName"
	user.FirstName = "firstName"
	user.Password = "password"
	user.Nickname = "testnickname"

	h := handler.UserHandler{Repository: &mongoMock, NotifyHandler: &snsMock, Logger: logger.ConfigureLogger()}

	jsonStr := []byte(`{"email":"test@test.com", "country" : "UK", "lastName" : "lastName", "firstName":"firstName", "password":"password", "nickname": "testnickname"}`)

	r, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()

	mongoMock.On("FindByNickname", "testnickname").Return(notFoundNick, errors.New("error on mongo"))
	mongoMock.On("Save", mock2.Anything).Return(user, nil)
	h.SaveUser(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"description\":\"error on mongo\"}", w.Body.String())
}

func TestSaveUserErrorOnSave(t *testing.T) {

	mongoMock := mock.MongoMock{}

	snsMock := mock.NotifyMock{}

	var notFoundNick *model.User

	user := new(model.User)
	user.Email = "test@test.com"
	user.Country = "UK"
	user.LastName = "lastName"
	user.FirstName = "firstName"
	user.Password = "password"
	user.Nickname = "testnickname"

	h := handler.UserHandler{Repository: &mongoMock, NotifyHandler: &snsMock, Logger: logger.ConfigureLogger()}

	jsonStr := []byte(`{"email":"test@test.com", "country" : "UK", "lastName" : "lastName", "firstName":"firstName", "password":"password", "nickname": "testnickname"}`)

	r, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()

	mongoMock.On("FindByNickname", "testnickname").Return(notFoundNick, nil)
	mongoMock.On("Save", mock2.Anything).Return(user, errors.New("error on mongo"))
	h.SaveUser(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"description\":\"error on mongo\"}", w.Body.String())
}

func TestSaveUserNickNameAlreadyExists(t *testing.T) {

	mongoMock := mock.MongoMock{}

	snsMock := mock.NotifyMock{}

	user := new(model.User)
	user.Email = "test@test.com"
	user.Country = "UK"
	user.LastName = "lastName"
	user.FirstName = "firstName"
	user.Password = "password"
	user.Nickname = "testnickname"

	h := handler.UserHandler{Repository: &mongoMock, NotifyHandler: &snsMock, Logger: logger.ConfigureLogger()}

	jsonStr := []byte(`{"email":"test@test.com", "country" : "UK", "lastName" : "lastName", "firstName":"firstName", "password":"password", "nickname": "testnickname"}`)

	r, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()

	mongoMock.On("FindByNickname", "testnickname").Return(user, nil)
	mongoMock.On("Save", mock2.Anything).Return(user, nil)
	h.SaveUser(w, r)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, "{\"description\":\"user with nick name testnickname already exist!\"}", w.Body.String())
}

func TestUpdateUserSuccess(t *testing.T) {

	mongoMock := mock.MongoMock{}

	snsMock := mock.NotifyMock{}

	user := new(model.User)
	user.Email = "test@test.com"
	user.Country = "UK"
	user.LastName = "lastName"
	user.FirstName = "firstName"
	user.Password = "password"
	user.Nickname = "testnickname"

	h := handler.UserHandler{Repository: &mongoMock, NotifyHandler: &snsMock, Logger: logger.ConfigureLogger()}

	jsonStr := []byte(`{"email":"test@test.com", "country" : "UK", "lastName" : "lastName", "firstName":"firstName", "password":"password", "nickname": "testnickname"}`)

	r, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonStr))

	vars := map[string]string{
		"nickname": "testnickname",
	}

	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	mongoMock.On("UpdateByNickname", "testnickname", mock2.Anything).Return(int64(1), nil)
	snsMock.On("Publish", mock2.Anything).Return(nil)
	h.UpdateUser(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestUpdateUserMongoErrorOnUpdate(t *testing.T) {

	mongoMock := mock.MongoMock{}

	snsMock := mock.NotifyMock{}

	user := new(model.User)
	user.Email = "test@test.com"
	user.Country = "UK"
	user.LastName = "lastName"
	user.FirstName = "firstName"
	user.Password = "password"
	user.Nickname = "testnickname"

	h := handler.UserHandler{Repository: &mongoMock, NotifyHandler: &snsMock, Logger: logger.ConfigureLogger()}

	jsonStr := []byte(`{"email":"test@test.com", "country" : "UK", "lastName" : "lastName", "firstName":"firstName", "password":"password", "nickname": "testnickname"}`)

	r, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonStr))

	vars := map[string]string{
		"nickname": "testnickname",
	}

	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	mongoMock.On("UpdateByNickname", "testnickname", mock2.Anything).Return(int64(0), errors.New("mongo error"))
	snsMock.On("Publish", mock2.Anything).Return(nil)
	h.UpdateUser(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"description\":\"mongo error\"}", w.Body.String())
}

func TestUpdateUserNotifyError(t *testing.T) {

	mongoMock := mock.MongoMock{}

	snsMock := mock.NotifyMock{}

	user := new(model.User)
	user.Email = "test@test.com"
	user.Country = "UK"
	user.LastName = "lastName"
	user.FirstName = "firstName"
	user.Password = "password"
	user.Nickname = "testnickname"

	h := handler.UserHandler{Repository: &mongoMock, NotifyHandler: &snsMock, Logger: logger.ConfigureLogger()}

	jsonStr := []byte(`{"email":"test@test.com", "country" : "UK", "lastName" : "lastName", "firstName":"firstName", "password":"password", "nickname": "testnickname"}`)

	r, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonStr))

	vars := map[string]string{
		"nickname": "testnickname",
	}

	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	mongoMock.On("UpdateByNickname", "testnickname", mock2.Anything).Return(int64(1), nil)
	snsMock.On("Publish", mock2.Anything).Return(errors.New("error on notify"))
	h.UpdateUser(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"description\":\"error on notify\"}", w.Body.String())
}

func TestUpdateUserSuccessWithoutNotify(t *testing.T) {

	mongoMock := mock.MongoMock{}

	snsMock := mock.NotifyMock{}

	user := new(model.User)
	user.Email = "test@test.com"
	user.Country = "UK"
	user.LastName = "lastName"
	user.FirstName = "firstName"
	user.Password = "password"
	user.Nickname = "testnickname"

	h := handler.UserHandler{Repository: &mongoMock, NotifyHandler: &snsMock, Logger: logger.ConfigureLogger()}

	jsonStr := []byte(`{"email":"test@test.com", "country" : "UK", "lastName" : "lastName", "firstName":"firstName", "password":"password", "nickname": "testnickname"}`)

	r, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonStr))

	vars := map[string]string{
		"nickname": "testnickname",
	}

	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	mongoMock.On("UpdateByNickname", "testnickname", mock2.Anything).Return(int64(0), nil)
	snsMock.AssertNotCalled(t, mock2.Anything)
	h.UpdateUser(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "", w.Body.String())
}
