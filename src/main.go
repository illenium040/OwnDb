package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"os/signal"
	"own-db/src/internal/config"
	"syscall"
	"time"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Read()
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	ctx := context.Background()

	con, err := pgx.Connect(ctx, cfg.DbUrl())
	if err != nil {
		return fmt.Errorf("db connection: %w", err)
	}

	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err := con.Close(ctxWithTimeout)
		if err != nil {
			log.Printf("db close: %v", err)
		}
		close(stopped)
	}()

	err = con.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("db ping: %w", err)
	}

	<-stopped

	return nil
}
