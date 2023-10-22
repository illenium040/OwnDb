package dto

import "time"

type Folder struct {
	ParentFolderId *int
	Name           string
	DtCreated      time.Time
	DtChanged      time.Time
}
