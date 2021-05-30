package tests

import (
	"github.com/Slidem/inplaceenvsubst"
	"os"
	"testing"
)

func TestSubstWithValueBiggerThanKey(t *testing.T) {
	t.Run("Test replacement of env variable value larger than key", func(t *testing.T) {
		// given
		_ = os.Setenv("VAR", "bigger_value")
		expectedContent := ` env bigger_value value length bigger than placeholder length`
		testFileTemplate := ` env ${VAR} value length bigger than placeholder length`
		testFile := createTempFileFromTemplate(testFileTemplate)
		defer func() {
			_ = os.Remove(testFile.Name())
			_ = os.Unsetenv("VAR")
		}()
		// when
		inplaceenvsubst.ProcessFiles([]string{testFile.Name()}, &inplaceenvsubst.Config{
			FailOnMissingVariables: false,
			RunInParallel:          false,
			ErrorListener:          nil,
		})
		// expect
		verifyExpectedContent(t, testFile, expectedContent)
	})
}