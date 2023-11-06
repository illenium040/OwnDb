package services

import (
	"context"
	"fmt"
	"own-db/src/internal/domain"
)

type FolderRepository interface {
	Create(ctx context.Context, parentFolderId domain.FolderId, name string) (domain.Folder, error)
	Rename(ctx context.Context, id domain.FolderId, name string) (domain.Folder, error)
	Move(ctx context.Context, id domain.FolderId, destFolderId domain.FolderId) (domain.Folder, error)
	Delete(ctx context.Context, id domain.FolderId) error
}

type FolderService struct {
	repo FolderRepository
}

func NewFolderService(repo FolderRepository) FolderService {
	return FolderService{repo: repo}
}

func (s FolderService) Create(ctx context.Context, parentFolderId domain.FolderId, name string) (domain.Folder, error) {
	return s.repo.Create(ctx, parentFolderId, name)
}

func (s FolderService) Rename(ctx context.Context, id domain.FolderId, name string) (domain.Folder, error) {
	if id.IsRoot() {
		return domain.Folder{}, fmt.Errorf("cann't rename a root folder")
	}

	return s.repo.Rename(ctx, id, name)
}

func (s FolderService) Move(ctx context.Context, id domain.FolderId, destFolderId domain.FolderId) (domain.Folder, error) {
	if id.IsRoot() {
		return domain.Folder{}, fmt.Errorf("cann't move a root folder")
	}

	return s.repo.Move(ctx, id, destFolderId)
}

func (s FolderService) Delete(ctx context.Context, id domain.FolderId) error {
	if id.IsRoot() {
		return fmt.Errorf("cann't remove a root folder")
	}

	return s.repo.Delete(ctx, id)
}
