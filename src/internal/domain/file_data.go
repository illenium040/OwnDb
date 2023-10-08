package domain

type FileData struct {
	id   uint
	hash string
}

func NewFileData(id uint, hash string) FileData {
	return FileData{id: id, hash: hash}
}

func (f FileData) Id() uint {
	return f.id
}

func (f FileData) Hash() string {
	return f.hash
}
