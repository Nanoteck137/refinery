package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/nanoteck137/refinery/database"
	"github.com/nanoteck137/refinery/types"
)

var userErr = NewServiceErrCreator("user")

var (
	ErrUserServiceUserNotFound        = userErr.New("user not found")
	ErrUserServiceApiTokenNotFound    = userErr.New("api token not found")
	ErrUserServiceUnauthorized        = userErr.New("unauthorized")
)

type UserService struct {
	logger *slog.Logger

	db *database.Database

	imageService *ImageService
}

func NewUserService(
	logger *slog.Logger,
	db *database.Database,
	imageService *ImageService,
) *UserService {
	return &UserService{
		logger:       logger,
		db:           db,
		imageService: imageService,
	}
}

type GetUserByIdParams struct {
	UserId string
}

func (s *UserService) GetUserById(
	ctx context.Context,
	params GetUserByIdParams,
) (database.User, error) {
	user, err := s.db.GetUserById(ctx, params.UserId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return database.User{}, ErrUserServiceUserNotFound
		}

		return database.User{}, userErr.Wrap("get user by id", err)
	}

	return user, nil
}

type UpdateMeParams struct {
	UserId string

	DisplayName *string
}

func (s *UserService) UpdateMe(
	ctx context.Context,
	params UpdateMeParams,
) error {
	user, err := s.GetUserById(ctx, GetUserByIdParams{
		UserId: params.UserId,
	})
	if err != nil {
		return userErr.Wrap("update me: get user", err)
	}

	changes := database.UserChanges{}

	if params.DisplayName != nil {
		changes.DisplayName = types.Change[string]{
			Value:   *params.DisplayName,
			Changed: *params.DisplayName != user.DisplayName,
		}
	}

	err = s.db.UpdateUser(ctx, user.Id, changes)
	if err != nil {
		return userErr.Wrap("update me: db update", err)
	}

	return nil
}

type GetApiTokensParams struct {
	UserId string
}

func (s *UserService) GetApiTokens(
	ctx context.Context,
	params GetApiTokensParams,
) ([]database.ApiToken, error) {
	tokens, err := s.db.GetAllApiTokensForUser(ctx, params.UserId)
	if err != nil {
		return nil, userErr.Wrap("get api tokens: db get", err)
	}

	return tokens, nil
}

type CreateApiTokenParams struct {
	UserId string

	Name string
}

func (s *UserService) CreateApiToken(
	ctx context.Context,
	params CreateApiTokenParams,
) (string, error) {
	id, err := s.db.CreateApiToken(ctx, database.CreateApiTokenParams{
		UserId: params.UserId,
		Name:   params.Name,
	})
	if err != nil {
		return "", userErr.Wrap("create api token: db create", err)
	}

	return id, nil
}

type DeleteApiTokenParams struct {
	TokenId string
	UserId  string
}

func (s *UserService) DeleteApiToken(
	ctx context.Context,
	params DeleteApiTokenParams,
) error {
	token, err := s.db.GetApiTokenById(ctx, params.TokenId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return ErrUserServiceApiTokenNotFound
		}

		return userErr.Wrap("delete api token: db get api token", err)
	}

	if token.UserId != params.UserId {
		return ErrUserServiceUnauthorized
	}

	err = s.db.DeleteApiToken(ctx, token.Id)
	if err != nil {
		return userErr.Wrap("delete api token: db delete", err)
	}

	return nil
}
