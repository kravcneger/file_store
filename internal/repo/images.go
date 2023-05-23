package repo

import (
	"context"
	"database/sql"
	"log"
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

func (r *FilesRepo) Create(ctx context.Context, file *storage.File) error {
	stmt, err := r.DB.PrepareContext(ctx, `INSERT INTO images (UUID, name, created_at, updated_at) VALUES ($1,$2,$3,$4)`)
	if err != nil {
		return err
	}
	time := time.Now().UTC()
	_, err = stmt.ExecContext(ctx, file.UUID, file.Name, time, time)
	return err
}

func (r *FilesRepo) GetList(ctx context.Context) ([]storage.File, error) {
	res := []storage.File{}
	stmt, err := r.DB.PrepareContext(ctx, `Select name,created_at,updated_at from images ORDER BY created_at DESC`)
	if err != nil {
		return res, err
	}

	rows, err := stmt.QueryContext(ctx)

	f := storage.File{}
	for rows.Next() {

		if err := rows.Scan(&f.Name, &f.CreatedAt, &f.UpdatedAt); err != nil {
			log.Println(err)
		}

		// Не нашёл в документации как узнать количество записей в ответа запроса
		// Чтобы с первого раза слайс под нужный размер сделать.
		res = append(res, f)
	}
	return res, err
}
