package main

import (
	"context"
	"crud/internal/config"
	"crud/internal/http_server"
	"crud/internal/http_server/handlers"
	"crud/internal/storage"
	"crud/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.NewConfig()

	lgr, err := logger.NewLogger(os.Stdout, cfg.LogLevel)
	if err != nil {
		log.Fatalln(err)
	}

	lgr = lgr.With().
		CallerWithSkipFrameCount(2).
		Str("app", "crud").
		Logger()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	stor := storage.NewStorage(cfg, lgr)

	handler := handlers.NewHandler(cfg, lgr, stor)
	httpServer, listenHTTPErr := http_server.NewServer(cfg, lgr, handler)

mainLoop:
	for {
		select {
		case <-ctx.Done():
			break mainLoop

		case err = <-listenHTTPErr:
			if err != nil {
				lgr.Error().Err(err).Msg("http server error")
				shutdownCh <- syscall.SIGTERM
			}

		case sig := <-shutdownCh:
			lgr.Info().Msgf("shutdown signal received: %s", sig.String())
			ctx, cancel = context.WithTimeout(ctx, 10*time.Second)

			if err = httpServer.Shutdown(); err != nil {
				lgr.Error().Err(err).Msg("shutdown http server error")
			}

			stor.Shutdown()

			lgr.Info().Msg("server loop stopped")
			cancel()
		}
	}
}
