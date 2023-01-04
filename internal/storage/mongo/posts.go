package mongo

import (
	"context"
	"crud/internal/config"
	"crud/internal/constants"
	"crud/internal/entities"
	"github.com/hendratommy/mongo-sequence/pkg/sequence"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Posts struct {
	Model
	coll *mongo.Collection
}

func NewPosts(cfg *config.Config, lgr zerolog.Logger, client *mongo.Client) *Posts {
	return &Posts{
		Model: Model{
			cfg:    cfg,
			lgr:    lgr,
			client: client,
		},
		coll: client.Database(cfg.Mongo.DB).Collection("posts"),
	}
}

func (p *Posts) Add(ctx context.Context, post *entities.Post) error {
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

	seq, err := sequence.NextVal("posts_seq")
	if err != nil {
		lgr.Error().Err(err).Msg("db sequence failed")
		return err
	}
	post.Id = uint64(seq)

	_, err = p.coll.InsertOne(ctx, post)
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

	cursor, err := p.coll.Find(ctx, bson.M{})
	if err != nil {
		lgr.Error().Err(err).Msg("db query failed")
		return nil, err
	}

	posts := make([]entities.Post, 0, 10)
	for cursor.Next(context.TODO()) {
		post := entities.Post{}
		err = cursor.Decode(&post)
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

	_, err = p.coll.UpdateOne(ctx,
		bson.M{"id": post.Id},
		bson.M{"$set": bson.M{
			"author_id": post.AuthorId,
			"title":     post.Title,
			"content":   post.Content},
			"created_at": post.CreatedAt,
		},
	)
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

	_, err = p.coll.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		lgr.Error().Err(err).Msg("db query failed")
		return err
	}

	lgr.Debug().Msg("executed")

	return nil
}
