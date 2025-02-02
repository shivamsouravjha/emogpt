package helpers

import (
    "context"
    "testing"
    "emogpt/types"
    "os"
)


// Test generated using Keploy
func TestGenerateChat_ContextCanceled_Error(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    cancel() // Immediately cancel the context

    prompt := types.Prompt{
        System: "Test system message",
        User:   "Test user message",
    }

    _, err := GenerateChat(ctx, prompt)
    if err == nil {
        t.Errorf("Expected error due to canceled context, got nil")
    }
}

// Test generated using Keploy
func TestUnmarshalYaml_ValidInput(t *testing.T) {
    yamlStr := "response: Test response\nzodiac: Test zodiac"

    result, err := unmarshalYaml(yamlStr)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }

    if result.Response != "Test response" || result.Zodiac != "Test zodiac" {
        t.Errorf("Unexpected result: %+v", result)
    }
}


// Test generated using Keploy
func TestUnmarshalYaml_InvalidInput(t *testing.T) {
    yamlStr := "invalid_yaml: ::::"

    _, err := unmarshalYaml(yamlStr)
    if err == nil {
        t.Errorf("Expected error for invalid YAML input, got nil")
    }
}


// Test generated using Keploy
func TestCallAi_MissingAPIKey(t *testing.T) {
    ctx := context.Background()

    os.Setenv("API_KEY", "") // Ensure API key is missing
    defer os.Unsetenv("API_KEY")

    completionParams := CompletionParams{}
    aiRequest := AIRequest{}

    _, err := CallAi(ctx, completionParams, aiRequest)
    if err == nil {
        t.Errorf("Expected error due to missing API key, got nil")
    }
}

