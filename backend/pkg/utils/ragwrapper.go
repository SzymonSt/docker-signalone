package utils

import (
	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RagWrapper struct {
	vectorDbClient pb.QdrantClient
}

func NewRagWrapper(vectorDbAddr string) *RagWrapper {
	conn, err := grpc.Dial(vectorDbAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return &RagWrapper{
		vectorDbClient: pb.NewQdrantClient(conn),
	}
}
