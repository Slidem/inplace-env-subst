package tests

import (
	"github.com/Slidem/inplaceenvsubst/v2"
	"os"
	"testing"
)

func TestSubstErrorWhenPlaceholderNotClosed(t *testing.T) {

	testFileTemplate := `opened but not closed ${NOT_CLOSING and other text`
	testFile := createTempFileFromTemplate(testFileTemplate)
	defer func() {
		_ = os.Remove(testFile.Name())
	}()

	// when
	mc := MessageCaptureErrorListener{}
	inplaceenvsubst.ProcessFiles([]string{testFile.Name()}, &inplaceenvsubst.Config{
		FailOnMissingVariables: false,
		RunInParallel:          false,
		ErrorListener:          &mc,
	})

	// expect
	verifyExpectedErrorMessage(t, "Expected } but found newline on line number 1", &mc)
}

