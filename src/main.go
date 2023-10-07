package main

import (
	"OwnDb/src/config"
	"context"
	"github.com/jackc/pgx"
	"log"
	"os"
	"os/signal"
	"strconv"
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
		log.Fatalf("reading config: %v", err)
	}

	port, err := strconv.Atoi(cfg.DbPort())
	if err != nil {
		log.Fatalf("port in not a number: port=%s", cfg.DbPort())
	}

	con, err := pgx.Connect(pgx.ConnConfig{
		User:     cfg.DbUser(),
		Password: cfg.DbPassword(),
		Host:     cfg.DbHost(),
		Database: cfg.DbName(),
		Port:     uint16(port),
	})
	if err != nil {
		log.Fatalf("db connection: %v", err)
	}

	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		err := con.Close()
		if err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(stopped)
	}()

	err = con.Ping(context.Background())
	if err != nil {
		log.Fatalf("db ping: %v", err)
	}

	<-stopped

	return nil
}
