package main

import (
	"context"
	"flag"
	"log"

	"github.com/file_store/internal/upload"
	"google.golang.org/grpc"
)

func main() {
	// Catch user input.
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatalln("Missing file path")
	}

	// Initialise gRPC connection.
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// Start uploading the file. Error if failed, otherwise echo download URL.
	client := upload.NewClient(conn)

	if len(flag.Args()) < 1 {
		log.Fatalf("fatal: you should specify the command")
	}

	command := flag.Arg(0)
	switch command {
	case "upload":
		res, er := client.Upload(context.Background(), flag.Arg(1))
		log.Println(res)
		log.Println(er)
	case "download":
		err = client.Download(context.Background(), flag.Arg(1), flag.Arg(2))
	case "list":
		err = client.GetList(context.Background())
	}
	log.Println(err)
}
