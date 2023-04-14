package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type Server struct {
	log     *logrus.Entry
	service Service
	address string
	version string
	server  *http.Server
}

func New(log *logrus.Logger, service Service, address string, version string) *Server {
	s := Server{
		log:     log.WithField("module", "server"),
		service: service,
		address: address,
		version: version,
	}
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", s.rootHandler)

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", s.usersHandler)
				r.Post("/", s.createUserHandler)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", s.userHandler)
					r.Patch("/", s.updateUserHandler)
					r.Delete("/", s.deleteUserHandler)
				})
			})
		})
	})

	s.server = &http.Server{
		Addr:              s.address,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return &s
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		gfCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		//nolint:contextcheck
		if err := s.server.Shutdown(gfCtx); err != nil {
			s.log.Warnf("err shutting down properly")
		}
	}()
	s.log.Infof("starting server on %s", s.address)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("ListenAndServe faild: %w", err)
	}
	return nil
}
