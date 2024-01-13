package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"signalone/pkg/models"

	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RagWrapper struct {
	vectorDbClient pb.PointsClient
	collectionName string
	retrivalLimit  uint64
	hfw            *HfWrapper
}

func NewRagWrapper(vectorDbAddr string, hfw *HfWrapper, collectionName string, retrivalLimit uint64) *RagWrapper {
	conn, err := grpc.Dial(vectorDbAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return &RagWrapper{
		vectorDbClient: pb.NewPointsClient(conn),
		collectionName: collectionName,
		retrivalLimit:  retrivalLimit,
		hfw:            hfw,
	}
}

func (rw *RagWrapper) Predict(input []float32) []models.PredictedSolutionSource {
	ctx := context.Background()
	searchResults, err := rw.vectorDbClient.Search(ctx, &pb.SearchPoints{
		CollectionName: rw.collectionName,
		Vector:         input,
		Limit:          uint64(rw.retrivalLimit),
	})
	if err != nil {
		panic(err)
	}
	results := searchResults.GetResult()
	parsedResults := make([]models.PredictedSolutionSource, 0)
	for _, result := range results {
		var parsedResult models.PredictedSolutionSource
		payloadBytes, _ := json.Marshal(result.Payload)
		json.Unmarshal(payloadBytes, &parsedResult)
		parsedResults = append(parsedResults, parsedResult)
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
