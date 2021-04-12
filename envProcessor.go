package inplaceenvsubst

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type EnvFoundEvent struct {
	start int
	end   int
	env   []rune
	line  []rune
}

type LineProcessingFinishedEvent struct {
	processedContent []string
	filePath         string
}

type EnvFoundResult struct {
	newLine []rune
	newIdx  int
}

type EnvProcessor interface {
	envFound(event *EnvFoundEvent) *EnvFoundResult
	lineProcessingFinished(*LineProcessingFinishedEvent) error
}

type ReplaceEnvVariables struct{}

func (p *ReplaceEnvVariables) envFound(e *EnvFoundEvent) *EnvFoundResult {
	env := string(e.env)
	value, exists := os.LookupEnv(env)
	if !exists {
		// this processor skips missing environment variables
		return &EnvFoundResult{e.line, e.end}
	}
	enValueRunes := []rune(value)
	newLine := append(e.line[:e.start], enValueRunes...)
	newLine = append(newLine, e.line[e.end:]...)
	newIdx := e.start + len(enValueRunes)
	return &EnvFoundResult{newLine, newIdx}
}

func (p *ReplaceEnvVariables) lineProcessingFinished(e *LineProcessingFinishedEvent) error {
	output := strings.Join(e.processedContent, "\n")
	return ioutil.WriteFile(e.filePath, []byte(output), os.ModeAppend)
}

type ProcessorWithValidation struct {
	processor           EnvProcessor
	missingEnvVariables []string
}

func (p *ProcessorWithValidation) envFound(e *EnvFoundEvent) *EnvFoundResult {
	env := string(e.env)
	value, exists := os.LookupEnv(env)
	if !exists {
		p.missingEnvVariables = append(p.missingEnvVariables, value)
	}
	return p.processor.envFound(e)
}

func (p *ProcessorWithValidation) lineProcessingFinished(e *LineProcessingFinishedEvent) error {
	if len(p.missingEnvVariables) == 0 {
		return p.processor.lineProcessingFinished(e)
	} else {
		return errors.New(fmt.Sprintf("Missing environment variables: %v", p.missingEnvVariables))
	}
}
