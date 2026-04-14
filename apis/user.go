package apis

import (
	"context"
	"errors"
	"net/http"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/anvil"
	"github.com/nanoteck137/refinery/core"
	"github.com/nanoteck137/refinery/database"
	"github.com/nanoteck137/refinery/service"
	"github.com/nanoteck137/refinery/types"
	"github.com/nanoteck137/validate"
)

type UserData struct {
	Id string `json:"id"`

	DisplayName string `json:"displayName"`
	Role        string `json:"role"`

	Picture types.Images `json:"picture"`

	Created string `json:"created"`
}

func ConvertDBUser(c pyrin.Context, user database.User) UserData {
	return UserData{
		Id:          user.Id,
		DisplayName: user.DisplayName,
		Role:        user.Role,
		Picture:     ConvertUserPictureURL(c, user.Id),
		Created:     formatTime(user.Created),
	}
}

type GetUser struct {
	User UserData `json:"user"`
}

type UpdateMeBody struct {
	DisplayName *string `json:"displayName,omitempty"`
}

func (b *UpdateMeBody) Transform() {
	b.DisplayName = anvil.StringPtr(b.DisplayName)
}

func (b UpdateMeBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.DisplayName, validate.Required.When(b.DisplayName != nil)),
	)
}

type ApiToken struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	Created string `json:"created"`
	Updated string `json:"updated"`
}

type GetApiTokens struct {
	Tokens []ApiToken `json:"tokens"`
}

type CreateApiToken struct {
	Token string `json:"token"`
}

type CreateApiTokenBody struct {
	Name string `json:"name"`
}

func (b *CreateApiTokenBody) Transform() {
	b.Name = anvil.String(b.Name)
}

func (b CreateApiTokenBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
	)
}

func handleUserServiceErrors(err error) error {
	switch {
	case errors.Is(err, service.ErrUserServiceUserNotFound):
		return UserNotFound()
	case errors.Is(err, service.ErrUserServiceUnauthorized):
		// TODO(patrik): Custom error
		return UserNotFound()
	}

	return err
}

func InstallUserHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetUser",
			Method:       http.MethodGet,
			Path:         "/users/:userId",
			ResponseType: GetUser{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.TODO()

				user, err := app.UserService().GetUserById(
					ctx,
					service.GetUserByIdParams{
						UserId: c.Param("userId"),
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				return GetUser{
					User: ConvertDBUser(c, user),
				}, nil
			},
		},
	)

	group.Register(
		pyrin.ApiHandler{
			Name:     "UpdateMe",
			Method:   http.MethodPatch,
			Path:     "/me",
			BodyType: UpdateMeBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[UpdateMeBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				err = app.UserService().UpdateMe(ctx, service.UpdateMeParams{
					UserId:      user.Id,
					DisplayName: body.DisplayName,
				})
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetApiTokens",
			Method:       http.MethodGet,
			Path:         "/me/apitokens",
			ResponseType: GetApiTokens{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				tokens, err := app.UserService().GetApiTokens(
					ctx,
					service.GetApiTokensParams{
						UserId: user.Id,
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				res := GetApiTokens{
					Tokens: make([]ApiToken, len(tokens)),
				}

				for i, token := range tokens {
					res.Tokens[i] = ApiToken{
						Id:      token.Id,
						Name:    token.Name,
						Created: formatTime(token.Created),
						Updated: formatTime(token.Updated),
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreateApiToken",
			Method:       http.MethodPost,
			Path:         "/me/apitokens",
			ResponseType: CreateApiToken{},
			BodyType:     CreateApiTokenBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[CreateApiTokenBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				tokenId, err := app.UserService().CreateApiToken(
					ctx,
					service.CreateApiTokenParams{
						UserId: user.Id,
						Name:   body.Name,
					},
				)
				if err != nil {
					return nil, err
				}

				return CreateApiToken{
					Token: tokenId,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "DeleteApiToken",
			Method: http.MethodDelete,
			Path:   "/me/apitokens/:tokenId",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				err = app.UserService().DeleteApiToken(
					ctx,
					service.DeleteApiTokenParams{
						TokenId: c.Param("tokenId"),
						UserId:  user.Id,
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				return nil, nil
			},
		},
	)
}
