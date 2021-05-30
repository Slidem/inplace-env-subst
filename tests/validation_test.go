package tests

import (
	"github.com/Slidem/inplaceenvsubst/v2"
	"os"
	"testing"
)

func TestWithValidation(t *testing.T) {

	t.Run("Test replacement of env variable value larger than key", func(t *testing.T) {
		// given
		testFileTemplate := ` env ${VAR} value length bigger than placeholder length`
		testFile := createTempFileFromTemplate(testFileTemplate)
		defer func() {
			_ = os.Remove(testFile.Name())
		}()

		mc := MessageCaptureErrorListener{}
		inplaceenvsubst.ProcessFiles([]string{testFile.Name()}, &inplaceenvsubst.Config{
			FailOnMissingVariables: true,
			RunInParallel:          false,
			ErrorListener:          &mc,
		})

		// when
		_ = readContentFromFile(testFile.Name())

		// expect
		verifyExpectedErrorMessage(t, "Missing environment variables: [VAR]", &mc)
	})
}