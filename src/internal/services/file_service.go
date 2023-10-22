package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"own-db/src/internal/domain"
	"own-db/src/internal/dto"
	"path/filepath"
)

type Repository interface {
	AddFile(ctx context.Context, file dto.FileMeta, fileReader io.Reader) (id uint, err error)
	ReadFile(ctx context.Context, id uint, readFn func(meta domain.FileMeta, loReader io.Reader) error) error
	DeleteFile(ctx context.Context, id uint) error
	GetFileList(ctx context.Context, folderId domain.FolderId) (fileList []domain.FileMeta, err error)
}

type FileService struct {
	repo Repository
}

func NewFileService(repo Repository) FileService {
	return FileService{repo: repo}
}

func (s FileService) AddFile(ctx context.Context, folderId domain.FolderId, selectedPath string) (fileId uint, err error) {
	file, err := os.OpenFile(selectedPath, os.O_RDONLY, 0400)
	if err != nil {
		return 0, fmt.Errorf("open file: %w", err)
	}

	stat, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("get file stat: %w", err)
	}

	fileMeta := dto.FileMeta{
		FolderId:     folderId,
		Name:         stat.Name(),
		Extension:    filepath.Ext(selectedPath),
		OriginalPath: selectedPath,
		Size:         uint32(stat.Size()),
		CreatedAt:    stat.ModTime(),
		ChangedAt:    nil,
	}

	fileId, err = s.repo.AddFile(ctx, fileMeta, file)
	if err != nil {
		return 0, fmt.Errorf("add file: %w", err)
	}

	return fileId, nil
}

func (s FileService) DownloadFile(ctx context.Context, fileId uint, selectedPath string) (err error) {
	err = s.repo.ReadFile(ctx, fileId, func(meta domain.FileMeta, loReader io.Reader) (err error) {
		file, err := os.Create(filepath.Join(selectedPath, meta.Name()))
		if err != nil {
			return fmt.Errorf("creating file: %w", err)
		}

		defer func() {
			closeErr := file.Close()
			if closeErr != nil {
				err = fmt.Errorf("file closing: %v, original err: %w", closeErr, err)
			}
		}()

		_, err = io.Copy(file, loReader)
		if err != nil {
			return fmt.Errorf("copy large object data to file: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	return nil
}

func (s FileService) DeleteFile(ctx context.Context, fileId uint) (err error) {
	err = s.repo.DeleteFile(ctx, fileId)
	if err != nil {
		return fmt.Errorf("deleting file: %w", err)
	}

	return nil
}
