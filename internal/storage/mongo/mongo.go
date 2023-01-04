package mongo

import (
	"context"
	"crud/internal/config"
	"github.com/hendratommy/mongo-sequence/pkg/sequence"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewClient(cfg *config.Config, lgr zerolog.Logger) (*mongo.Client, *mongo.Collection) {
	lgr = lgr.With().Str("db", "mongo").Logger()

	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.Mongo.URI))
	if err != nil {
		lgr.Fatal().Err(err).Msg("failed to create mongo client")
	}
	if err = client.Connect(context.TODO()); err != nil {
		lgr.Fatal().Err(err).Msg("failed to connect mongo client")
	}
	if err = client.Ping(context.TODO(), nil); err != nil {
		lgr.Fatal().Err(err).Msg("failed to ping mongo client")
	}

	seqColl := client.Database(cfg.Mongo.DB).Collection("sequences")
	sequence.SetupDefaultSequence(seqColl, 1*time.Second)

	lgr.Debug().Msg("connection established")

	return client, seqColl
}

type Model struct {
	cfg     *config.Config
	lgr     zerolog.Logger
	client  *mongo.Client
	seqColl *mongo.Collection
}
