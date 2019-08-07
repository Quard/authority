package rest_api

import (
	"net/http"
	"time"

	"github.com/thedevsaddam/govalidator"

	"github.com/Quard/authority/internal/storage"
	"github.com/Quard/authority/internal/user"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (srv RestAPIServer) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var requestData registerRequest
	validator := govalidator.New(govalidator.Options{
		Request: r,
		Data:    &requestData,
		Rules: govalidator.MapData{
			"email":    []string{"required", "email", "email_not_registered"},
			"password": []string{"required", "min:8"},
		},
	})
	validationError := validator.ValidateJSON()
	if len(validationError) > 0 {
		responseValidationError(w, validationError)
	} else {
		err := registerUser(
			srv.storage,
			requestData.Email,
			requestData.Password,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func registerUser(storage storage.Storage, email, password string) error {
	user := user.User{Email: email, DateJoined: time.Now()}
	if err := user.SetPassword(password); err != nil {
		return err
	}

	return storage.AddUser(user)
}
