package postgres

import (
	"context"
	"crud/internal/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

func NewConn(cfg *config.Config, lgr zerolog.Logger) *pgxpool.Pool {
	lgr = lgr.With().Str("db", "postgres").Logger()

	pgConf, err := pgxpool.ParseConfig(cfg.Postgres.URI)
	if err != nil {
		lgr.Fatal().Err(err).Msg("failed parse PostgreSQL config")
	}

	pgPool, err := pgxpool.ConnectConfig(context.Background(), pgConf)
	if err != nil {
		lgr.Fatal().Err(err).Msg("failed connect to PostgreSQL")
	}

	err = pgPool.Ping(context.Background())
	if err != nil {
		lgr.Fatal().Err(err).Msg("unsuccessful ping attempt")
	}

	lgr.Debug().Msg("connection established")

	return pgPool
}

type Model struct {
	cfg  *config.Config
	lgr  zerolog.Logger
	conn *pgxpool.Pool
}
