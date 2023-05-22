package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/file_store/internal/repo"
	"github.com/file_store/internal/storage"
	"github.com/file_store/internal/upload"
	uploadpb "github.com/file_store/proto"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// load Enviroments
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Db connection

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("DB is successfully connected!")
	repo := repo.NewFilesRepository(db)

	// Initialise TCP listener.
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()

	// Bootstrap upload server.
	uplSrv := upload.NewServer(storage.New(os.Getenv("FILES_DIRECTORY")), repo)

	// Bootstrap gRPC server.
	rpcSrv := grpc.NewServer()

	// Register and start gRPC server.
	uploadpb.RegisterUploadServiceServer(rpcSrv, uplSrv)
	log.Fatal(rpcSrv.Serve(lis))
}
