package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/refinery/types"
)

var createProjectId = createIdGenerator(8)

type Project struct {
	Id string `db:"id"`

	Name string `db:"name"`

	CloneUrl string `db:"clone_url"`

	Installable string `db:"installable"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func ProjectQuery() *goqu.SelectDataset {
	query := dialect.From("projects").
		Select(
			"projects.id",

			"projects.name",

			"projects.clone_url",

			"projects.installable",

			"projects.created",
			"projects.updated",
		)

	return query
}

func (db DB) GetProjectById(ctx context.Context, projectId string) (Project, error) {
	query := ProjectQuery().
		Where(goqu.I("projects.id").Eq(projectId))

	return ember.Single[Project](db.db, ctx, query)
}

func (db DB) GetAllProjectIds(ctx context.Context) ([]string, error) {
	query := dialect.From("projects").
		Select("projects.id")

	return ember.Multiple[string](db.db, ctx, query)
}

type GetProjectsParams struct {
	Page types.PageParams
}

func (db DB) GetProjects(
	ctx context.Context,
	params GetProjectsParams,
) ([]Project, types.Page, error) {
	query := ProjectQuery()

	var err error

	page, err := buildPage(ctx, db.db, params.Page, query, "projects.id")
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := ember.Multiple[Project](db.db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

type CreateProjectParams struct {
	Id string

	Name string

	CloneUrl string

	Installable string

	Created int64
	Updated int64
}

func (db DB) CreateProject(
	ctx context.Context,
	params CreateProjectParams,
) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = createProjectId()
	}

	query := dialect.Insert("projects").
		Rows(goqu.Record{
			"id": params.Id,

			"name": params.Name,

			"clone_url": params.CloneUrl,

			"installable": params.Installable,

			"created": params.Created,
			"updated": params.Updated,
		})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return "", err
	}

	return params.Id, nil
}

type ProjectChanges struct {
	Name types.Change[string]

	CloneUrl types.Change[string]

	Installable types.Change[string]

	Created types.Change[int64]
}

func (db DB) UpdateProject(
	ctx context.Context,
	projectId string,
	changes ProjectChanges,
) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "clone_url", changes.CloneUrl)

	addToRecord(record, "installable", changes.Installable)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("projects").
		Set(record).
		Where(goqu.I("projects.id").Eq(projectId))

	_, err := db.db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteProject(ctx context.Context, projectId string) error {
	query := dialect.Delete("projects").
		Where(goqu.I("projects.id").Eq(projectId))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
