package processors

// Event used when an env placeholder is found for a given line
type EnvPlaceholderFoundEvent struct {

	// line runes where the placeholder was found
	Line []rune

	// start idx of the line of runes where the placeholder was found
	Start int

	// end idx of the line of runes where the placeholder was found
	End int

	// placeholder value or env variable key as a rune slice
	Placeholder []rune
}

// Used to signal the result of processing an env placeholder for a given line
type EnvPlaceholderResult struct {

	// new line runes if the placeholder is replaced with a env placehodler value
	NewLine []rune

	// returns the new idx where the algorithm should continue to search for placeholders
	NewIdx int
}

// Called when a processing is finished for a file
type FileProcessingFinishedEvent struct {

	// Contains the processed content as lines of string
	ProcessedContent []string

	// The filepath for which the processing had taken place
	FilePath         string
}
