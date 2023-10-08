package domain

import (
	"time"
)

type FileMeta struct {
	id           uint
	dataId       uint
	extension    string
	originalPath string
	size         uint32
	createdAt    time.Time
	changedAt    *time.Time
}

func NewFileMeta(
	id uint,
	dataId uint,
	extension string,
	originalPath string,
	size uint32,
	createdAt time.Time,
	changedAt *time.Time,
) (FileMeta, error) {
	return FileMeta{
		id:           id,
		dataId:       dataId,
		extension:    extension,
		originalPath: originalPath,
		size:         size,
		createdAt:    createdAt,
		changedAt:    changedAt,
	}, nil
}

func (f FileMeta) Id() uint {
	return f.id
}

func (f FileMeta) DataId() uint {
	return f.dataId
}

func (f FileMeta) Extension() string {
	return f.extension
}

func (f FileMeta) OriginalPath() string {
	return f.originalPath
}

func (f FileMeta) Size() uint32 {
	return f.size
}

func (f FileMeta) CreatedAt() time.Time {
	return f.createdAt
}

func (f FileMeta) ChangedAt() *time.Time {
	return f.changedAt
}
