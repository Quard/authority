package rest_api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/go-pkgz/rest"
	"github.com/thedevsaddam/govalidator"

	"github.com/Quard/authority/internal/storage"
)

type Opts struct {
	bind string
}

type RestAPIServer struct {
	opts    Opts
	storage storage.Storage
}

func NewRestAPIServer(bind string, stor storage.Storage) RestAPIServer {
	opts := Opts{bind: bind}

	srv := RestAPIServer{
		opts:    opts,
		storage: stor,
	}

	srv.initCustomValidators()

	return srv
}

func (srv RestAPIServer) Run() {
	router := srv.getRouter()

	log.Printf("REST API server listen on: %s", srv.opts.bind)
	log.Fatal(http.ListenAndServe(srv.opts.bind, router))
}

func (srv RestAPIServer) getRouter() chi.Router {
	router := chi.NewRouter()
	router.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register/", srv.Register)
		r.Post("/login/", srv.Login)
	})

	return router
}

func (srv RestAPIServer) initCustomValidators() {
	govalidator.AddCustomRule("email_not_registered", func(field, rule, message string, value interface{}) error {
		email := value.(string)
		_, err := srv.storage.GetUserByEmail(email)
		if err == nil {
			return errors.New("this email already registered")
		} else if err != storage.ErrUserNotFound {
			return err
		}

		return nil
	})
}

func responseValidationError(w http.ResponseWriter, validationError url.Values) {
	w.WriteHeader(http.StatusBadRequest)
	err := map[string]interface{}{"validationError": validationError}
	json.NewEncoder(w).Encode(err)
}

func responseError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusBadRequest)
	rest.RenderJSON(w, r, rest.JSON{"error": err.Error()})
}
