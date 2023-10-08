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

	//repo := repository.NewFileRepository(con)

	//path := "C:\\Users\\serj-\\Рабочий стол\\Main\\Аниме\\_--5SsruWa0.jpg"
	//file, err := os.OpenFile(path, os.O_RDONLY, 0400)
	//if err != nil {
	//	return fmt.Errorf("open file: %w", err)
	//}

	//stat, err := file.Stat()
	//if err != nil {
	//	return fmt.Errorf("get file stat: %w", err)
	//}

	//fileMeta, _ := domain.NewFileMeta(0, 0, stat.Name(), filepath.Ext(path), path, uint32(stat.Size()), stat.ModTime(), nil)

	//id, err := repo.AddFile(ctx, fileMeta, file)
	//if err != nil {
	//	return fmt.Errorf("add file: %w", err)
	//}
	//
	//fmt.Printf("id: %d", id)

	//err = repo.ReadFile(ctx, 2, func(meta domain.FileMeta, loReader io.Reader) (err error) {
	//	file, err := os.Create(filepath.Join("./bin", meta.Name()))
	//	if err != nil {
	//		return fmt.Errorf("creating file: %w", err)
	//	}
	//
	//	defer func() {
	//		closeErr := file.Close()
	//		if closeErr != nil {
	//			err = fmt.Errorf("file closing: %v, original err: %w", closeErr, err)
	//		}
	//	}()
	//
	//	_, err = io.Copy(file, loReader)
	//	if err != nil {
	//		return fmt.Errorf("copy large object data to file: %w", err)
	//	}
	//
	//	return nil
	//})
	//if err != nil {
	//	return err
	//}

	<-stopped

	return nil
}
