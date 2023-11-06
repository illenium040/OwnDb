package domain

import "time"

type Folder struct {
	id             FolderId
	parentFolderId FolderId
	name           string
	createdAt      time.Time
	changedAt      *time.Time
}

func NewFolder(id FolderId, parentFolderId FolderId, name string, dtCreated time.Time, dtChanged *time.Time) Folder {
	return Folder{
		id:             id,
		parentFolderId: parentFolderId,
		name:           name,
		createdAt:      dtCreated,
		changedAt:      dtChanged,
	}
}

func (f Folder) Id() FolderId {
	return f.id
}

func (f Folder) ParentFolderId() FolderId {
	return f.parentFolderId
}

func (f Folder) Name() string {
	return f.name
}

func (f Folder) CreatedAt() time.Time {
	return f.createdAt
}

func (f Folder) ChangedAt() *time.Time {
	return f.changedAt
}
