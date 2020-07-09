package handler

import (
	"encoding/json"
	"github.com/bernardoms/user-api/internal/logger"
	"github.com/bernardoms/user-api/internal/model"
	"github.com/bernardoms/user-api/internal/repository"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
)

var decoder = schema.NewDecoder()

type UserHandler struct {
	Repository    repository.UserRepository
	NotifyHandler NotifyInterface
	Logger        *logger.Logger
}

// GetAllUsers godoc
// @Summary Retrieves user based on a given filter
// @Description Get all users
// @Produce json
// @Param nickname query string false "User nickname"
// @Success 200 {object} model.User
// @Router /users [get]
// @Tags users
func (u *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	filter := new(model.Filter)
	err := decoder.Decode(filter, r.URL.Query())

	if err != nil {
		f := map[string]interface{}{"msg": "Error in GET parameters " + err.Error(), "parameters": r.URL.Query()}
		u.Logger.LogWithFields(r, "error", f)
		log.Println("Error in GET parameters : ", err)
	}

	results, err := u.Repository.FindAllByFilter(filter)

	if err != nil {
		f := map[string]interface{}{"msg": err}
		u.Logger.LogWithFields(r, "error", f)
		respondWithJson(w, http.StatusInternalServerError, model.ResponseError{Description: err.Error()})
		return
	}

	respondWithJson(w, http.StatusOK, results)
}

// GetAllUsers godoc
// @Summary Retrieves an user by a given nickname
// @Description Retrieves an user by a given nickname
// @Produce json
// @Param nickname path string true "User nickname"
// @Success 200 {object} model.User
// @Router /users/{nickname} [get]
// @Tags users
func (u *UserHandler) GetUserByNickname(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	result, err := u.Repository.FindByNickname(vars["nickname"])

	if err != nil {
		f := map[string]interface{}{"msg": err}
		u.Logger.LogWithFields(r, "error", f)
		respondWithJson(w, http.StatusInternalServerError, model.ResponseError{Description: err.Error()})
		return
	}

	if result == nil {
		f := map[string]interface{}{"msg": "user with nickname " + vars["nickname"] + " not found!", "nickname": vars["nickname"]}
		u.Logger.LogWithFields(r, "info", f)
		respondWithJson(w, http.StatusNotFound, model.ResponseError{Description: "user with nickname " + vars["nickname"] + " not found!"})
		return
	}
	respondWithJson(w, http.StatusOK, result)
}

// DeleteUserByNickname godoc
// @Summary Deletes an user by a given nickname
// @Description Deletes an user by a given nickname
// @Produce json
// @Param nickname path string true "User nickname"
// @Success 204
// @Router /users/{nickname} [delete]
// @Tags users
func (u *UserHandler) DeleteUserByNickname(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := u.Repository.Delete(vars["nickname"])

	if err != nil {
		f := map[string]interface{}{"msg": err}
		u.Logger.LogWithFields(r, "error", f)
		respondWithJson(w, http.StatusInternalServerError, model.ResponseError{Description: err.Error()})
		return
	}

	respondWithEmpty(w, http.StatusNoContent, "")
}

// SaveUser godoc
// @Summary create an user
// @Description create an user
// @Produce json
// @Param user body model.User true "Create user"
// @Success 201
// @Header 201 {string} Location "/v1/users/{nickname}
// @Router /users [post]
// @Tags users
func (u *UserHandler) SaveUser(w http.ResponseWriter, r *http.Request) {
	v := validator.New()

	var user *model.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		f := map[string]interface{}{"msg": err}
		u.Logger.LogWithFields(r, "error", f)
		respondWithJson(w, http.StatusInternalServerError, model.ResponseError{Description: err.Error()})
		return
	}

	err = v.Struct(user)

	if err != nil {
		f := map[string]interface{}{"msg": err}
		u.Logger.LogWithFields(r, "info", f)
		respondWithJson(w, http.StatusBadRequest, model.ResponseError{Description: err.Error()})
		return
	}

	result, err := u.Repository.FindByNickname(user.Nickname)

	if err != nil {
		f := map[string]interface{}{"msg": err}
		u.Logger.LogWithFields(r, "error", f)
		respondWithJson(w, http.StatusInternalServerError, model.ResponseError{Description: err.Error()})
	}

	if result != nil {
		f := map[string]interface{}{"msg": "user with nick name " + user.Nickname + " already exist!"}
		u.Logger.LogWithFields(r, "info", f)
		respondWithJson(w, http.StatusConflict, model.ResponseError{Description: "user with nick name " + user.Nickname + " already exist!"})
		return
	}

	user.Password = hashAndSalt([]byte(user.Password))
	user.Id = primitive.NewObjectID()

	inserted, err := u.Repository.Save(user)

	if err != nil {

		respondWithJson(w, http.StatusInternalServerError, model.ResponseError{Description: err.Error()})
		return
	}
	respondWithEmpty(w, http.StatusCreated, "v1/users/"+inserted.Nickname)
}

// UpdateUser godoc
// @Summary Update an user by a given nickname and notify to a topic
// @Description Update an user by a given nickname and notify to a topic
// @Produce json
// @Param nickname path string true "User nickname"
// @Success 204
// @Router /users/{nickname} [put]
// @Tags users
func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	vars := mux.Vars(r)

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		f := map[string]interface{}{"msg": err}
		u.Logger.LogWithFields(r, "error", f)
		respondWithJson(w, http.StatusInternalServerError, model.ResponseError{Description: err.Error()})
		return
	}

	user.Password = hashAndSalt([]byte(user.Password))

	result, err := u.Repository.UpdateByNickname(vars["nickname"], user)

	if err != nil {
		f := map[string]interface{}{"msg": err}
		u.Logger.LogWithFields(r, "error", f)
		respondWithJson(w, http.StatusInternalServerError, model.ResponseError{Description: err.Error()})
		return
	}

	if result > 0 {
		err := u.NotifyHandler.Publish(user)
		if err != nil {
			f := map[string]interface{}{"msg": err}
			u.Logger.LogWithFields(r, "error", f)
			respondWithJson(w, http.StatusInternalServerError, model.ResponseError{Description: err.Error()})
			return
		}
	}
	respondWithEmpty(w, http.StatusNoContent, "")
}

func respondWithEmpty(w http.ResponseWriter, code int, location string) {
	w.Header().Set("Content-Type", "application/json")
	if location != "" {
		w.Header().Set("Location", location)
	}
	w.WriteHeader(code)
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
