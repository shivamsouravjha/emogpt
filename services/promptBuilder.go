package services

import (
	"bytes"
	"embed"
	"emogpt/types"
	"fmt"
	"html"
	"html/template"
	"log"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type PromptBuilder struct {
	UserInput string
}

func NewPromptBuilder(UserInput string) (*PromptBuilder, error) {
	promptBuilder := &PromptBuilder{
		UserInput: UserInput,
	}
	return promptBuilder, nil
}
func (pb *PromptBuilder) BuildPrompt(file string) (*types.Prompt, error) {
	variables := map[string]interface{}{
		"user_input": pb.UserInput,
	}

	settings := GetSettings()

	prompt := &types.Prompt{}
	systemPrompt, err := renderTemplate(settings.GetString(file+".system"), variables)
	if err != nil {
		prompt.System = ""
		prompt.User = ""
		return prompt, fmt.Errorf("Error rendering system prompt: %v", err)
	}

	userPrompt, err := renderTemplate(settings.GetString(file+".user"), variables)
	if err != nil {
		prompt.System = ""
		prompt.User = ""
		return prompt, fmt.Errorf("Error rendering user prompt: %v", err)
	}
	prompt.System = systemPrompt
	userPrompt = html.UnescapeString(userPrompt)
	prompt.User = userPrompt
	return prompt, nil
}

// SingletonSettings manages the singleton instance of the configuration settings
type SingletonSettings struct {
	viper *viper.Viper
}

var instance *SingletonSettings
var once sync.Once

//go:embed *.toml
var settings embed.FS

// NewSingletonSettings initializes the singleton settings instance
func NewSingletonSettings() *SingletonSettings {
	once.Do(func() {

		settingsFiles := []string{
			"emo_gpt.toml",
		}

		v := viper.New()
		v.SetConfigType("toml")
		for _, file := range settingsFiles {
			fileContent, err := settings.ReadFile(file)
			if err != nil {
				log.Fatalf("Failed to read settings file %s: %v", file, err)
			}
			v.SetConfigFile(file)
			if err := v.MergeConfig(bytes.NewBuffer(fileContent)); err != nil {
				log.Fatalf("Error loading config file : %v", err)
			}
		}

		instance = &SingletonSettings{
			viper: v,
		}
	})
	return instance
}

// GetSettings returns the singleton settings instance
func GetSettings() *viper.Viper {
	return NewSingletonSettings().viper
}

func renderTemplate(templateText string, variables map[string]interface{}) (string, error) {
	funcMap := template.FuncMap{
		"trim": strings.TrimSpace,
	}
	tmpl, err := template.New("prompt").Funcs(funcMap).Parse(templateText)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, variables)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
