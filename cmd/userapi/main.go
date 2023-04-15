package main

import (
	"context"
	"os/signal"
	"syscall"

	"refactoring/internal/logger"
	"refactoring/internal/server"
	"refactoring/pkg/service"
	"refactoring/pkg/store"
)

const (
	port    = ":3333"
	version = "0.0.1"
)

const storePath = `users.json`

func main() {
	log := logger.New()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	defer stop()

	s := store.New(log, storePath)
	app := service.New(log, s)
	http := server.New(log, app, port, version)

	if err := http.Run(ctx); err != nil {
		log.Panic(err)
	}
}
