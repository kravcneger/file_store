package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/file_store/internal/storage"
)

type FilesRepo struct {
	DB *sql.DB
}

func NewFilesRepository(db *sql.DB) *FilesRepo {
	return &FilesRepo{
		DB: db,
	}
}

func (r *FilesRepo) Create(ctx context.Context, file *storage.File) (err error) {
	stmt, err := r.DB.PrepareContext(ctx, `INSERT INTO images (UUID, name, created_at, updated_at) VALUES ($1,$2,$3,$4)`)
	if err != nil {
		return err
	}
	time := time.Now()
	_, err = stmt.ExecContext(ctx, file.UUID, file.Name, time, time)
	return err
}
