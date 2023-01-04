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

type Authors struct {
	Model
	coll *mongo.Collection
}

func NewAuthors(cfg *config.Config, lgr zerolog.Logger, client *mongo.Client, seqColl *mongo.Collection) *Authors {
	return &Authors{
		Model: Model{
			cfg:     cfg,
			lgr:     lgr,
			client:  client,
			seqColl: seqColl,
		},
		coll: client.Database(cfg.Mongo.DB).Collection("authors"),
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

	seq, err := sequence.NextVal("authors_seq")
	if err != nil {
		lgr.Error().Err(err).Msg("db sequence failed")
		return err
	}
	author.Id = uint64(seq)

	_, err = a.coll.InsertOne(ctx, author)
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

	cursor, err := a.coll.Find(ctx, bson.M{})
	if err != nil {
		lgr.Error().Err(err).Msg("db query failed")
		return nil, err
	}

	authors := make([]entities.Author, 0, 10)
	for cursor.Next(context.TODO()) {
		author := entities.Author{}
		err = cursor.Decode(&author)
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

	_, err = a.coll.UpdateOne(ctx,
		bson.M{"id": author.Id},
		bson.M{"$set": bson.M{"name": author.Name}},
	)
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

	_, err = a.coll.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		lgr.Error().Err(err).Msg("db query failed")
		return err
	}

	lgr.Debug().Msg("executed")

	return nil
}
