package domain

type FolderId struct {
	id int
}

func NewFolderId(id int) FolderId {
	return FolderId{id: id}
}

func (f FolderId) Value() uint {
	if f.id <= 0 {
		return 0
	}

	return uint(f.id)
}
