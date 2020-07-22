package tester

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func Test_defaultLicenseTester_HasSetupLicense(t *testing.T) {
	t.Run("should find old license in file", func(t *testing.T) {
		// given: a config file with old license
		configFile, err := ioutil.TempFile("", "confluence.*.cfg.xml")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(configFile.Name())
		oldContent := buildConfigFileContent(t, getSetupLicense())
		_, _ = configFile.WriteString(oldContent)
		_ = configFile.Sync()

		// when
		sut := New()
		actualChanged, err := sut.HasLicenseChanged(configFile.Name(), getSetupLicense())

		// then
		require.NoError(t, err)
		assert.False(t, actualChanged)
	})
	t.Run("should find new license in file", func(t *testing.T) {
		// given: a config file with new license
		configFile, err := ioutil.TempFile("", "confluence.*.cfg.xml")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(configFile.Name())

		newContent := buildConfigFileContent(t, getProductionLicense())
		_, _ = configFile.WriteString(newContent)
		_ = configFile.Sync()

		// when
		sut := New()
		actualChanged, err := sut.HasLicenseChanged(configFile.Name(), getSetupLicense())

		// then
		require.NoError(t, err)
		assert.True(t, actualChanged)
	})
	t.Run("should return error on weird config file", func(t *testing.T) {
		// given: a config file with new license
		configFile, err := ioutil.TempFile("", "confluence.*.cfg.xml")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(configFile.Name())

		malformedContent := "\n\n\n"
		_, _ = configFile.WriteString(malformedContent)
		_ = configFile.Sync()

		// when
		sut := New()
		_, actualErr := sut.HasLicenseChanged(configFile.Name(), getSetupLicense())

		// then
		require.Error(t, actualErr)
		assert.Contains(t, actualErr.Error(), "failed to find property")
	})
}

func Test_defaultLicenseTester_HasSetupLicense1(t *testing.T) {
	t.Run("should return true", func(t *testing.T) {
		// given: a config file with old license
		configFile, err := ioutil.TempFile("", "confluence.*.cfg.xml")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(configFile.Name())
		oldContent := buildConfigFileContent(t, getSetupLicense())
		_, _ = configFile.WriteString(oldContent)
		_ = configFile.Sync()

		// when
		sut := New()
		actualChanged, err := sut.HasSetupLicense(configFile.Name(), getSetupLicense())

		// then
		require.NoError(t, err)
		assert.True(t, actualChanged)
	})
	t.Run("should return false", func(t *testing.T) {
		// given: a config file with old license
		configFile, err := ioutil.TempFile("", "confluence.*.cfg.xml")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(configFile.Name())
		oldContent := buildConfigFileContent(t, getProductionLicense())
		_, _ = configFile.WriteString(oldContent)
		_ = configFile.Sync()

		// when
		sut := New()
		actualChanged, err := sut.HasSetupLicense(configFile.Name(), getSetupLicense())

		// then
		require.NoError(t, err)
		assert.False(t, actualChanged)
	})
	t.Run("should return error", func(t *testing.T) {
		// given: a config file with new license
		configFile, err := ioutil.TempFile("", "confluence.*.cfg.xml")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(configFile.Name())

		malformedContent := "\n\n\n"
		_, _ = configFile.WriteString(malformedContent)
		_ = configFile.Sync()

		// when
		sut := New()
		_, actualErr := sut.HasSetupLicense(configFile.Name(), getSetupLicense())

		// then
		require.Error(t, actualErr)
		assert.Contains(t, actualErr.Error(), "failed to find property")
	})
}

func Test_readLicenseLineFrom(t *testing.T) {
	t.Run("should return error on opening file", func(t *testing.T) {
		mockedOpener := new(fileOpenerMock)
		mockedOpener.On("Open", "some/file").Return(nil, assert.AnError)
		// when
		_, err := readLicenseLineFrom("some/file", mockedOpener)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error while opening")
		mockedOpener.AssertExpectations(t)
	})
}

// test util stuff

func getSetupLicense() string {
	return `AAABOA0ODAoPeNp9UVtPwjAUfu+SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP+X++SETUP++SETUP++SETUP+ +SETUP++SETUP++SETUP++SETUP++SETUP+X+++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP+XXXXX +SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP+XX/+SETUP++SETUP++SETUP++SETUP++SETUP+XXXX/XX +SETUP+XXXXX++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP++SETUP+XX BNy0Y6X9vBnwzAW2/+QZTOctX88X02fj`
}

func getProductionLicense() string {
	return `AAABOAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA AAAAAAA+AAAAAAAAAAAAAAAAAAAAAAAAAA/AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA AAAAAAAAAAAAAAAAAAAAAAAAAAAA//AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA AAAAAAAAAAAAAAAAAAA+AAAAAAAAAAAAAAAAAAAA++AAAAAAAAAAAAAAAAAAAAAAAAAAAAAA+AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA BNy0Y6X9vBnwzAW2/+QZTOctX88X02fj`
}

func buildConfigFileContent(t *testing.T, license string) string {
	t.Helper()

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>

<confluence-configuration>
  <setupStep>complete</setupStep>
  <setupType>custom</setupType>
  <buildNumber>8501</buildNumber>
  <properties>
    <property name="admin.ui.allow.daily.backup.custom.location">false</property>
    <property name="atlassian.license.message">%s</property>
    <property name="attachments.dir">${confluenceHome}/attachments</property>
  </properties>
</confluence-configuration>`,
		license)
}

type fileOpenerMock struct {
	mock.Mock
}

func (f *fileOpenerMock) Open(filePath string) (io.ReadCloser, error) {
	args := f.Called(filePath)
	file := args.Get(0)
	if file == nil {
		return nil, args.Error(1)
	}
	return file.(io.ReadCloser), args.Error(1)
}
