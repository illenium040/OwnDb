package transport

import (
	"own-db/src/internal/domain"
	"time"
)

type fileMeta struct {
	Id           uint       `json:"id"`
	DataId       uint       `json:"dataId"`
	FolderId     *int       `json:"folderId"`
	Name         string     `json:"name"`
	Extension    string     `json:"extension"`
	OriginalPath string     `json:"originalPath"`
	Size         uint32     `json:"size"`
	CreatedAt    time.Time  `json:"createdAt"`
	ChangedAt    *time.Time `json:"changedAt"`
}

func fileMetaFromDomain(fm domain.FileMeta) fileMeta {
	return fileMeta{
		Id:           fm.Id(),
		DataId:       fm.DataId(),
		FolderId:     fm.FolderId(),
		Name:         fm.Name(),
		Extension:    fm.Extension(),
		OriginalPath: fm.OriginalPath(),
		Size:         fm.Size(),
		CreatedAt:    fm.CreatedAt(),
		ChangedAt:    fm.ChangedAt(),
	}
}
