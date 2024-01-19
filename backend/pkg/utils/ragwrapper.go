package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		WithPayload: &pb.WithPayloadSelector{
			SelectorOptions: &pb.WithPayloadSelector_Enable{
				Enable: true,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	results := searchResults.GetResult()
	parsedResults := make([]models.PredictedSolutionSource, 0)
	for _, result := range results {
		payload := result.GetPayload()
		parsedResult := models.PredictedSolutionSource{
			Url:            payload["url"].GetStringValue(),
			Title:          payload["title"].GetStringValue(),
			Description:    payload["description"].GetStringValue(),
			FeaturedAnswer: payload["featuredAnswer"].GetStringValue(),
		}
		parsedResults = append(parsedResults, parsedResult)
	}
	fmt.Println(len(parsedResults))
	fmt.Println(parsedResults)
	return parsedResults
}

func (rw *RagWrapper) Tokenize(input string) []float32 {
	payload := map[string]interface{}{
		"inputs": append(make([]string, 0), input),
		"options": map[string]interface{}{
			"wait_for_model": true,
		},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST",
		rw.hfw.Url+"/pipeline/feature-extraction/sentence-transformers/all-MiniLM-L12-v2",
		bytes.NewBuffer(jsonData))
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

	var tokenized [][]float32
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &tokenized)
	if err != nil {
		panic(err)
	}
	return tokenized[0]
}
