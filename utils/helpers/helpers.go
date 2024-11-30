package helpers

import (
	"bytes"
	"context"
	"emogpt/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

func GenerateChat(ctx context.Context, prompt types.Prompt) (string, error) {
	fmt.Println("Generating chat...")

	select {
	case <-ctx.Done():
		err := ctx.Err()
		zap.L().Error("Context done", zap.Error(err))
		return "", err
	default:
	}

	aiRequest := AIRequest{
		MaxTokens: 4096,
		Prompt:    prompt,
	}
	response, err := CallAi(ctx, CompletionParams{}, aiRequest)
	if err != nil {
		zap.L().Error("Error calling AI", zap.Error(err))
		return "", err
	}

	select {
	case <-ctx.Done():
		err := ctx.Err()
		zap.L().Error("Context done", zap.Error(err))
		return "", err
	default:
	}

	select {
	case <-ctx.Done():
		err := ctx.Err()
		zap.L().Error("Context done", zap.Error(err))
		return "", err
	default:
	}

	chat, err := unmarshalYaml(response)

	if err != nil {
		zap.L().Error("Error unmarshalling yaml", zap.Error(err))
		return "", err
	}

	select {
	case <-ctx.Done():
		err := ctx.Err()
		zap.L().Error("Context done", zap.Error(err))
		return "", err
	default:
	}

	return chat, nil
}

type Response struct {
	Response string `yaml:"response"`
}

func unmarshalYaml(yamlStr string) (string, error) {
	yamlStr = strings.TrimSpace(yamlStr)
	yamlStr = strings.TrimPrefix(yamlStr, "```yaml")
	yamlStr = strings.TrimSuffix(yamlStr, "```")
	var data Response
	err := yaml.Unmarshal([]byte(yamlStr), &data)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling yaml: %s", err)
	}
	return data.Response, nil
}

type AIResponse struct {
	IsSuccess        bool   `json:"isSuccess"`
	Error            string `json:"error"`
	FinalContent     string `json:"finalContent"`
	PromptTokens     int    `json:"promptTokens"`
	CompletionTokens int    `json:"completionTokens"`
	APIKey           string `json:"apiKey"`
}

type AIRequest struct {
	MaxTokens int          `json:"maxTokens"`
	Prompt    types.Prompt `json:"prompt"`
	SessionID string       `json:"sessionId"`
	Iteration int          `json:"iteration"`
}

type CompletionParams struct {
	Model               string    `json:"model"`
	Messages            []Message `json:"messages"`
	MaxTokens           int       `json:"max_tokens,omitempty"`
	MaxCompletionTokens int       `json:"max_completion_tokens,omitempty"`
	Stream              *bool     `json:"stream,omitempty"`
	Temperature         float32   `json:"temperature,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func CallAi(ctx context.Context, completionParams CompletionParams, aiRequest AIRequest) (string, error) {

	var apiBaseURL string

	var apiKey string

	apiBaseURL = os.Getenv("API_BASE_URL")

	if completionParams.Model == "" {
		completionParams.Model = "gpt-4o"
	}

	if len(completionParams.Messages) == 0 {
		var messages []Message
		if aiRequest.Prompt.System == "" {
			messages = []Message{
				{Role: "user", Content: aiRequest.Prompt.User},
			}
		} else {
			messages = []Message{
				{Role: "system", Content: aiRequest.Prompt.System},
				{Role: "user", Content: aiRequest.Prompt.User},
			}
		}

		completionParams.Messages = messages
	}

	requestBody, err := json.Marshal(completionParams)
	if err != nil {
		zap.L().Error("Error marshalling request body", zap.Error(err))
		return "", fmt.Errorf("error marshalling request body: %v", err)
	}

	queryParams := "?api-version=2024-08-01-preview"

	req, err := http.NewRequestWithContext(ctx, "POST", apiBaseURL+"/chat/completions"+queryParams, bytes.NewBuffer(requestBody))
	if err != nil {
		zap.L().Error("Error creating request", zap.Error(err))
		return "", fmt.Errorf("error creating request: %v", err)
	}

	apiKey = os.Getenv("API_KEY")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Error("Error making request", zap.Error(err))
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			zap.L().Error("Error closing response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		zap.L().Error("Unexpected status code", zap.Int("status_code", resp.StatusCode), zap.String("response_body", bodyString))
		return "", fmt.Errorf("unexpected status code: %v, response body: %s", resp.StatusCode, bodyString)
	}

	var contentBuilder strings.Builder
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("Error reading response body", zap.Error(err))
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var chunk types.ModelResponse
	err = json.Unmarshal(bodyBytes, &chunk)
	if err != nil {
		zap.L().Error("Error unmarshalling response body", zap.Error(err))
		return "", fmt.Errorf("error unmarshalling response body: %v", err)
	}

	// Process the data
	if len(chunk.Choices) > 0 {
		if chunk.Choices[0].Message.Content != "" {
			contentBuilder.WriteString(chunk.Choices[0].Message.Content)
		}
	}
	result := contentBuilder.String()

	return result, nil
}
