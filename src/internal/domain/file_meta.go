package domain

import (
	"time"
)

type FileMeta struct {
	id           uint
	dataId       uint
	folderId     *int
	name         string
	extension    string
	originalPath string
	size         uint32
	createdAt    time.Time
	changedAt    *time.Time
}

func NewFileMeta(
	id uint,
	dataId uint,
	folderId *int,
	name string,
	extension string,
	originalPath string,
	size uint32,
	createdAt time.Time,
	changedAt *time.Time,
) FileMeta {
	return FileMeta{
		id:           id,
		dataId:       dataId,
		folderId:     folderId,
		name:         name,
		extension:    extension,
		originalPath: originalPath,
		size:         size,
		createdAt:    createdAt,
		changedAt:    changedAt,
	}
}

func (f FileMeta) Id() uint {
	return f.id
}

func (f FileMeta) DataId() uint {
	return f.dataId
}

func (f FileMeta) FolderId() *int {
	return f.folderId
}

func (f FileMeta) Name() string {
	return f.name
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
