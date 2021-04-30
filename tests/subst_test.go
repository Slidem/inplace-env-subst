package tests

import (
	"github.com/Slidem/inplaceenvsubst"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestSubstPanicWhenPlaceholderNotClosed(t *testing.T) {

	testFileTemplate := `
		opened but not closed ${NOT_CLOSING and other text
	`
	testFile := createTempFileFromTemplate(testFileTemplate)
	defer os.Remove(testFile.Name())

	// expect to panic
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Expected panic")
		}
	}()

	// when
	inplaceenvsubst.ProcessFiles([]string{testFile.Name()}, &inplaceenvsubst.Config{
		FailOnMissingVariables: false,
		RunInParallel:          false,
		ErrorListener:          nil,
	})
}

func TestSubstWithValueBiggerThanKey(t *testing.T) {
	t.Run("Test without validation", func(t *testing.T) {
		// given
		os.Setenv("VAR", "bigger_value")
		expectedContent := ` env bigger_value value length bigger than placeholder length`
		testFileTemplate := ` env ${VAR} value length bigger than placeholder length`
		testFile := createTempFileFromTemplate(testFileTemplate)
		defer func() {
			os.Remove(testFile.Name())
			os.Unsetenv("VAR")
		}()
		// when
		inplaceenvsubst.ProcessFiles([]string{testFile.Name()}, &inplaceenvsubst.Config{
			FailOnMissingVariables: false,
			RunInParallel:          false,
			ErrorListener:          nil,
		})
		// expect
		replacedContent := readContentFromFile(testFile.Name())
		if expectedContent != replacedContent {
			t.Fatalf("Replaced content does not match expected one")
		}
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
		os.Remove(tempFile.Name())
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
			replacedContent := readContentFromFile(tempFile.Name())
			if expectedContent != replacedContent {
				t.Fatalf("Replaced content does not match expected one")
			}
		})
	}
}

func readContentFromFile(name string) string {
	content, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func createTempFileFromTemplate(testFileTemplate string) *os.File {
	tempFile, err := ioutil.TempFile(os.TempDir(), "test-template-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	if _, err = tempFile.Write([]byte(testFileTemplate)); err != nil {
		log.Fatal("Failed to write to test file", err)
	}
	defer tempFile.Close()
	return tempFile
}
