package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/refinery/tools/utils"
)

type ApiToken struct {
	Id     string `db:"id"`
	UserId string `db:"user_id"`

	Name string `db:"name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func ApiTokenQuery() *goqu.SelectDataset {
	query := dialect.From("api_tokens").
		Select(
			"api_tokens.id",
			"api_tokens.user_id",

			"api_tokens.name",

			"api_tokens.updated",
			"api_tokens.created",
		).
		Prepared(true)

	return query
}

func (db DB) GetApiTokenById(
	ctx context.Context,
	tokenId string,
) (ApiToken, error) {
	query := ApiTokenQuery().
		Where(goqu.I("api_tokens.id").Eq(tokenId))

	return ember.Single[ApiToken](db.db, ctx, query)
}

func (db DB) GetAllApiTokensForUser(
	ctx context.Context,
	userId string,
) ([]ApiToken, error) {
	query := ApiTokenQuery().
		Where(goqu.I("api_tokens.user_id").Eq(userId))

	return ember.Multiple[ApiToken](db.db, ctx, query)
}

type CreateApiTokenParams struct {
	Id     string
	UserId string
	Name   string

	Created int64
	Updated int64
}

func (db DB) CreateApiToken(
	ctx context.Context,
	params CreateApiTokenParams,
) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = utils.CreateApiTokenId()
	}

	query := dialect.Insert("api_tokens").Rows(goqu.Record{
		"id":      params.Id,
		"user_id": params.UserId,

		"name": params.Name,

		"created": params.Created,
		"updated": params.Updated,
	})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return "", err
	}

	return params.Id, nil
}

func (db DB) DeleteApiToken(ctx context.Context, tokenId string) error {
	query := dialect.Delete("api_tokens").
		Where(goqu.I("api_tokens.id").Eq(tokenId))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
