package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMoveElementDoesNotReplaceConcurrentDestination(t *testing.T) {
	tempDir := t.TempDir()
	destination := filepath.Join(tempDir, "destination.txt")
	sources := []string{
		filepath.Join(tempDir, "first.txt"),
		filepath.Join(tempDir, "second.txt"),
	}
	contents := []string{"first", "second"}
	contentBySource := make(map[string]string, len(sources))
	for index, source := range sources {
		contentBySource[source] = contents[index]
		require.NoError(t, os.WriteFile(source, []byte(contents[index]), 0o644))
	}

	type result struct {
		source string
		err    error
	}
	start := make(chan struct{})
	results := make(chan result, len(sources))
	for _, source := range sources {
		go func() {
			<-start
			results <- result{source: source, err: moveElement(source, destination)}
		}()
	}
	close(start)

	var succeeded, rejected result
	for range sources {
		moveResult := <-results
		if moveResult.err == nil {
			succeeded = moveResult
		} else {
			rejected = moveResult
		}
	}

	require.NotEmpty(t, succeeded.source)
	require.NotEmpty(t, rejected.source)
	require.ErrorIs(t, rejected.err, os.ErrExist)
	_, err := os.Stat(succeeded.source)
	assert.True(t, os.IsNotExist(err))
	assert.FileExists(t, rejected.source)
	destinationContents, err := os.ReadFile(destination)
	require.NoError(t, err)
	assert.Equal(t, contentBySource[succeeded.source], string(destinationContents))
}

func TestRenameIfAbsent(t *testing.T) {
	dir := t.TempDir()
	src, dst := filepath.Join(dir, "src.txt"), filepath.Join(dir, "dst.txt")
	require.NoError(t, os.WriteFile(src, []byte("src"), 0o644))
	require.NoError(t, os.WriteFile(dst, []byte("dst"), 0o644))

	require.ErrorIs(t, renameIfAbsent(src, dst), os.ErrExist)
	data, err := os.ReadFile(dst)
	require.NoError(t, err)
	assert.Equal(t, "dst", string(data), "existing destination must be kept")

	free := filepath.Join(dir, "free.txt")
	require.NoError(t, renameIfAbsent(src, free))
	assert.FileExists(t, free)
	assert.NoFileExists(t, src)
}
