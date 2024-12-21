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
	"time"

	"github.com/spf13/viper"
)

type PromptBuilder struct {
	UserInput       string
	Mood            string
	DateOfBirth     string
	PlaceOfBirth    string
	CurrentDeath    string
	TimeOfBirth     string
	CurrentDate     string
	CurrentLocation string
}

type AstroRequestBody struct {
	PlaceOfBirth    string `json:"placeOfBirth"`
	DateOfBirth     string `json:"dateOfBirth"`
	TimeOfBirth     string `json:"timeOfBirth"`
	CurrentLocation string `json:"currentLocation"`
}

func NewPromptBuilder(UserInput, mood string) (*PromptBuilder, error) {
	mood = strings.Split(mood, "-")[1]
	promptBuilder := &PromptBuilder{
		UserInput: UserInput,
		Mood:      mood,
	}
	return promptBuilder, nil
}

func NewAstroPromptBuilder(astroRequestBody AstroRequestBody) (*PromptBuilder, error) {
	promptBuilder := &PromptBuilder{
		DateOfBirth:     astroRequestBody.DateOfBirth,
		PlaceOfBirth:    astroRequestBody.PlaceOfBirth,
		TimeOfBirth:     astroRequestBody.TimeOfBirth,
		CurrentLocation: astroRequestBody.CurrentLocation,
	}
	return promptBuilder, nil
}
func (pb *PromptBuilder) BuildPrompt(file string) (*types.Prompt, error) {
	variables := map[string]interface{}{
		"user_input": pb.UserInput,
		"mood":       pb.Mood,
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

func (pb *PromptBuilder) BuildAstroPrompt(file string) (*types.Prompt, error) {
	date := time.Now().Format("2006-01-02")
	variables := map[string]interface{}{
		"dob":              pb.DateOfBirth,
		"birth_location":   pb.PlaceOfBirth,
		"current_date":     date,
		"time_of_birth":    pb.TimeOfBirth,
		"current_location": pb.CurrentLocation,
	}
	settings := GetAstroSettings()

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

//go:embed *.toml
var settings embed.FS

var (
	onceEmo     sync.Once
	instanceEmo *SingletonSettings

	onceAstro     sync.Once
	instanceAstro *SingletonSettings
)

// NewSingletonSettings initializes the singleton settings instance
func NewSingletonSettings() *SingletonSettings {
	onceEmo.Do(func() {
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

		instanceEmo = &SingletonSettings{
			viper: v,
		}
	})
	return instanceEmo
}

func NewAstroSingletonSettings() *SingletonSettings {
	onceAstro.Do(func() {

		settingsFiles := []string{
			"astro_gpt.toml",
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

		instanceAstro = &SingletonSettings{
			viper: v,
		}
	})
	return instanceAstro
}

// GetSettings returns the singleton settings instance
func GetSettings() *viper.Viper {
	return NewSingletonSettings().viper
}

func GetAstroSettings() *viper.Viper {
	return NewAstroSingletonSettings().viper
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
