package upload

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/file_store/internal/storage"
	uploadpb "github.com/file_store/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	client uploadpb.UploadServiceClient
}

func NewClient(conn grpc.ClientConnInterface) Client {
	return Client{
		client: uploadpb.NewUploadServiceClient(conn),
	}
}

func (c Client) Upload(ctx context.Context, file string) (string, error) {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Second))
	defer cancel()

	stream, err := c.client.Upload(ctx)
	if err != nil {
		return "", err
	}

	fil, err := os.Open(file)
	if err != nil {
		return "", err
	}

	// Maximum 1mb size per stream.
	buf := make([]byte, 1024*1024)

	for {
		num, err := fil.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if err := stream.Send(&uploadpb.UploadRequest{Chunk: buf[:num]}); err != nil {
			return "", err
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return "", err
	}

	return res.GetUuid(), nil
}

func (c Client) Download(ctx context.Context, uuid string, path string) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Second))
	defer cancel()

	stream, err := c.client.Download(ctx, &uploadpb.DownloadRequest{Uuid: uuid})
	if err != nil {
		return err
	}

	store := storage.New(path)

	file := storage.NewFile(uuid)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			if err := store.WriteToFile(file); err != nil {
				return status.Error(codes.Internal, err.Error())
			}

			return stream.CloseSend()
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		if err := file.Write(req.GetChunk()); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

	}

}
