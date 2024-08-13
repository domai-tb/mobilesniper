package models

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type OpenAPI struct {
	Info  Info                              `yaml:"info"`
	Paths map[string]map[string]interface{} `yaml:"paths"`
}

type Info struct {
	Title       string `yaml:"title"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

func ParseOpenAPIFile(filePath string) (*OpenAPI, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAPI file: %v", err)
	}

	var openapi OpenAPI
	err = yaml.Unmarshal(data, &openapi)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI file: %v", err)
	}

	return &openapi, nil
}

func ValidateOpenAPIFile(filePath string) (*OpenAPI, error) {
	// Check if the file has a YAML extension
	if !(strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml")) {
		return nil, fmt.Errorf("file %s is not a YAML file", filePath)
	}

	swagger, err := ParseOpenAPIFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("invalid OpenAPI file: %s - %v", filePath, err)
	}

	return swagger, nil
}
