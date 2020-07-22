package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_checkMainError(t *testing.T) {
	t.Run("should not exit with a nil error", func(t *testing.T) {
		mockedExiter := new(mockExiter)

		checkMainError(nil, mockedExiter)

		mockedExiter.AssertExpectations(t)
	})
	t.Run("should exit 1 with any error and not print stacktrace", func(t *testing.T) {
		mockedExiter := new(mockExiter)
		mockedExiter.On("exit", 1)

		checkMainError(assert.AnError, mockedExiter)

		mockedExiter.AssertExpectations(t)
	})
	t.Run("should exit 1 with any error and print stacktrace", func(t *testing.T) {
		os.Args = []string{"--show-stack"}
		defer func() { os.Args = []string{} }()

		mockedExiter := new(mockExiter)
		mockedExiter.On("exit", 1)

		checkMainError(assert.AnError, mockedExiter)

		mockedExiter.AssertExpectations(t)
	})
}

func TestWatchCommand(t *testing.T) {
	actual := WatchCommand()

	require.NotNil(t, actual)
}

func TestTestLicenseCommand(t *testing.T) {
	actual := TestLicenseCommand()

	require.NotNil(t, actual)
}

func Test_createGlobalFlags(t *testing.T) {
	t.Run("should return three flags", func(t *testing.T) {
		actual := createGlobalFlags()

		require.Len(t, actual, 3)
	})
}

// test util stuff
type mockExiter struct {
	mock.Mock
}

func (m *mockExiter) exit(exitCode int) {
	m.Called(exitCode)
}

func Test_newExiter(t *testing.T) {
	t.Run("should create an exiter instance", func(t *testing.T) {
		sut := newExiter()

		require.NotNil(t, sut)
		assert.Implements(t, (*exiter)(nil), sut)
	})
}
