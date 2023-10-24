package transport

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"own-db/src/internal/domain"
	"strconv"
)

type FileService interface {
	AddFile(ctx context.Context, folderId domain.FolderId, selectedPath string) (fileId uint, err error)
	DownloadFile(ctx context.Context, fileId uint, selectedPath string) (err error)
	DeleteFile(ctx context.Context, fileId uint) (err error)
	GetFileList(ctx context.Context, folderId domain.FolderId) (fileList []domain.FileMeta, err error)
}

type FileServer struct {
	service FileService
}

func NewFileServer(service FileService) FileServer {
	return FileServer{service: service}
}

func (s FileServer) UploadFile(ctx *gin.Context) {
	type path struct {
		SelectedPath string `json:"selected_path"`
	}

	folderId, err := strconv.Atoi(ctx.Param("folderId"))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("folderId is not a number: %w", err))
		return
	}

	var data path
	err = ctx.Bind(&data)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("reading body: %w", err))
		return
	}

	id, err := s.service.AddFile(ctx.Request.Context(), domain.NewFolderId(folderId), data.SelectedPath)
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

func (s FileServer) DeleteFile(ctx *gin.Context) {
	fileId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fileId is not a number: %w", err))
		return
	}

	err = s.service.DeleteFile(ctx.Request.Context(), uint(fileId))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("delete file: %w", err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (s FileServer) GetFileList(ctx *gin.Context) {
	folderId, err := strconv.Atoi(ctx.Param("folderId"))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fileId is not a number: %w", err))
		return
	}

	fileList, err := s.service.GetFileList(ctx, domain.NewFolderId(folderId))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("get file list of folder %d: %w", folderId, err))
		return
	}

	var list []fileMeta
	for _, fm := range fileList {
		list = append(list, fileMetaFromDomain(fm))
	}

	ctx.JSON(200, list)
}
