package processors

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Used for processing environment variables placeholders
type EnvPlaceholderProcessor interface {
	// Called when an env placeholder is found in a given line
	EnvPlaceholderFound(e *EnvPlaceholderFoundEvent) *EnvPlaceholderResult
	// Called when a file is completely processed
	FileProcessingFinished(e *FileProcessingFinishedEvent) error
}

// Processor used for replacing the env variables placeholder with their actual value
// The replacement is done at the end, so decorators (like the EnvVariableExistsDecorator) can be used to either
// replace the whole file with the found env variables or not
type ReplaceEnvVariablesProcessor struct{}

func (p *ReplaceEnvVariablesProcessor) EnvPlaceholderFound(e *EnvPlaceholderFoundEvent) *EnvPlaceholderResult {
	env := string(e.Placeholder)
	value, exists := os.LookupEnv(env)
	if !exists {
		return &EnvPlaceholderResult{e.Line, e.End}
	}
	enValueRunes := []rune(value)
	var newLine []rune
	newLine = append(newLine, e.Line[:e.Start]...)
	newLine = append(newLine, enValueRunes...)
	newLine = append(newLine, e.Line[e.End:]...)
	newIdx := e.Start + len(enValueRunes)
	return &EnvPlaceholderResult{newLine, newIdx}
}

func (p *ReplaceEnvVariablesProcessor) FileProcessingFinished(e *FileProcessingFinishedEvent) error {
	output := strings.Join(e.ProcessedContent, "\n")
	return ioutil.WriteFile(e.FilePath, []byte(output), os.ModeAppend)
}

// Decorator that also validates if env variable exists
type EnvVariableExistsDecorator struct {

	// the original processor which will be decorated
	Processor EnvPlaceholderProcessor

	// collections of strings representing the missing env variables found
	MissingEnvVariables []string
}

func (p *EnvVariableExistsDecorator) EnvPlaceholderFound(e *EnvPlaceholderFoundEvent) *EnvPlaceholderResult {
	env := string(e.Placeholder)
	value, exists := os.LookupEnv(env)
	if !exists {
		p.MissingEnvVariables = append(p.MissingEnvVariables, value)
	}
	return p.Processor.EnvPlaceholderFound(e)
}

func (p *EnvVariableExistsDecorator) FileProcessingFinished(e *FileProcessingFinishedEvent) error {
	if len(p.MissingEnvVariables) == 0 {
		return p.Processor.FileProcessingFinished(e)
	} else {
		return errors.New(fmt.Sprintf("Missing environment variables: %v", p.MissingEnvVariables))
	}
}
