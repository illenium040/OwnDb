package domain

import "time"

type Folder struct {
	id             int
	parentFolderId *int
	name           string
	dtCreated      time.Time
	dtChanged      time.Time
}

func NewFolder(id int, parentFolderId *int, name string, dtCreated time.Time, dtChanged time.Time) Folder {
	return Folder{
		id:             id,
		parentFolderId: parentFolderId,
		name:           name,
		dtCreated:      dtCreated,
		dtChanged:      dtChanged,
	}
}

func (f Folder) Id() int {
	return f.id
}

func (f Folder) ParentFolderId() *int {
	return f.parentFolderId
}

func (f Folder) Name() string {
	return f.name
}

func (f Folder) DtCreated() time.Time {
	return f.dtCreated
}

func (f Folder) DtChanged() time.Time {
	return f.dtChanged
}
