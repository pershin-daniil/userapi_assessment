package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := store.New(log, storePath)
	app := service.New(log, s)
	http := server.New(log, app, port, version)

	go func() {
		signCh := make(chan os.Signal, 1)
		signal.Notify(signCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
		<-signCh
		log.Infof("Received signal, shutting down...")
		cancel()
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := http.Run(ctx); err != nil {
			log.Panic(err)
		}
	}()
	wg.Wait()
}
