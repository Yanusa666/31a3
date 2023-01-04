package http_server

import (
	"context"
	"crud/internal/config"
	"crud/internal/http_server/handlers"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"net"
	"net/http"
	"time"
)

type Server struct {
	cfg        *config.Config
	lgr        zerolog.Logger
	httpServer *http.Server
}

func NewServer(cfg *config.Config, lgr zerolog.Logger, handler *handlers.Handler,
) (*Server, chan error) {
	netListener, err := net.Listen("tcp", cfg.HttpServer.ListenAddress)
	if err != nil {
		lgr.Fatal().Err(err).Msgf("failed to start listener for http server on %s", cfg.HttpServer.ListenAddress)
	}
	lgr.Debug().Msgf("start listener for http server success on %s", cfg.HttpServer.ListenAddress)

	server := &Server{
		cfg: cfg,
		lgr: lgr,
		httpServer: &http.Server{
			ReadTimeout:       30 * time.Second,
			ReadHeaderTimeout: 30 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       30 * time.Second,
		},
	}

	router := httprouter.New()
	router.RedirectFixedPath = true
	router.RedirectTrailingSlash = true

	router.POST("/authors", handler.Middlware(handler.AddAuthor))
	router.GET("/authors", handler.Middlware(handler.ListAuthors))
	router.PUT("/authors/:id", handler.Middlware(handler.UpdateAuthor))
	router.DELETE("/authors/:id", handler.Middlware(handler.DeleteAuthor))

	router.POST("/posts", handler.Middlware(handler.AddPost))
	router.GET("/posts", handler.Middlware(handler.ListPosts))
	router.PUT("/posts/:id", handler.Middlware(handler.UpdatePost))
	router.DELETE("/posts/:id", handler.Middlware(handler.DeletePost))

	server.httpServer.Handler = router

	listenErrCh := make(chan error, 1)
	go func() {
		listenErrCh <- server.httpServer.Serve(netListener)
	}()

	return server, listenErrCh
}

func (srv *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		cancel()
	}()

	err := srv.httpServer.Shutdown(ctx)
	if err != nil {
		srv.lgr.Error().Err(err).Msg("http server grace shutdown finished with error")
		return err
	}

	srv.lgr.Debug().Msg("http server grace shutdown success")
	time.Sleep(10 * time.Millisecond)
	return nil
}
