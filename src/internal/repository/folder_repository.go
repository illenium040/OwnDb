package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"own-db/src/internal/domain"
)

type FolderRepository struct {
	con *pgx.Conn
}

func (r FileRepository) AddFolder(ctx context.Context, folderId domain.FolderId, name string) error {

	return nil
}
