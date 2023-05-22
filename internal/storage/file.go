package storage

import (
	"bytes"
)

type File struct {
	UUID   string
	Name   string
	buffer *bytes.Buffer
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
