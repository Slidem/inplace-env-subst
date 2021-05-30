package tests

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

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

func verifyExpectedContent(t *testing.T, file *os.File, expectedContent string) {
	actualContent := readContentFromFile(file.Name())
	if expectedContent != actualContent {
		t.Fatalf("Replaced content does not match expected one. \n#### Expected: \n%s\n#### Got: \n%s\n", expectedContent, actualContent)
	}
}

func verifyExpectedErrorMessage(t *testing.T, expectedErrMessage string, mc *MessageCaptureErrorListener) {
	if mc.message != expectedErrMessage {
		t.Fatalf("Actual error message different from expected one. \n#### Expected: \n%s\n#### Got: \n%s\n", expectedErrMessage, mc.message)
	}
}

type MessageCaptureErrorListener struct {
	message string
}

func (m *MessageCaptureErrorListener) ErrorFound(filepath string, err error) {
	m.message = err.Error()
}
