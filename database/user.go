package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/refinery/tools/utils"
	"github.com/nanoteck137/refinery/types"
	"github.com/nanoteck137/pyrin/ember"
)

type UserSettings struct {
	Id            string         `db:"id"`
	QuickPlaylist sql.NullString `db:"quick_playlist"`
}

type User struct {
	Id    string `db:"id"`
	Email string `db:"email"`

	DisplayName string `db:"display_name"`
	Role        string `db:"role"`

	Picture sql.NullString `db:"picture"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	QuickPlaylist sql.NullString `db:"quick_playlist"`
}

func (u User) ToUserSettings() UserSettings {
	return UserSettings{
		Id:            u.Id,
		QuickPlaylist: u.QuickPlaylist,
	}
}

func UserQuery() *goqu.SelectDataset {
	query := dialect.From("users").
		Select(
			"users.id",
			"users.email",

			"users.display_name",
			"users.role",

			"users.picture",

			"users.created",
			"users.updated",

			"users_settings.quick_playlist",
		).
		LeftJoin(
			goqu.I("users_settings"),
			goqu.On(goqu.I("users.id").Eq(goqu.I("users_settings.id"))),
		)

	return query
}

// TODO(patrik): This needs fixing
func UserSettingsQuery() *goqu.SelectDataset {
	query := dialect.From("users_settings").
		Select(
			"users_settings.id",

			"users_settings.quick_playlist",
		)

	return query
}

func (db DB) GetAllUsers(ctx context.Context) ([]User, error) {
	query := UserQuery()

	return ember.Multiple[User](db.db, ctx, query)
}

func (db DB) GetUsersIn(ctx context.Context, in any) ([]User, error) {
	query := UserQuery().
		Where(goqu.I("users.id").In(in))

	return ember.Multiple[User](db.db, ctx, query)
}

func (db DB) CountUsers(ctx context.Context) (int, error) {
	query := UserQuery().Select(goqu.COUNT("users.id").As("count"))

	return ember.Single[int](db.db, ctx, query)
}

func (db DB) GetUserById(ctx context.Context, id string) (User, error) {
	query := UserQuery().
		Where(goqu.I("users.id").Eq(id))

	return ember.Single[User](db.db, ctx, query)
}

func (db DB) GetUserByUsername(ctx context.Context, username string) (User, error) {
	query := UserQuery().
		Where(goqu.I("users.username").Eq(username))

	return ember.Single[User](db.db, ctx, query)
}

func (db DB) GetUserByEmail(ctx context.Context, email string) (User, error) {
	query := UserQuery().
		Where(goqu.I("users.email").Eq(email))

	return ember.Single[User](db.db, ctx, query)
}

func (db DB) GetUserSettingsById(ctx context.Context, id string) (UserSettings, error) {
	query := UserSettingsQuery().
		Where(goqu.I("users_settings.id").Eq(id))

	return ember.Single[UserSettings](db.db, ctx, query)
}

type CreateUserParams struct {
	Id    string
	Email string

	DisplayName string
	Role        string

	Picture sql.NullString

	Created int64
	Updated int64
}

// TODO(patrik): Change to return id
func (db DB) CreateUser(ctx context.Context, params CreateUserParams) (User, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = utils.CreateUserId()
	}

	query := dialect.
		Insert("users").
		Rows(goqu.Record{
			"id":    params.Id,
			"email": params.Email,

			"display_name": params.DisplayName,
			"role":         params.Role,

			"picture": params.Picture,

			"created": params.Created,
			"updated": params.Updated,
		}).
		// TODO(patrik): Fix this
		Returning(
			"users.id",
			"users.email",

			"users.display_name",
			"users.role",

			"users.picture",

			"users.created",
			"users.updated",
		)

	return ember.Single[User](db.db, ctx, query)
}

type UserChanges struct {
	DisplayName types.Change[string]
	Role        types.Change[string]

	Picture types.Change[sql.NullString]

	Created types.Change[int64]
}

func (db DB) UpdateUser(ctx context.Context, id string, changes UserChanges) error {
	record := goqu.Record{}

	addToRecord(record, "display_name", changes.DisplayName)
	addToRecord(record, "role", changes.Role)

	addToRecord(record, "picture", changes.Picture)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("users").
		Set(record).
		Where(goqu.I("users.id").Eq(id))

	_, err := db.db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) UpdateUserSettings(ctx context.Context, settings UserSettings) error {
	query := dialect.Insert("users_settings").
		Rows(goqu.Record{
			"id":             settings.Id,
			"quick_playlist": settings.QuickPlaylist,
		}).
		OnConflict(goqu.DoUpdate("id", goqu.Record{
			"quick_playlist": settings.QuickPlaylist,
		}))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
