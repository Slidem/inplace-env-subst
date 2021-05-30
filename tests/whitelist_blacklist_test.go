package tests

import (
	"github.com/Slidem/inplaceenvsubst/v2"
	"os"
	"testing"
)

func TestWhitelistAndBlacklist(t *testing.T) {

	mc := MessageCaptureErrorListener{}

	config := inplaceenvsubst.Config{
		FailOnMissingVariables: false,
		RunInParallel:          false,
		ErrorListener:          &mc,
	}

	teardown := func(filename string) {
		mc.message = ""
		config.FailOnMissingVariables = false
		config.WhitelistEnvs = nil
		config.BlacklistEnvs = nil
		_ = os.Unsetenv("this")
		_ = os.Unsetenv("that")
		_ = os.Unsetenv("replaced")
		_ = os.Remove(filename)

	}

	t.Run("Given whitelist and blacklist config, then error is thrown", func(t *testing.T) {
		// given
		defer teardown("")

		config.WhitelistEnvs = inplaceenvsubst.NewStringSet("this", "that")
		config.BlacklistEnvs = inplaceenvsubst.NewStringSet("this", "that")

		// when
		inplaceenvsubst.ProcessFiles([]string{"dummyPath"}, &config)

		// expect
		verifyExpectedErrorMessage(t, "cannot have both whitelist and blacklist environment keys", &mc)
	})

	t.Run("Given whitelist variables when subst then non whitelisted variables ignored", func(t *testing.T) {
		// given
		testFileTemplate := `
		${this} and ${that} are ignore but ${replaced} is replaced 
		`
		testFile := createTempFileFromTemplate(testFileTemplate)
		defer teardown(testFile.Name())

		config.WhitelistEnvs = inplaceenvsubst.NewStringSet("replaced")
		config.FailOnMissingVariables = true
		_ = os.Setenv("this", "thisValue")
		_ = os.Setenv("that", "thatValue")
		_ = os.Setenv("replaced", "replacedValue")

		expectedContent := `
		${this} and ${that} are ignore but replacedValue is replaced 
		`
		// when
		inplaceenvsubst.ProcessFiles([]string{testFile.Name()}, &config)

		// expect
		verifyExpectedContent(t, testFile, expectedContent)
	})

	t.Run("Given whitelist variables when FailOnMissingVariables true and whitelist variable not found, then error returned", func(t *testing.T) {

		// given
		config.FailOnMissingVariables = true

		testFileTemplate := `
		${this} and ${that} are ignore but ${replaced} is replaced 
		`
		testFile := createTempFileFromTemplate(testFileTemplate)
		defer teardown(testFile.Name())

		config.WhitelistEnvs = inplaceenvsubst.NewStringSet("replaced")

		// when
		inplaceenvsubst.ProcessFiles([]string{testFile.Name()}, &config)

		// expect
		verifyExpectedErrorMessage(t, "Missing environment variables: [replaced]", &mc)
	})

	t.Run("Given blacklist variables when subst then blacklist variables ignored", func(t *testing.T) {
		// given
		testFileTemplate := `
		${this} and ${that} are ignore but ${replaced} is replaced 
		`
		testFile := createTempFileFromTemplate(testFileTemplate)
		defer teardown(testFile.Name())

		config.BlacklistEnvs = inplaceenvsubst.NewStringSet("this", "that")

		_ = os.Setenv("this", "thisValue")
		_ = os.Setenv("that", "thatValue")
		_ = os.Setenv("replaced", "replacedValue")

		expectedContent := `
		${this} and ${that} are ignore but replacedValue is replaced 
		`
		// when
		inplaceenvsubst.ProcessFiles([]string{testFile.Name()}, &config)

		// expect
		verifyExpectedContent(t, testFile, expectedContent)
	})

	t.Run("Given blacklist variables when FailOnMissingVariables true and blacklist variable not found, then blacklist variable ignored", func(t *testing.T) {

		// given
		config.FailOnMissingVariables = true

		testFileTemplate := `${replaced} is blacklisted, thus it's ignored`
		testFile := createTempFileFromTemplate(testFileTemplate)
		defer teardown(testFile.Name())

		_ = os.Setenv("replaced", "replacedValue")

		config.BlacklistEnvs = inplaceenvsubst.NewStringSet("replaced")

		// when
		inplaceenvsubst.ProcessFiles([]string{testFile.Name()}, &config)

		// expect
		verifyExpectedContent(t, testFile, testFileTemplate)
	})

	t.Run("Given blacklist variables, when FailOnMissingVariables true and non blacklisted variable not found, then error returned", func(t *testing.T) {

		// given
		config.FailOnMissingVariables = true

		testFileTemplate := `${replaced} is blacklisted, thus it's ignored, but ${this} should be replaced`
		testFile := createTempFileFromTemplate(testFileTemplate)
		defer teardown(testFile.Name())

		config.BlacklistEnvs = inplaceenvsubst.NewStringSet("replaced")

		// when
		inplaceenvsubst.ProcessFiles([]string{testFile.Name()}, &config)

		// expect
		verifyExpectedErrorMessage(t, "Missing environment variables: [this]", &mc)
	})
}
