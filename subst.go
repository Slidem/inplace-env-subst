package inplaceenvsubst

import (
	"github.com/Slidem/inplaceenvsubst/processors"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

var replaceEnvsProcessor = &processors.ReplaceEnvVariablesProcessor{}

type ErrorListener interface {
	ErrorFound(filepath string, err error)
}

type ConsoleErrorListener struct{}

func (c *ConsoleErrorListener) ErrorFound(filepath string, err error) {
	log.Print(err)
}

func ProcessFiles(filePaths []string, config *Config) {
	if config.RunInParallel {
		processInParallel(filePaths, config)
	} else {
		processSequentially(filePaths, config)
	}
}

func processSequentially(filePaths []string, config *Config) {

	for _, p := range filePaths {
		processFile(p, config)
	}
}

func processInParallel(filePaths []string, config *Config) {

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

	var processor processors.EnvPlaceholderProcessor
	processor = replaceEnvsProcessor
	if config.FailOnMissingVariables {
		processor = &processors.EnvVariableExistsDecorator{
			Processor: processor,
		}
	}
	err := process(filePath, &processors.LinesProcessor{
		ProcessedLines:       []string{},
		PlaceholderProcessor: replaceEnvsProcessor,
	})
	if err != nil && config.ErrorListener != nil {
		config.ErrorListener.ErrorFound(filePath, err)
	}
}

func process(filePath string, s *processors.LinesProcessor) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	fileLines := strings.Split(string(data), "\n")
	for number, line := range fileLines {
		s.ProcessLine(number, line)
	}
	return s.ProcessFinishedForPath(filePath)
}
