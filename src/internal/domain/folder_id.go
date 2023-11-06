package domain

type FolderId struct {
	id int
}

func NewFolderId(id int) FolderId {
	return FolderId{id: id}
}

func (f FolderId) Value() int {
	if f.id <= 0 {
		return 0
	}

	return f.id
}

func (f FolderId) IsRoot() bool {
	return f.id == 0
}
