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

type FileRepository interface {
	Add(ctx context.Context, file dto.FileMeta, fileReader io.Reader) (fm domain.FileMeta, err error)
	Read(ctx context.Context, id int, readFn func(meta domain.FileMeta, loReader io.Reader) error) error
	Delete(ctx context.Context, id int) error
	Move(ctx context.Context, id int, destFolderId domain.FolderId) (domain.FileMeta, error)
	Rename(ctx context.Context, id int, name string) (domain.FileMeta, error)
}

type FileService struct {
	repo FileRepository
}

func NewFileService(repo FileRepository) FileService {
	return FileService{repo: repo}
}

func (s FileService) Create(ctx context.Context, folderId domain.FolderId, selectedPath string) (fm domain.FileMeta, err error) {
	file, err := os.OpenFile(selectedPath, os.O_RDONLY, 0400)
	if err != nil {
		return domain.FileMeta{}, fmt.Errorf("open file: %w", err)
	}

	stat, err := file.Stat()
	if err != nil {
		return domain.FileMeta{}, fmt.Errorf("get file stat: %w", err)
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

	fm, err = s.repo.Add(ctx, fileMeta, file)
	if err != nil {
		return domain.FileMeta{}, fmt.Errorf("add file: %w", err)
	}

	return fm, nil
}

func (s FileService) Download(ctx context.Context, fileId int, selectedPath string) (err error) {
	err = s.repo.Read(ctx, fileId, func(meta domain.FileMeta, loReader io.Reader) (err error) {
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

func (s FileService) Delete(ctx context.Context, fileId int) (err error) {
	return s.repo.Delete(ctx, fileId)
}

func (s FileService) Move(ctx context.Context, id int, destFolderId domain.FolderId) (fileMeta domain.FileMeta, err error) {
	return s.repo.Move(ctx, id, destFolderId)
}

func (s FileService) Rename(ctx context.Context, id int, name string) (fileMeta domain.FileMeta, err error) {
	return s.repo.Rename(ctx, id, name)
}
