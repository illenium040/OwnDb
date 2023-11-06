package repository

import (
	"own-db/src/internal/domain"
	"time"
)

type folder struct {
	Id             int
	ParentFolderId int
	Name           string
	CreatedAt      time.Time
	ChangedAt      *time.Time
}

func folderToDomain(f folder) domain.Folder {
	return domain.NewFolder(
		domain.NewFolderId(f.Id),
		domain.NewFolderId(f.ParentFolderId),
		f.Name,
		f.CreatedAt,
		f.ChangedAt,
	)
}

func folderFromDomain(f domain.Folder) folder {
	return folder{
		Id: f.Id().Value(),
	}
}
