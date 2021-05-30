package inplaceenvsubst

import (
	"bufio"
	"errors"
	"github.com/Slidem/inplaceenvsubst/processors"
	"log"
	"os"
	"sync"
)

var envPlaceholderProcessor = &processors.ReplaceEnvVariablesProcessor{}

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

	if !config.WhitelistEnvs.IsEmpty() && !config.BlacklistEnvs.IsEmpty() {
		config.ErrorListener.ErrorFound(filePath, errors.New("cannot have both whitelist and blacklist environment keys"))
		return
	}

	var processor processors.EnvPlaceholderProcessor
	processor = envPlaceholderProcessor
	if config.FailOnMissingVariables {
		processor = &processors.EnvVariableExistsDecorator{
			Processor: processor,
		}
	}

	err := processLineByLine(filePath, &processors.LinesProcessor{
		ProcessedLines:       []string{},
		PlaceholderProcessor: processor,
		Config:               config,
	})

	if err != nil && config.ErrorListener != nil {
		config.ErrorListener.ErrorFound(filePath, err)
	}
}

func processLineByLine(filePath string, s *processors.LinesProcessor) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scn := bufio.NewScanner(file)

	var line string
	var lineNumber = 1
	for scn.Scan() {
		line = scn.Text()
		err = s.ProcessLine(lineNumber, line)
		if err != nil {
			return err
		}
		lineNumber += 1
	}

	if err = scn.Err(); err != nil {
		return err
	}

	return s.ProcessFinishedForPath(filePath)
}
