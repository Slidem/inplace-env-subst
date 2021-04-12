package inplaceenvsubst

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

var replaceEnvsProcessor = &ReplaceEnvVariables{}

type ErrorListener interface {
	ErrorFound(filepath string, err error)
}

type ConsoleErrorListener struct {}

func (c *ConsoleErrorListener) ErrorFound(filepath string, err error) {
	log.Print(err)
}

func ProcessFiles(filePaths []string, config *Config) {
	if config.RunInParallel {
		ProcessInParallel(filePaths, config)
	} else {
		ProcessSequentially(filePaths, config)
	}
}

func ProcessSequentially(filePaths []string, config *Config) {

	for _, p := range filePaths {
		processFile(p, config)
	}
}

func ProcessInParallel(filePaths []string, config *Config) {

	var wg sync.WaitGroup
	wg.Add(len(filePaths))
	for _, p := range filePaths {
		go processFileAsync(&wg, p, config)
	}
	wg.Wait()
}

func processFileAsync(w *sync.WaitGroup, filePath string, config *Config) {
	defer w.Done()
	processFile(filePath, config)
}

func processFile(filePath string, config *Config) {

	var processor EnvProcessor
	processor = replaceEnvsProcessor
	if config.FailOnMissingVariables {
		processor = &ProcessorWithValidation{
			processor: processor,
		}
	}
	err := process(filePath, &LineProcessor{[]string{}, replaceEnvsProcessor})
	if err != nil && config.ErrorListener != nil {
		config.ErrorListener.ErrorFound(filePath, err)
	}
}

func process(fp string, s *LineProcessor) error {
	input, err := ioutil.ReadFile(fp)
	if err != nil {
		panic(err)
	}
	fileLines := strings.Split(string(input), "\n")
	for number, line := range fileLines {
		s.processLine(number, line)
	}
	return s.processFinishedForPath(fp)
}

type LineProcessor struct {
	processedLines []string
	envProcessor   EnvProcessor
}

func (p *LineProcessor) processLine(number int, line string) {

	lineRunes := []rune(line)
	idx := 0
	for idx < len(lineRunes) {
		for idx < len(lineRunes) && lineRunes[idx] != '$' {
			idx++
		}
		for idx < len(lineRunes) && lineRunes[idx] == '$' {
			idx++
		}
		if idx == len(lineRunes) {
			break
		}
		if lineRunes[idx] == '{' {
			start := idx - 1
			var envRunes []rune
			for {
				idx++
				if idx == len(lineRunes) {
					panic(fmt.Sprintf("Expected } but found newline on line number %d", number+1))
				}
				if lineRunes[idx] == '}' {
					idx++
					for idx < len(lineRunes) && lineRunes[idx] == '}' {
						envRunes = append(envRunes, '}')
						idx++
					}
					end := idx
					if len(envRunes) != 0 {
						response := p.envProcessor.envFound(&EnvFoundEvent{
							start: start,
							end:   end,
							env:   envRunes,
							line:  lineRunes,
						})

						lineRunes = response.newLine
						idx = response.newIdx
					}
					break
				}
				envRunes = append(envRunes, lineRunes[idx])
			}
		}
	}
	p.processedLines = append(p.processedLines, string(lineRunes))
}

func (p *LineProcessor) processFinishedForPath(filePath string) error {
	return p.envProcessor.lineProcessingFinished(
		&LineProcessingFinishedEvent{
			processedContent: p.processedLines,
			filePath:         filePath,
		})
}
