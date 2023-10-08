package main

import (
	"OwnDb/src/config"
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	con, err := pgx.Connect(pgx.ConnConfig{
		User:     cfg.DbUser(),
		Password: cfg.DbPassword(),
		Host:     cfg.DbHost(),
		Database: cfg.DbName(),
		Port:     cfg.DbPort(),
	})
	if err != nil {
		return fmt.Errorf("db connection: %w", err)
	}

	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		err := con.Close()
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
