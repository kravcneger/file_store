package storage

import (
	"bytes"
	"time"
)

type File struct {
	UUID      string    `pg:"uuid"`
	Name      string    `pg:"name"`
	CreatedAt time.Time `pg:"created_at"`
	UpdatedAt time.Time `pg:"updated_at"`
	buffer    *bytes.Buffer
}

func NewFile(uuid string) *File {
	return &File{
		UUID:   uuid,
		buffer: &bytes.Buffer{},
	}
}

func (f *File) Write(chunk []byte) error {
	_, err := f.buffer.Write(chunk)

	return err
}
