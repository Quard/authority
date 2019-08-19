package internal_api

import (
	context "context"

	"github.com/getsentry/sentry-go"

	"github.com/Quard/authority/internal/user"
)

func (srv internalAPIServer) GetUserByAuthToken(ctx context.Context, authToken *AuthToken) (*User, error) {
	authUser, err := srv.storage.GetUserBySession(authToken.GetAuthToken())
	if err != nil {
		return nil, err
	}
	return &User{ID: authUser.ID.Hex(), Email: authUser.Email, Name: authUser.Name}, nil
}

func (srv internalAPIServer) GetUserByProp(ctx context.Context, prop *UserProp) (*User, error) {
	authUser, err := srv.storage.GetUserByProp(prop.GetName(), prop.GetValue())
	if err != nil {
		return nil, err
	}

	return &User{ID: authUser.ID.Hex(), Email: authUser.Email, Name: authUser.Name}, nil
}

func (srv internalAPIServer) SetUserProps(ctx context.Context, userProps *UserProps) (*Error, error) {
	var currentUser user.User
	var err error

	switch userProps.User.(type) {
	case *UserProps_Id:
		currentUser, err = srv.storage.GetUserByID(userProps.GetId())
		if err != nil {
			sentry.CaptureException(err)
			return &Error{IsOk: false, Msg: err.Error()}, err
		}
	case *UserProps_AuthToken:
		currentUser, err = srv.storage.GetUserBySession(userProps.GetAuthToken())
		if err != nil {
			sentry.CaptureException(err)
			return &Error{IsOk: false, Msg: err.Error()}, err
		}
	}

	for _, prop := range userProps.GetProps() {
		err := srv.storage.SetUserProp(currentUser, prop.GetName(), prop.GetValue())
		if err != nil {
			return &Error{IsOk: false, Msg: err.Error()}, err
		}
	}

	return &Error{IsOk: true}, nil
}
