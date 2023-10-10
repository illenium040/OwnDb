package transport

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FileService interface {
	AddFile(ctx context.Context, selectedPath string) (fileId uint, err error)
	DownloadFile(ctx context.Context, fileId uint, selectedPath string) (err error)
}

type FileServer struct {
	service FileService
}

func NewFileServer(service FileService) FileServer {
	return FileServer{service: service}
}

func (s FileServer) AddFile(ctx *gin.Context) {
	type path struct {
		SelectedPath string `json:"selected_path"`
	}

	var data path
	err := ctx.Bind(&data)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("reading body: %w", err))
		return
	}

	id, err := s.service.AddFile(ctx.Request.Context(), data.SelectedPath)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("adding file: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"fileId": id,
	})
}

func (s FileServer) DownloadFile(ctx *gin.Context) {
	selectedPath := ctx.Query("selectedPath")
	fileId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fileId is not a number: %w", err))
		return
	}

	err = s.service.DownloadFile(ctx.Request.Context(), uint(fileId), selectedPath)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("download file: %w", err))
		return
	}

	ctx.Status(http.StatusOK)
}
