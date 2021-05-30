package processors

import (
	"errors"
	"fmt"
)

type LinesProcessor struct {
	ProcessedLines       []string
	PlaceholderProcessor EnvPlaceholderProcessor
	Config               Config
}

// TODO: improve this algorithm
// it can be cleaned-up and modified to support different type of placeholder delimiters
// right now it just supports ${} placeholder delimiter
// another options, is to have different tipes of Lines Processors, based on the configuration
// this would allow an easier way to maintain both $, ${}, and __variable__ types
func (p *LinesProcessor) ProcessLine(number int, line string) error {

	lineRunes := []rune(line)
	idx := 0

	for idx < len(lineRunes) {

		// increment until the start of placeholder
		for idx < len(lineRunes) && lineRunes[idx] != '$' {
			idx++
		}

		// if duplicate starts of placeholder, increment until we find the next non matching character
		for idx < len(lineRunes) && lineRunes[idx] == '$' {
			idx++
		}

		// make sure we break if end of line reached
		if idx == len(lineRunes) {
			break
		}

		// if second character is valid, this marks the start of our placeholder name, or the key of the env var
		if lineRunes[idx] == '{' {

			// the start is actually where the $ sign is found, so go back an index
			start := idx - 1

			var placeHolderRunes []rune
			for {
				idx++
				if idx == len(lineRunes) {
					// special case, where a placeholder was opened, but not closed at the end of file (eg: ${something \n)
					return errors.New(fmt.Sprintf("Expected } but found newline on line number %d", number))
				}

				// if closing character found
				if lineRunes[idx] == '}' {
					idx++
					// find the last closing character (we might get }}}, in which case we consider the last } the end)
					for idx < len(lineRunes) && lineRunes[idx] == '}' {

						// we consider the duplicates part of the placeholder (useful for debugging too !! :D)
						placeHolderRunes = append(placeHolderRunes, '}')
						idx++
					}
					end := idx
					// trigger placeholder found event
					if len(placeHolderRunes) != 0 && !p.Config.ShouldIgnoreEnv(string(placeHolderRunes)){
						response := p.PlaceholderProcessor.EnvPlaceholderFound(&EnvPlaceholderFoundEvent{
							Start:       start,
							End:         end,
							Placeholder: placeHolderRunes,
							Line:        lineRunes,
						})

						// update the line runes and the new index from which to continue (when line length changed , etc)
						lineRunes = response.NewLine
						idx = response.NewIdx
					}
					break
				}

				placeHolderRunes = append(placeHolderRunes, lineRunes[idx])
			}
		}
	}
	p.ProcessedLines = append(p.ProcessedLines, string(lineRunes))

	return nil
}

func (p *LinesProcessor) ProcessFinishedForPath(filePath string) error {
	return p.PlaceholderProcessor.FileProcessingFinished(
		&FileProcessingFinishedEvent{
			ProcessedContent: p.ProcessedLines,
			FilePath:         filePath,
		})
}
