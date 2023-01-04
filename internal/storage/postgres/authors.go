package postgres

import (
	"context"
	"crud/internal/config"
	"crud/internal/constants"
	"crud/internal/entities"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"time"
)

type Authors struct {
	Model
}

func NewAuthors(cfg *config.Config, lgr zerolog.Logger, conn *pgxpool.Pool) *Authors {
	lgr = lgr.With().Str("model", "authors").Logger()

	return &Authors{
		Model: Model{
			cfg:  cfg,
			lgr:  lgr,
			conn: conn,
		},
	}
}

func (a *Authors) Add(ctx context.Context, author *entities.Author) (err error) {
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)
	lgr := a.lgr.With().
		Str("api", "Add").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Str("name", author.Name),
		).Logger()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err = a.Model.conn.Exec(ctx,
		`INSERT INTO public.authors(name) 
			 VALUES ($1)`, author.Name)
	if err != nil {
		lgr.Error().Err(err).Msg("db query failed")
		return err
	}

	lgr.Debug().Msg("executed")

	return nil
}

func (a *Authors) List(ctx context.Context) ([]entities.Author, error) {
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)
	lgr := a.lgr.With().
		Str("api", "List").
		Str(constants.RequestIdKey, requestId).
		Logger()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := a.Model.conn.Query(ctx,
		`SELECT id, name
			 FROM public.authors`)
	if err != nil {
		lgr.Error().Err(err).Msg("db query failed")
		return nil, err
	}

	authors := make([]entities.Author, 0, 10)
	for rows.Next() {
		author := entities.Author{}
		err = rows.Scan(&(author.Id), &(author.Name))
		if err != nil {
			lgr.Error().Err(err).Msg("db scan failed")
			return nil, err
		}
		authors = append(authors, author)
	}

	lgr.Debug().Msg("executed")

	return authors, nil
}

func (a *Authors) Update(ctx context.Context, author *entities.Author) (err error) {
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)
	lgr := a.lgr.With().
		Str("api", "Update").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Uint64("id", author.Id).
			Str("name", author.Name),
		).Logger()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err = a.Model.conn.Exec(ctx,
		`UPDATE public.authors
			 SET name = $2
			 WHERE id = $1`, author.Id, author.Name)
	if err != nil {
		lgr.Error().Err(err).Msg("db query failed")
		return err
	}

	lgr.Debug().Msg("executed")

	return nil
}

func (a *Authors) Delete(ctx context.Context, id uint64) (err error) {
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)
	lgr := a.lgr.With().
		Str("api", "Delete").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Uint64("id", id),
		).Logger()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err = a.Model.conn.Exec(ctx,
		`DELETE FROM public.authors
			 WHERE id = $1`, id)
	if err != nil {
		lgr.Error().Err(err).Msg("db query failed")
		return err
	}

	lgr.Debug().Msg("executed")

	return nil
}
