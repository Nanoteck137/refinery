package database

import (
	"context"
	"log/slog"
	"os"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/refinery/tools/utils"
	"github.com/nanoteck137/refinery/types"
	"github.com/nrednav/cuid2"
)

func addToRecord[T any](record goqu.Record, name string, change types.Change[T]) {
	if change.Changed {
		record[name] = change.Value
	}
}

func createIdGenerator(length int) func() string {
	res, err := cuid2.Init(cuid2.WithLength(length))
	if err != nil {
		slog.Error("failed to create id generator", "err", err)
		os.Exit(1)
	}

	return res
}

func buildPage(
	ctx context.Context,
	db ember.DB,
	params types.PageParams,
	query *goqu.SelectDataset,
	countCol any,
) (types.Page, error) {
	countQuery := query.Select(goqu.COUNT(countCol))

	totalItems, err := ember.Single[int](db, ctx, countQuery)
	if err != nil {
		return types.Page{}, err
	}

	return types.Page{
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalItems: totalItems,
		TotalPages: utils.TotalPages(params.PerPage, totalItems),
	}, nil
}

func applyPageParams(
	params types.PageParams,
	query *goqu.SelectDataset,
) *goqu.SelectDataset {
	return query.
		Limit(uint(params.PerPage)).
		Offset(uint(params.Page * params.PerPage))
}
