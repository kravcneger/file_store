package upload

import (
	"context"
	"errors"
	"io"
	"path/filepath"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/file_store/internal/storage"
	storepb "github.com/file_store/proto"
	"github.com/google/uuid"
)

const (
	UploadWorkerPool  = 10
	GetListWorkerPool = 100
)

type FilesRepo interface {
	Create(context.Context, *storage.File) (err error)
	GetList(context.Context) ([]storage.File, error)
}

type Server struct {
	uploadWorkerGroup chan struct{}
	readWorkerGroup   chan struct{}
	storage           storage.Manager
	repo              FilesRepo
	storepb.UnimplementedUploadServiceServer
}

func NewServer(storage storage.Manager, repo FilesRepo) *Server {
	s := Server{
		storage: storage,
		repo:    repo,
	}
	s.uploadWorkerGroup = make(chan struct{}, UploadWorkerPool)
	s.readWorkerGroup = make(chan struct{}, GetListWorkerPool)
	return &s
}

func (s *Server) Upload(stream storepb.UploadService_UploadServer) error {
	s.uploadWorkerGroup <- struct{}{}
	defer func() { <-s.uploadWorkerGroup }()

	uuid := uuid.New()
	file := storage.NewFile(uuid.String())

	// В первом партишине получаем название файла
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	// For debug
	//time.Sleep(30 * time.Second)
	// Не стал добавлять строгую проверку на mime type
	// В продакшене бы добавил.
	switch filepath.Ext(req.GetName()) {
	case ".jpg", ".png", ".jpeg", ".gif":
	default:
		return errors.New("Invalid file")
	}

	file.Name = filepath.Base(req.GetName())

	err = s.repo.Create(stream.Context(), file)
	if err != nil {
		return err
	}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			if err := s.storage.WriteToFile(file); err != nil {
				return status.Error(codes.Internal, err.Error())
			}

			return stream.SendAndClose(&storepb.UploadResponse{Uuid: uuid.String()})
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		if err := file.Write(req.GetChunk()); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

	}
}

func (s *Server) Download(req *storepb.DownloadRequest, stream storepb.UploadService_DownloadServer) error {
	s.readWorkerGroup <- struct{}{}
	defer func() { <-s.readWorkerGroup }()

	file, err := s.storage.Open(req.GetUuid())
	defer file.Close()
	if err != nil {
		return err
	}

	// Maximum 1mb size per stream.
	buf := make([]byte, 1024*1024)

	for {
		num, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&storepb.DownloadResponse{Chunk: buf[:num]}); err != nil {
			return err
		}
	}

	return nil

}

func (s *Server) GetList(req *storepb.GetListRequest, stream storepb.UploadService_GetListServer) error {
	s.uploadWorkerGroup <- struct{}{}
	defer func() { <-s.uploadWorkerGroup }()

	files, err := s.repo.GetList(stream.Context())
	if err != nil {
		return err
	}

	for _, file := range files {
		err := stream.Send(&storepb.GetListResponse{
			Name:      file.Name,
			CreatedAt: file.CreatedAt.String(),
			UpdatedAt: file.UpdatedAt.String(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
