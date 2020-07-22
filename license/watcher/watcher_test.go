package watcher

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_defaultWatcher_doWatchWork(t *testing.T) {
	const licFile = "/var/atlassian/confluence/confluence.cfg.xml"
	const license = "AAAB/testLicense+=okBf"
	commandArgs := []string{"/opt/atlassian/confluence/bin/shutdown.sh"}

	t.Run("should call command executor on license change", func(t *testing.T) {
		// given
		args := &ProcessArgs{
			CommandArgs:          commandArgs,
			WatchIntervalInSecs:  30,
			ConfluenceConfigFile: licFile,
			SetupLicense:         license,
		}
		const licenseHasChanged = true
		mockedLicenseChecker := new(licenseTesterMock)
		mockedLicenseChecker.On("HasLicenseChanged", licFile, license).Return(licenseHasChanged, nil)

		mockedExecutor := new(executorMock)
		mockedExecutor.On("execute", commandArgs).Return("", nil)

		sut := defaultWatcher{
			args:          args,
			cmdExecutor:   mockedExecutor,
			licenseTester: mockedLicenseChecker,
		}

		// when
		finishWatcher, err := sut.doWatchWork()

		// then
		require.NoError(t, err)
		assert.True(t, finishWatcher)
		mockedLicenseChecker.AssertExpectations(t)
		mockedExecutor.AssertExpectations(t)
	})
	t.Run("should do nothing when license is still the same", func(t *testing.T) {
		// given
		args := &ProcessArgs{
			CommandArgs:          commandArgs,
			WatchIntervalInSecs:  30,
			ConfluenceConfigFile: licFile,
			SetupLicense:         license,
		}
		const licenseHasNotChanged = false
		mockedLicenseChecker := new(licenseTesterMock)
		mockedLicenseChecker.On("HasLicenseChanged", licFile, license).Return(licenseHasNotChanged, nil)

		mockedExecutor := new(executorMock)
		// no cmdExecutor modelling -> cmdExecutor will not be called

		sut := defaultWatcher{
			args:          args,
			cmdExecutor:   mockedExecutor,
			licenseTester: mockedLicenseChecker,
		}

		// when
		finishWatcher, err := sut.doWatchWork()

		// then
		require.NoError(t, err)
		assert.False(t, finishWatcher)
		mockedLicenseChecker.AssertExpectations(t)
		mockedExecutor.AssertExpectations(t)
	})
	t.Run("should fail on error that cannot be handled", func(t *testing.T) {
		// given
		args := &ProcessArgs{
			CommandArgs:          commandArgs,
			WatchIntervalInSecs:  30,
			ConfluenceConfigFile: licFile,
			SetupLicense:         license,
		}
		anError := assert.AnError
		mockedLicenseChecker := new(licenseTesterMock)
		mockedLicenseChecker.On("HasLicenseChanged", licFile, license).Return(false, anError)
		mockedExecutor := new(executorMock)

		sut := defaultWatcher{
			args:          args,
			cmdExecutor:   mockedExecutor,
			licenseTester: mockedLicenseChecker,
		}

		// when
		finishWatcher, err := sut.doWatchWork()

		// then
		require.Error(t, err)
		assert.True(t, finishWatcher)
		mockedLicenseChecker.AssertExpectations(t)
		mockedExecutor.AssertExpectations(t)
	})
}

// test util stuff
type executorMock struct {
	mock.Mock
}

func (e *executorMock) execute(shellCommandArgs []string) (string, error) {
	args := e.Called(shellCommandArgs)
	return args.String(0), args.Error(1)
}

type licenseTesterMock struct {
	mock.Mock
}

func (l *licenseTesterMock) HasLicenseChanged(configFile string, knownLicense string) (changed bool, err error) {
	args := l.Called(configFile, knownLicense)
	return args.Bool(0), args.Error(1)
}

func (l *licenseTesterMock) HasSetupLicense(configFile string, setupLicense string) (unchanged bool, err error) {
	args := l.Called(configFile, setupLicense)
	return args.Bool(0), args.Error(1)
}

func Test_defaultWatcher_Watch(t *testing.T) {
	t.Run("should create instance from defaultWatcher", func(t *testing.T) {
		args := &ProcessArgs{
			CommandArgs:          []string{},
			WatchIntervalInSecs:  30,
			ConfluenceConfigFile: "licFile",
			SetupLicense:         "license",
		}

		// when
		sut := New(args)

		// then
		require.IsType(t, &defaultWatcher{}, sut)
	})
}
