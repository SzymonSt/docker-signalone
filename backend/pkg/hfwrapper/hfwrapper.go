package hfwrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var DEFAULT_ZEPHYR_PROMPT_TEMPLATE = "<|user|>Summarize these logs and generate a single paragraph summary of what is happening in these logs in high technical detail: \n%s\n</s>\n<|assistant|>"

type HfWrapper struct {
	Url            string
	ApiKey         string
	Temperature    float64
	TopK           int
	TopP           float64
	DoSample       bool
	MaxNewTokens   int
	PromptTemplate string
}

func NewHfWrapper(
	url string,
	apiKey string,
	temperature float64,
	topK int, topP float64,
	doSample bool, maxNewTokens int) *HfWrapper {
	return &HfWrapper{
		Url:            url,
		ApiKey:         apiKey,
		Temperature:    temperature,
		TopK:           topK,
		TopP:           topP,
		DoSample:       doSample,
		MaxNewTokens:   maxNewTokens,
		PromptTemplate: DEFAULT_ZEPHYR_PROMPT_TEMPLATE,
	}
}

func (hfw *HfWrapper) Predict(logs string) string {
	prompt := fmt.Sprintf(hfw.PromptTemplate, logs)
	payload := map[string]interface{}{
		"prompt": prompt,
		"parameters": map[string]interface{}{
			"temperature":    hfw.Temperature,
			"top_k":          hfw.TopK,
			"top_p":          hfw.TopP,
			"do_sample":      hfw.DoSample,
			"max_new_tokens": hfw.MaxNewTokens,
		},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(hfw.Url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+hfw.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}
