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

	key := getKey(e)
	value, exists := os.LookupEnv(key)
	if !exists {
		defaultValue, hasDefault := getDefault(e)
		if hasDefault {
			value = defaultValue
		} else {
			return &EnvPlaceholderResult{e.Line, e.End}
		}
	}

	// transform value to runes
	enValueRunes := []rune(value)
	var newLine []rune
	newLine = append(newLine, e.Line[:e.Start]...)
	newLine = append(newLine, enValueRunes...)
	newLine = append(newLine, e.Line[e.End:]...)

	// update index to initial start + len of the value
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
	key := getKey(e)
	_, exists := os.LookupEnv(key)
	if !exists {
		_, hasDefault := getDefault(e)
		if !hasDefault {
			p.MissingEnvVariables = append(p.MissingEnvVariables, key)
		}
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

func getKey(e *EnvPlaceholderFoundEvent) string {
	placeHolderParts := strings.Split(string(e.Placeholder), ":-")
	return placeHolderParts[0]
}

func getDefault(e *EnvPlaceholderFoundEvent) (string, bool) {
	placeHolderParts := strings.Split(string(e.Placeholder), ":-")
	if len(placeHolderParts) != 2 {
		return "", false
	}
	return placeHolderParts[1], true
}
