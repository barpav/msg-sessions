package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"go.elastic.co/ecszerolog"

	"github.com/barpav/msg-sessions/internal/data"
	"github.com/barpav/msg-sessions/internal/rest"
	"github.com/barpav/msg-sessions/internal/users"
)

func main() {
	log.Logger = ecszerolog.New(os.Stdout)

	app := microservice{}
	err := app.launch()

	if err == nil {
		log.Info().Msg("Microservice launched.")
	} else {
		log.Err(err).Msg("Failed to launch microservice")
		app.abort()
	}

	err = app.serveAndShutdownGracefully()

	if err == nil {
		log.Info().Msg("Microservice stopped.")
	} else {
		log.Err(err).Msg("Failed to shutdown microservice gracefully.")
	}
}

type microservice struct {
	api struct {
		public *rest.Service // Specification: https://barpav.github.io/msg-api-spec/#/sessions
		// private *pb.Service   // TODO
	}
	clients struct {
		users *users.Client // github.com/barpav/msg-users
	}
	storage  *data.Storage
	shutdown chan os.Signal
}

func (m *microservice) launch() (err error) {
	m.shutdown = make(chan os.Signal, 2)
	signal.Notify(m.shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)

	m.storage = &data.Storage{}
	err = errors.Join(err, m.storage.Open())

	m.clients.users = &users.Client{}
	err = errors.Join(err, m.clients.users.Connect())

	// m.api.private = &pb.Service{}
	// m.api.private.Start(m.storage)

	m.api.public = &rest.Service{}
	m.api.public.Start(m.storage, m.clients.users)

	return err
}

func (m *microservice) abort() {
	m.shutdown <- syscall.SIGINT
}

func (m *microservice) serveAndShutdownGracefully() (err error) {
	select {
	case <-m.shutdown:
	case <-m.api.public.Shutdown:
		// case <-m.api.private.Shutdown:
	}

	log.Info().Msg("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = errors.Join(err, m.api.public.Stop(ctx))
	// err = errors.Join(err, m.api.private.Stop(ctx))
	err = errors.Join(err, m.clients.users.Disconnect(ctx))
	err = errors.Join(err, m.storage.Close(ctx))

	return err
}
