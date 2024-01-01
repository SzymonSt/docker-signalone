package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RagWrapper struct {
	vectorDbClient pb.PointsClient
	hfw            *HfWrapper
}

func NewRagWrapper(vectorDbAddr string, hfw *HfWrapper) *RagWrapper {
	conn, err := grpc.Dial(vectorDbAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return &RagWrapper{
		vectorDbClient: pb.NewPointsClient(conn),
		hfw:            hfw,
	}
}

func (rw *RagWrapper) Predict(input []float32) []string {
	ctx := context.Background()
	collectionName := "resources"
	searchResults, err := rw.vectorDbClient.Search(ctx, &pb.SearchPoints{
		CollectionName: collectionName,
		Vector:         input,
		Limit:          5,
	})
	if err != nil {
		panic(err)
	}
	results := searchResults.GetResult()
	parsedResults := make([]string, 0)
	for _, result := range results {
		parsedResults = append(parsedResults, result.Payload["solution"].GetStringValue())
	}
	return parsedResults
}

func (rw *RagWrapper) Tokenize(input string) []float32 {
	payload := map[string]interface{}{
		"prompt": input,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(rw.hfw.Url+"/sentence-transformers/all-MiniLM-L12-v2", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+rw.hfw.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var tokenized []float32
	if err := json.NewDecoder(resp.Body).Decode(&tokenized); err != nil {
		panic(err)
	}
	return tokenized
}
