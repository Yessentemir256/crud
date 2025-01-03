package main

import (
	"context"
	"github.com/Yessentemir256/crud/cmd/server/app"
	"github.com/Yessentemir256/crud/pkg/customers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/dig"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	dsn := "postgres://app:pass@localhost:5435/db"

	if err := execute(host, port, dsn); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func execute(host string, port string, dsn string) (err error) {
	deps := []interface{}{
		app.NewServer,
		mux.NewRouter, // mux "github.com/gorilla/mux"
		func() (*pgxpool.Pool, error) {
			ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
			return pgxpool.Connect(ctx, dsn)
		},
		customers.NewService,
		func(server *app.Server) *http.Server {
			return &http.Server{
				Addr:    net.JoinHostPort(host, port),
				Handler: server,
			}
		},
	}

	container := dig.New()
	for _, dep := range deps {
		err = container.Provide(dep)
		if err != nil {
			return err
		}
	}

	err = container.Invoke(func(server *app.Server) {
		server.Init()
	})
	if err != nil {
		return err
	}

	return container.Invoke(func(server *http.Server) error {
		return server.ListenAndServe()
	})
}
