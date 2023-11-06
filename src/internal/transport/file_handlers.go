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
	Create(ctx context.Context, folderId domain.FolderId, selectedPath string) (fm domain.FileMeta, err error)
	Download(ctx context.Context, id int, selectedPath string) (err error)
	Move(ctx context.Context, id int, destFolderId domain.FolderId) (domain.FileMeta, error)
	Rename(ctx context.Context, id int, name string) (domain.FileMeta, error)
	Delete(ctx context.Context, id int) (err error)
}

type FileHandlers struct {
	fileService FileService
}

func NewFileHandlers(service FileService) FileHandlers {
	return FileHandlers{fileService: service}
}

func (h FileHandlers) UploadFile(ctx *gin.Context) {
	type path struct {
		SelectedPath string `json:"selected_path"`
	}

	folderId, err := strconv.Atoi(ctx.Param("folderId"))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("folderId is not a number: %w", err))
		return
	}

	var data path
	err = ctx.Bind(&data)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("reading body: %w", err))
		return
	}

	fm, err := h.fileService.Create(ctx.Request.Context(), domain.NewFolderId(folderId), data.SelectedPath)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("adding file: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, fileMetaFromDomain(fm))
}

func (h FileHandlers) DownloadFile(ctx *gin.Context) {
	selectedPath := ctx.Query("selectedPath")
	fileId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fileId is not a number: %w", err))
		return
	}

	err = h.fileService.Download(ctx.Request.Context(), fileId, selectedPath)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("download file: %w", err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (h FileHandlers) DeleteFile(ctx *gin.Context) {
	fileId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fileId is not a number: %w", err))
		return
	}

	err = h.fileService.Delete(ctx.Request.Context(), fileId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("delete file: %w", err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (h FileHandlers) MoveFile(ctx *gin.Context) {
	fileId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fileId is not a number: %w", err))
		return
	}

	destFolderId, err := strconv.Atoi(ctx.Param("folderId"))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("destFolderId is not a number: %w", err))
		return
	}

	fm, err := h.fileService.Move(ctx, fileId, domain.NewFolderId(destFolderId))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("moving file: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, fileMetaFromDomain(fm))
}

func (h FileHandlers) RenameFile(ctx *gin.Context) {
	type body struct {
		Name string `json:"name"`
	}

	var data body
	err := ctx.Bind(&data)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("reading body: %w", err))
		return
	}

	fileId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fileId is not a number: %w", err))
		return
	}

	fm, err := h.fileService.Rename(ctx, fileId, data.Name)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("renaming file: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, fileMetaFromDomain(fm))
}
