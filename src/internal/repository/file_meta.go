package repository

import (
	"own-db/src/internal/domain"
	"time"
)

type fileMeta struct {
	Id           uint
	DataId       uint
	Name         string
	Extension    string
	OriginalPath string
	Size         uint32
	CreatedAt    time.Time
	ChangedAt    *time.Time
}

func fileMetaToDomain(fm fileMeta) domain.FileMeta {
	return domain.NewFileMeta(
		fm.Id,
		fm.DataId,
		fm.Name,
		fm.Extension,
		fm.OriginalPath,
		fm.Size,
		fm.CreatedAt,
		fm.ChangedAt,
	)
}

func fileMetaFromDomain(fm domain.FileMeta) (res fileMeta) {
	return fileMeta{
		Id:           fm.Id(),
		DataId:       fm.DataId(),
		Name:         fm.Name(),
		Extension:    fm.Extension(),
		OriginalPath: fm.OriginalPath(),
		Size:         fm.Size(),
		CreatedAt:    fm.CreatedAt(),
		ChangedAt:    fm.ChangedAt(),
	}
}
