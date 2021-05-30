package tests

import (
	"github.com/Slidem/inplaceenvsubst/v2"
	"os"
	"testing"
)

func TestInvalidPlaceholderFormatIgnored(t *testing.T) {


	t.Run("Test replacement of env variable with invalid placeholder format ignored", func(t *testing.T) {
		// given
		testFileTemplate := ` env ${VAR}}}} has multiple closing brackets, and should be ignored`
		testFile := createTempFileFromTemplate(testFileTemplate)
		_ = os.Setenv("VAR", "testValue")

		defer func() {
			_ = os.Remove(testFile.Name())
			_ = os.Unsetenv("VAR")
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
		verifyExpectedContent(t, testFile, testFileTemplate)
	})
}

func TestSubstWithoutValidation(t *testing.T) {

	config := &inplaceenvsubst.Config{
		FailOnMissingVariables: false,
		RunInParallel:          false,
		ErrorListener:          nil,
	}

	var tempFile *os.File

	setup := func(testFileTemplate string) {
		tempFile = createTempFileFromTemplate(testFileTemplate)
		_ = os.Setenv("testA", "a")
		_ = os.Setenv("testB", "b")
	}

	teardown := func() {
		_ = os.Remove(tempFile.Name())
		_ = os.Unsetenv("testA")
		_ = os.Unsetenv("testB")
	}

	inputs := [][]string{
		{
			`
			testing ${testA}
			and ${testB} value
			-->${testA}${testB}<--
			${testA} at the beginning of line
			at the end of line ${testA}
			`,
			`
			testing a
			and b value
			-->ab<--
			a at the beginning of line
			at the end of line a
			`,
		},
		{
			`
			testing ${testA}
			and $$$${testB} value
			-->${testA}${testB}<--
			${testA} at the beginning of line
			at the end of line ${testA}
			`,
			`
			testing a
			and $$$b value
			-->ab<--
			a at the beginning of line
			at the end of line a
			`,
		},
		{
			`
			testing ${testA}
			and ${testB} value
			-->${testA}${testB}}}}<--
			${testA} at the beginning of line
			at the end of line ${testA}
			`,
			`
			testing a
			and b value
			-->a${testB}}}}<--
			a at the beginning of line
			at the end of line a
			`,
		},
	}

	for _, i := range inputs {
		t.Run("Test without validation", func(t *testing.T) {

			// given
			expectedContent := i[1]
			testFileTemplate := i[0]
			setup(testFileTemplate)
			defer teardown()

			// when
			inplaceenvsubst.ProcessFiles([]string{tempFile.Name()}, config)

			// expect
			verifyExpectedContent(t, tempFile, expectedContent)
		})
	}
}
