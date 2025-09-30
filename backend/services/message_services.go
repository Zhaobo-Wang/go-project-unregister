package services

import (
	"bytes"
	"encoding/json"
	"github.com/Zhaobo-Wang/go-project-unregister/models"
	"io/ioutil"
	"net/http"
)

type LLMRequest struct {
	Model     string           `json:"model"`
	Messages  []models.Message `json:"messages"`
	MaxTokens int              `json:"max_tokens,omitempty"`
}

type LLMResponse struct {
	// Define the response structure based on the API documentation
}

func SendMessageToLLMV1(message models.Message) (LLMResponse, error) {
	url := "https://api.siliconflow.cn/v1/messages"
	payload := LLMRequest{
		Model: "deepseek-ai/DeepSeek-V3.1",
		Messages: []models.Message{
			message,
		},
		MaxTokens: 8192,
	}

	return sendRequest(url, payload)
}

func SendMessageToLLMV2(message models.Message) (LLMResponse, error) {
	url := "https://api.siliconflow.cn/v1/chat/completions"
	payload := LLMRequest{
		Model: "Qwen/QwQ-32B",
		Messages: []models.Message{
			message,
		},
	}

	return sendRequest(url, payload)
}

func sendRequest(url string, payload LLMRequest) (LLMResponse, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return LLMResponse{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return LLMResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+"sk-xvcdofkkvlndxowhiwdptngunifwxxvclccffxsapquawiur") //token
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return LLMResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return LLMResponse{}, err
	}

	var response LLMResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return LLMResponse{}, err
	}

	return response, nil
}
