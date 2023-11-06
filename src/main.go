package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"own-db/src/internal/config"
	"own-db/src/internal/repository"
	"own-db/src/internal/services"
	"own-db/src/internal/transport"
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

	fileRepo := repository.NewFileRepository(con)
	folderRepo := repository.NewFolderRepository(con)

	fileService := services.NewFileService(fileRepo)
	folderService := services.NewFolderService(folderRepo)

	fileHandlers := transport.NewFileHandlers(fileService)
	folderHandlers := transport.NewFolderHandlers(folderService)

	router := gin.Default()

	httpServer := &http.Server{
		Addr:           ":8000",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// files
	router.POST("/api/v1/folder/:folderId/file/add", fileHandlers.UploadFile)
	router.GET("/api/v1/file/:id", fileHandlers.DownloadFile)
	router.DELETE("/api/v1/file/:id", fileHandlers.DeleteFile)
	router.PATCH("/api/v1/file/:id/rename", fileHandlers.RenameFile)
	router.PATCH("/api/v1/file/:id/move/:folderId", fileHandlers.MoveFile)

	// folders
	router.POST("/api/v1/folder", folderHandlers.CreateFolder)
	router.DELETE("/api/v1/folder/:id", folderHandlers.DeleteFolder)
	router.PATCH("/api/v1/folder/:id/rename", folderHandlers.RenameFolder)
	router.PATCH("/api/v1/folder/:id/move/:folderId", folderHandlers.MoveFolder)

	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err = httpServer.Shutdown(ctxWithTimeout)
		if err != nil {
			log.Printf("http server close: %v", err)
		}

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

	err = httpServer.ListenAndServe()
	if err != nil {
		log.Printf("db close: %v", err)
	}

	<-stopped

	return nil
}
