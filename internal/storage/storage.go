package storage

import (
	"os"
)

type Manager interface {
	WriteToFile(file *File) error
	Open(uuid string) (*os.File, error)
}

var _ Manager = &Storage{}

type Storage struct {
	Dir string
}

func New(dir string) Storage {
	return Storage{
		Dir: dir,
	}
}

func (s Storage) Open(uuid string) (*os.File, error) {
	f, err := os.Open(s.Dir + uuid)
	return f, err
}

func (s Storage) WriteToFile(file *File) error {
	f, err := os.OpenFile(s.Dir+file.UUID, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.Write(file.buffer.Bytes())
	return err
}
