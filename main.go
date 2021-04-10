package main

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

var logger = GetLogger()

func main() {

	filePaths := os.Args[1:]
	for _, fp := range filePaths {
		processFile(fp, &LineProcessor{[]string{}, &ReplaceEnvVariables{}})
	}
}

func processFile(fp string, s *LineProcessor) {
	input, err := ioutil.ReadFile(fp)
	if err != nil {
		panic(err)
	}
	fileLines := strings.Split(string(input), "\n")
	for number, line := range fileLines {
		s.processLine(number, line)
	}
	s.processFinishedForPath(fp)
}

type EnvProcessor interface {

	envFound(int, int, string, []rune) (error, []rune)
	lineProcessingFinished([]string, string) error
}

type ReplaceEnvVariables struct {

}

func (e *ReplaceEnvVariables) envFound(start, end int, env string, lineRunes []rune) (error, []rune) {
	value, exists := os.LookupEnv(env)
	if !exists {
		return errors.New("Env variable " + env + " does not exist"), nil
	}
	newRunes := append(lineRunes[:start], []rune(value)...)
	newRunes = append(newRunes, lineRunes[end:]...)
	return nil, newRunes
}

func (e *ReplaceEnvVariables) lineProcessingFinished(lines []string, filePath string) error {
	output := strings.Join(lines, "\n")
	return ioutil.WriteFile(filePath, []byte(output), os.ModeAppend)
}

type LineProcessor struct {
	processedLines []string
	envProcessor EnvProcessor
}

func (p *LineProcessor) processLine(number int, line string) {

	runes := []rune(line)
	idx := 0
	var err error

	for idx < len(runes) {
		for idx < len(runes) && runes[idx] != '$' {
			idx++
		}
		if idx == len(runes) {
			break
		}
		if runes[idx] == '$' {
			for runes[idx] == '$' {
				idx++
			}
			if idx == len(runes) {
				break
			}
			if runes[idx] == '{' {
				start := idx - 1
				var envSb strings.Builder
				for {
					idx++
					if idx == len(runes) {
						logger.Fatalf("Expected } but found newline on line number %d", number + 1)
						return
					}
					if runes[idx] == '}' {
						idx++
						for idx < len(runes) && runes[idx] == '}' {
							envSb.WriteRune('}')
							idx++
						}
						end := idx
						env := envSb.String()
						if env != "" {
							err, runes = p.envProcessor.envFound(start, end, env, runes)
							if err != nil {
								logger.Fatal(err)
							}
							break
						}
					}
					envSb.WriteRune(runes[idx])
				}
			}
		}
		idx++
	}
	p.processedLines = append(p.processedLines, string(runes))
}

func (p *LineProcessor) processFinishedForPath(filePath string) {
	p.envProcessor.lineProcessingFinished(p.processedLines, filePath)
}
