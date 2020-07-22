package watcher

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_defaultExecutor_execute(t *testing.T) {
	t.Run("should echo to stdout", func(t *testing.T) {
		sut := &defaultExecutor{}

		actual, err := sut.execute([]string{"/bin/echo", "-n", "hello", "world"})

		require.NoError(t, err)
		assert.Equal(t, "hello world", actual)
	})
	t.Run("should fail with output from stderr", func(t *testing.T) {
		sut := &defaultExecutor{}

		actual, err := sut.execute([]string{"/bin/something", "--not", "existing"})

		require.Error(t, err)
		assert.Contains(t, actual, "no such file or directory")
	})
}
