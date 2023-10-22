package dto

import "time"

type FileMeta struct {
	Name         string
	Extension    string
	OriginalPath string
	Size         uint32
	CreatedAt    time.Time
	ChangedAt    *time.Time
}
