package storage

import (
	"context"
	"crud/internal/config"
	"crud/internal/entities"
	"crud/internal/storage/mongo"
	"crud/internal/storage/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	_mongo "go.mongodb.org/mongo-driver/mongo"
)

type IAuthors interface {
	Add(context.Context, *entities.Author) error
	List(context.Context) ([]entities.Author, error)
	Update(context.Context, *entities.Author) error
	Delete(context.Context, uint64) error
}

type IPosts interface {
	Add(context.Context, *entities.Post) error
	List(context.Context) ([]entities.Post, error)
	Update(context.Context, *entities.Post) error
	Delete(context.Context, uint64) error
}

type Storage struct {
	Authors  IAuthors
	Posts    IPosts
	pgConn   *pgxpool.Pool
	mgClient *_mongo.Client
}

func NewStorage(cfg *config.Config, lgr zerolog.Logger) *Storage {
	var (
		authors  IAuthors
		posts    IPosts
		pgConn   *pgxpool.Pool
		mgClient *_mongo.Client
	)

	switch cfg.Database.Name {
	case "postgres":
		pgConn = postgres.NewConn(cfg, lgr)
		authors = postgres.NewAuthors(cfg, lgr, pgConn)
		posts = postgres.NewPosts(cfg, lgr, pgConn)
	case "mongo":
		var seqColl *_mongo.Collection
		mgClient, seqColl = mongo.NewClient(cfg, lgr)
		authors = mongo.NewAuthors(cfg, lgr, mgClient, seqColl)
		posts = mongo.NewPosts(cfg, lgr, mgClient)
	default:
		lgr.Fatal().Msg("incorrect database name")
	}

	return &Storage{
		Authors:  authors,
		Posts:    posts,
		pgConn:   pgConn,
		mgClient: mgClient,
	}
}

func (s *Storage) Shutdown() {
	if s.pgConn != nil {
		s.pgConn.Close()
	}

	if s.mgClient != nil {
		s.mgClient.Disconnect(context.TODO())
	}
}
