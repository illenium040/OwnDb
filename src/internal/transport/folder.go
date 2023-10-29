package transport

import (
	"own-db/src/internal/domain"
	"time"
)

type folder struct {
	Id             int       `json:"id"`
	ParentFolderId *int      `json:"parentFolderId"`
	Name           string    `json:"name"`
	DtCreated      time.Time `json:"dtCreated"`
	DtChanged      time.Time `json:"dtChanged"`
}

func folderFromDomain(domainFolder domain.Folder) folder {
	return folder{
		Id:             domainFolder.Id(),
		ParentFolderId: domainFolder.ParentFolderId(),
		Name:           domainFolder.Name(),
		DtCreated:      domainFolder.DtCreated(),
		DtChanged:      domainFolder.DtChanged(),
	}
}
