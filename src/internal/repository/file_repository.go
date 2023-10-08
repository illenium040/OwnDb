package repository

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/jackc/pgx/v5"
	"io"
	"own-db/src/internal/domain"
)

type FileRepository struct {
	con *pgx.Conn
}

func NewFileRepository(con *pgx.Conn) FileRepository {
	return FileRepository{con: con}
}

func (r FileRepository) AddFile(ctx context.Context, file domain.FileMeta, fileReader io.Reader) (id uint, err error) {
	tx, err := r.con.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = fmt.Errorf("rollback error: %v, original error: %w", rollbackErr, err)
			}
		}
	}()

	loStorage := tx.LargeObjects()

	loId, err := loStorage.Create(ctx, 0)
	if err != nil {
		return 0, fmt.Errorf("creating large object: %w", err)
	}

	lo, err := loStorage.Open(ctx, loId, pgx.LargeObjectModeWrite)
	if err != nil {
		return 0, fmt.Errorf("opening large object: %w", err)
	}

	hash := sha256.New()
	teeReader := io.TeeReader(fileReader, hash)

	_, err = io.Copy(lo, teeReader)
	if err != nil {
		return 0, fmt.Errorf("copying data to db large object from file: %w", err)
	}

	var dataId uint
	err = r.con.QueryRow(
		ctx,
		`
			insert into main.file_data (hash, data_oid)
			values (@hash, @oid)
			returning id
		`,
		pgx.NamedArgs{
			"hash": base64.URLEncoding.EncodeToString(hash.Sum(nil)),
			"oid":  loId,
		},
	).Scan(&dataId)
	if err != nil {
		return 0, fmt.Errorf("inserting file data row: %w", err)
	}

	fm := fileMetaFromDomain(file)
	err = r.con.QueryRow(
		ctx,
		`
			insert into main.file_meta (file_data_id, name, extension, original_path, size, dt_created, dt_changed)
			values (@dataId, @name, @extension, @originalPath, @size, @dtCreated, @dtChanged)
			returning id
		`,
		pgx.NamedArgs{
			"dataId":       dataId,
			"name":         fm.Name,
			"extension":    fm.Extension,
			"originalPath": fm.OriginalPath,
			"size":         fm.Size,
			"dtCreated":    fm.CreatedAt,
			"dtChanged":    fm.ChangedAt,
		},
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("inserting file data row: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("commit tx: %w", err)
	}

	return id, nil
}

func (r FileRepository) ReadFile(ctx context.Context, id uint, readFn func(meta domain.FileMeta, loReader io.Reader) error) error {
	tx, err := r.con.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = fmt.Errorf("rollback error: %v, original error: %w", rollbackErr, err)
			}
		}
	}()

	loStorage := tx.LargeObjects()

	var loId uint32
	var fm fileMeta
	err = r.con.QueryRow(
		ctx,
		`
			select 
			    fd.data_oid,
			    fm.id,
			    fm.file_data_id,
			    fm.name,
			    fm.extension,
			    fm.original_path,
			    fm.size,
			    fm.dt_created,
			    fm.dt_changed
			from main.file_meta fm 
			    inner join main.file_data fd 
			        on fd.id = fm.file_data_id
			where fd.id = @id
		`,
		pgx.NamedArgs{
			"id": id,
		},
	).Scan(
		&loId,
		&fm.Id,
		&fm.DataId,
		&fm.Name,
		&fm.Extension,
		&fm.OriginalPath,
		&fm.Size,
		&fm.CreatedAt,
		&fm.ChangedAt,
	)
	if err != nil {
		return fmt.Errorf("getting loId by file id: %w", err)
	}

	lo, err := loStorage.Open(ctx, loId, pgx.LargeObjectModeRead)
	if err != nil {
		return fmt.Errorf("opening large object with id=%d: %w", loId, err)
	}

	err = readFn(fileMetaToDomain(fm), lo)
	if err != nil {
		return fmt.Errorf("reading large object: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
