package dto

import (
	"own-db/src/internal/domain"
	"time"
)

type FileMeta struct {
	FolderId     domain.FolderId
	Name         string
	Extension    string
	OriginalPath string
	Size         uint32
	CreatedAt    time.Time
	ChangedAt    *time.Time
}
