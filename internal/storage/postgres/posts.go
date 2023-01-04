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

type Posts struct {
	Model
}

func NewPosts(cfg *config.Config, lgr zerolog.Logger, conn *pgxpool.Pool) *Posts {
	lgr = lgr.With().Str("model", "posts").Logger()

	return &Posts{
		Model: Model{
			cfg:  cfg,
			lgr:  lgr,
			conn: conn,
		},
	}
}

func (p *Posts) Add(ctx context.Context, post *entities.Post) (err error) {
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)
	lgr := p.lgr.With().
		Str("api", "Add").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Uint64("author_id", post.AuthorId).
			Str("title", post.Title).
			Str("content", post.Content).
			Time("created_at", post.CreatedAt),
		).Logger()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err = p.Model.conn.Exec(ctx,
		`INSERT INTO public.posts(author_id, title, content, created_at) 
			 VALUES ($1, $2, $3, $4)`, post.AuthorId, post.Title, post.Content, post.CreatedAt)
	if err != nil {
		lgr.Error().Err(err).Msg("db query failed")
		return err
	}

	lgr.Debug().Msg("executed")

	return nil
}

func (p *Posts) List(ctx context.Context) ([]entities.Post, error) {
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)
	lgr := p.lgr.With().
		Str("api", "List").
		Str(constants.RequestIdKey, requestId).
		Logger()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := p.Model.conn.Query(ctx,
		`SELECT id, author_id, title, content, created_at
			 FROM public.posts`)
	if err != nil {
		lgr.Error().Err(err).Msg("db query failed")
		return nil, err
	}

	posts := make([]entities.Post, 0, 10)
	for rows.Next() {
		post := entities.Post{}
		err = rows.Scan(&(post.Id), &(post.AuthorId), &(post.Title), &(post.Content), &(post.CreatedAt))
		if err != nil {
			lgr.Error().Err(err).Msg("db scan failed")
			return nil, err
		}
		posts = append(posts, post)
	}

	lgr.Debug().Msg("executed")

	return posts, nil
}

func (p *Posts) Update(ctx context.Context, post *entities.Post) (err error) {
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)
	lgr := p.lgr.With().
		Str("api", "Update").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Uint64("id", post.Id).
			Uint64("author_id", post.AuthorId).
			Str("title", post.Title).
			Str("content", post.Content).
			Time("created_at", post.CreatedAt),
		).Logger()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err = p.Model.conn.Exec(ctx,
		`UPDATE public.posts
			 SET author_id = $2, title = $3, content = $4, created_at = $5
			 WHERE id = $1`, post.Id, post.AuthorId, post.Title, post.Content, post.CreatedAt)
	if err != nil {
		lgr.Error().Err(err).Msg("db query failed")
		return err
	}

	lgr.Debug().Msg("executed")

	return nil
}

func (p *Posts) Delete(ctx context.Context, id uint64) (err error) {
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)
	lgr := p.lgr.With().
		Str("api", "Delete").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Uint64("id", id),
		).Logger()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err = p.Model.conn.Exec(ctx,
		`DELETE public.authors
			 WHERE id = $1`, id)
	if err != nil {
		lgr.Error().Err(err).Msg("db query failed")
		return err
	}

	lgr.Debug().Msg("executed")

	return nil
}
