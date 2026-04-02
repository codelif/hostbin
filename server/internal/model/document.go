package model

import "time"

type Document struct {
	Slug      string
	Content   []byte
	SHA256    string
	SizeBytes int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DocumentMeta struct {
	Slug      string
	SHA256    string
	SizeBytes int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
