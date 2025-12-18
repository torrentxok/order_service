package http

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/torrentxok/order_service/internal/http/handler"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func NewServer(addr string, orderhandler *handler.OrderHandler, logger *zap.Logger) *Server {
	r := chi.NewRouter()

	// базовые middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Use(middleware.Logger)

	r.Route("/order", func(r chi.Router) {
		r.Get("/{order_uid}", orderhandler.GetOrder)
	})

	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return &Server{
		httpServer: httpServer,
		logger:     logger,
	}
}

func (s *Server) Start() {
	s.logger.Info("http server started", zap.String("addr", s.httpServer.Addr))

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("http server error", zap.Error(err))
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down http server")

	return s.httpServer.Shutdown(ctx)
}
