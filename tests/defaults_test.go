package tests

import (
	"github.com/Slidem/inplaceenvsubst/v2"
	"os"
	"testing"
)

func TestDefaults(t *testing.T) {

	t.Run("Given file with defaults when env variables missing then defaults added", func(t *testing.T) {

		// given
		testFile := createTempFileFromTemplate("Testing default ${missing:-defaultValue}")
		defer func() {
			_ = os.Remove(testFile.Name())
		}()

		// when
		inplaceenvsubst.ProcessFiles([]string{testFile.Name()}, &inplaceenvsubst.Config{
			FailOnMissingVariables: false,
			RunInParallel:          false,
		})

		// expect
		verifyExpectedContent(t, testFile, "Testing default defaultValue")
	})

	t.Run("Given file with defaults when env variables present then env value added", func(t *testing.T) {

		// given
		testFile := createTempFileFromTemplate("Testing default ${env:-defaultValue}")
		_ = os.Setenv("env", "envValue")
		defer func() {
			_ = os.Remove(testFile.Name())
			_ = os.Unsetenv("env")
		}()

		// when
		inplaceenvsubst.ProcessFiles([]string{testFile.Name()}, &inplaceenvsubst.Config{
			FailOnMissingVariables: false,
			RunInParallel:          false,
		})

		// expect
		verifyExpectedContent(t, testFile, "Testing default envValue")
	})

	t.Run("Given file with defaults when env variable missing and FailOnMissingVariables true then defaults added", func(t *testing.T) {
		// given
		testFile := createTempFileFromTemplate("Testing default ${missing:-defaultValue}")
		defer func() {
			_ = os.Remove(testFile.Name())
		}()

		// when
		inplaceenvsubst.ProcessFiles([]string{testFile.Name()}, &inplaceenvsubst.Config{
			FailOnMissingVariables: true,
			RunInParallel:          false,
		})

		// expect
		verifyExpectedContent(t, testFile, "Testing default defaultValue")
	})

	t.Run("Given file with default value empty when env variable missing and FailOnMissingVariables true then defaults added", func(t *testing.T) {
		// given
		testFile := createTempFileFromTemplate("Testing default ${missing:-}")
		defer func() {
			_ = os.Remove(testFile.Name())
		}()

		// when
		inplaceenvsubst.ProcessFiles([]string{testFile.Name()}, &inplaceenvsubst.Config{
			FailOnMissingVariables: true,
			RunInParallel:          false,
		})

		// expect
		verifyExpectedContent(t, testFile, "Testing default ")
	})
}