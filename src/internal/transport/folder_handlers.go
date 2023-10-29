package transport

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"own-db/src/internal/domain"
)

type FolderService interface {
	Create(ctx context.Context, parentFolderId domain.FolderId, name string) (domain.Folder, error)
	Rename(ctx context.Context, id domain.FolderId, name string) (domain.Folder, error)
	Move(ctx context.Context, id domain.FolderId, destFolderId domain.FolderId) (domain.Folder, error)
	Delete(ctx context.Context, id domain.FolderId) error
}

type FolderHandlers struct {
	folderService FolderService
}

func NewFolderHandlers(folderService FolderService) FolderHandlers {
	return FolderHandlers{folderService: folderService}
}

func (h FolderHandlers) CreateFolder(ctx *gin.Context) {
	type body struct {
		Name           string `json:"name"`
		ParentFolderId int    `json:"parentFolderId"`
	}

	var data body
	err := ctx.Bind(&data)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("reading body: %w", err))
	}

	result, err := h.folderService.Create(ctx, domain.NewFolderId(data.ParentFolderId), data.Name)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("creating fodler: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, folderFromDomain(result))
}

func (h FolderHandlers) RenameFolder(ctx *gin.Context) {
	type body struct {
		FolderId int    `json:"folderId"`
		Name     string `json:"name"`
	}

	var data body
	err := ctx.Bind(&data)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("reading body: %w", err))
		return
	}

	result, err := h.folderService.Rename(ctx, domain.NewFolderId(data.FolderId), data.Name)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("rename folder: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, folderFromDomain(result))
}

func (h FolderHandlers) MoveFolder(ctx *gin.Context) {
	type body struct {
		FolderId     int `json:"folderId"`
		DestFolderId int `json:"destFolderId"`
	}

	var data body
	err := ctx.Bind(&data)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("reading body: %w", err))
		return
	}

	result, err := h.folderService.Move(ctx, domain.NewFolderId(data.FolderId), domain.NewFolderId(data.DestFolderId))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("move folder: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, folderFromDomain(result))
}

func (h FolderHandlers) DeleteFolder(ctx *gin.Context) {
	type body struct {
		FolderId int `json:"folderId"`
	}

	var data body
	err := ctx.Bind(&data)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("reading body: %w", err))
		return
	}

	err = h.folderService.Delete(ctx, domain.NewFolderId(data.FolderId))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("rename folder: %w", err))
		return
	}

	ctx.Status(http.StatusOK)
}
