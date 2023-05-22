package upload

import (
	"context"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/file_store/internal/storage"
	uploadpb "github.com/file_store/proto"
	"github.com/google/uuid"
)

type FilesRepo interface {
	Create(context.Context, *storage.File) (err error)
}

type Server struct {
	storage storage.Manager
	repo    FilesRepo
	uploadpb.UnimplementedUploadServiceServer
}

func NewServer(storage storage.Manager, repo FilesRepo) Server {
	return Server{
		storage: storage,
		repo:    repo,
	}
}

func (s Server) Upload(stream uploadpb.UploadService_UploadServer) error {
	uuid := uuid.New()
	file := storage.NewFile(uuid.String())

	err := s.repo.Create(stream.Context(), file)
	if err != nil {
		return err
	}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			if err := s.storage.WriteToFile(file); err != nil {
				return status.Error(codes.Internal, err.Error())
			}

			return stream.SendAndClose(&uploadpb.UploadResponse{Uuid: uuid.String()})
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		if err := file.Write(req.GetChunk()); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

	}
}

func (s Server) Download(req *uploadpb.DownloadRequest, stream uploadpb.UploadService_DownloadServer) error {
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

		if err := stream.Send(&uploadpb.DownloadResponse{Chunk: buf[:num]}); err != nil {
			return err
		}
	}

	return nil

}
