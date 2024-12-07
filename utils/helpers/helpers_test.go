package helpers

import (
    "context"
    "testing"
    "emogpt/types"
    "os"
)


// Test generated using Keploy
func TestGenerateChat_ContextCanceledBeforeAI(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    cancel() // Cancel the context immediately

    prompt := types.Prompt{User: "Hello", System: "System message"}
    result, err := GenerateChat(ctx, prompt)

    if err == nil {
        t.Errorf("Expected error due to canceled context, got nil")
    }
    if result != "" {
        t.Errorf("Expected empty result due to canceled context, got %s", result)
    }
}

// Test generated using Keploy
func TestGenerateChat_CallAiError(t *testing.T) {
    ctx := context.Background()

    prompt := types.Prompt{User: "Hello", System: "System message"}
    // Simulate CallAi returning an error by using a mock or a stub
    // For this example, assume CallAi is modified to return an error
    result, err := GenerateChat(ctx, prompt)

    if err == nil {
        t.Errorf("Expected error from CallAi, got nil")
    }
    if result != "" {
        t.Errorf("Expected empty result due to CallAi error, got %s", result)
    }
}


// Test generated using Keploy
func TestUnmarshalYaml_InvalidInput(t *testing.T) {
    invalidYaml := "invalid: yaml: content"
    result, err := unmarshalYaml(invalidYaml)

    if err == nil {
        t.Errorf("Expected error for invalid YAML, got nil")
    }
    if result != "" {
        t.Errorf("Expected empty result for invalid YAML, got %s", result)
    }
}


// Test generated using Keploy
func TestCallAi_MissingAPIKey(t *testing.T) {
    ctx := context.Background()
    completionParams := CompletionParams{}
    aiRequest := AIRequest{Prompt: types.Prompt{User: "Hello"}}

    // Unset the API_KEY environment variable
    os.Setenv("API_KEY", "")

    result, err := CallAi(ctx, completionParams, aiRequest)

    if err == nil {
        t.Errorf("Expected error due to missing API key, got nil")
    }
    if result != "" {
        t.Errorf("Expected empty result due to missing API key, got %s", result)
    }
}

