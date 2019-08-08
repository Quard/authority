package rest_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Quard/authority/internal/storage"
	"github.com/Quard/authority/internal/user"
	"github.com/thedevsaddam/govalidator"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AuthToken string `json:"auth_token"`
	UserName  string `json:"user_name"`
}

func (srv RestAPIServer) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var requestData loginRequest
	validator := govalidator.New(govalidator.Options{
		Request: r,
		Data:    &requestData,
		Rules: govalidator.MapData{
			"email":    []string{"required", "email"},
			"password": []string{"required", "min:8"},
		},
	})
	validationError := validator.ValidateJSON()
	if len(validationError) > 0 {
		responseValidationError(w, validationError)
	} else {
		user, session, err := login(
			srv.storage,
			requestData.Email,
			requestData.Password,
		)
		if err != nil {
			responseError(w, r, err)
		} else {
			resp := loginResponse{
				AuthToken: session.AuthToken,
				UserName:  user.Name,
			}
			json.NewEncoder(w).Encode(resp)
		}
	}
}

func login(stor storage.Storage, email, password string) (user.User, user.Session, error) {
	var currentUser user.User
	var userSession user.Session

	var err error
	currentUser, err = stor.GetUserByEmail(email)
	if err != nil && err != storage.ErrUserNotFound {
		return currentUser, userSession, errors.New("unable to login")
	}
	if err == storage.ErrUserNotFound {
		return currentUser, userSession, errors.New("such email not registered in application")
	}

	passwd := user.HashPassword(currentUser.Salt, []byte(password))
	if !bytes.Equal(passwd, currentUser.Password) {
		log.Printf("pwd equal: %x & %x", passwd, currentUser.Password)
		return currentUser, userSession, errors.New("wrong email or password")
	}

	userSession, err = currentUser.CreateSession()
	if err != nil {
		return currentUser, userSession, errors.New("unable to login")
	}
	if err := stor.AddSession(currentUser, userSession); err != nil {
		return currentUser, userSession, errors.New("unable to login")
	}

	return currentUser, userSession, nil
}
