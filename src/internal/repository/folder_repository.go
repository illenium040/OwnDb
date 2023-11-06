package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"own-db/src/internal/domain"
)

type FolderRepository struct {
	con *pgx.Conn
}

func NewFolderRepository(con *pgx.Conn) FolderRepository {
	return FolderRepository{con: con}
}

func (s FolderRepository) Create(ctx context.Context, parentFolderId domain.FolderId, name string) (domain.Folder, error) {
	var f folder
	err := s.con.QueryRow(
		ctx,
		`insert into main.folders(parent_folder_id, name)
		values (@parentFolderId, @name)
		returning id, parent_folder_id, name, dt_created, dt_changed
		`,
		pgx.NamedArgs{
			"parentFolderId": parentFolderId.Value(),
			"name":           name,
		},
	).Scan(
		&f.Id,
		&f.ParentFolderId,
		&f.Name,
		&f.CreatedAt,
		&f.ChangedAt,
	)
	if err != nil {
		return domain.Folder{}, fmt.Errorf("create new folder: %w", err)
	}

	return folderToDomain(f), nil
}

func (s FolderRepository) Rename(ctx context.Context, id domain.FolderId, name string) (domain.Folder, error) {
	var f folder
	err := s.con.QueryRow(
		ctx,
		`update main.folders
		set name = @name
		where id = @id
		returning id, parent_folder_id, name, dt_created, dt_changed
		`,
		pgx.NamedArgs{
			"id":   id.Value(),
			"name": name,
		},
	).Scan(
		&f.Id,
		&f.ParentFolderId,
		&f.Name,
		&f.CreatedAt,
		&f.ChangedAt,
	)
	if err != nil {
		return domain.Folder{}, fmt.Errorf("rename folder: %w", err)
	}

	return folderToDomain(f), nil
}

func (s FolderRepository) Move(ctx context.Context, id domain.FolderId, destFolderId domain.FolderId) (domain.Folder, error) {
	var f folder
	err := s.con.QueryRow(
		ctx,
		`update main.folders
		set parent_folder_id = @destFolderId
		where id = @id
		returning id, parent_folder_id, name, dt_created, dt_changed
		`,
		pgx.NamedArgs{
			"id":           id.Value(),
			"destFolderId": destFolderId.Value(),
		},
	).Scan(
		&f.Id,
		&f.ParentFolderId,
		&f.Name,
		&f.CreatedAt,
		&f.ChangedAt,
	)
	if err != nil {
		return domain.Folder{}, fmt.Errorf("move folder: %w", err)
	}

	return folderToDomain(f), nil
}

func (s FolderRepository) Delete(ctx context.Context, id domain.FolderId) error {
	_, err := s.con.Exec(
		ctx,
		`delete from main.folders
		where id = @id
		`,
		pgx.NamedArgs{
			"id": id.Value(),
		},
	)
	if err != nil {
		return fmt.Errorf("delete folder: %w", err)
	}

	return nil
}
