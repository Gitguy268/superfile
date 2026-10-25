package trash

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRestore(t *testing.T) {
	setup := func(t *testing.T) Result {
		dir := t.TempDir()
		r := Result{
			OriginalPath: filepath.Join(dir, "orig.txt"),
			TrashedPath:  filepath.Join(dir, "trashed.txt"),
			InfoPath:     filepath.Join(dir, "trashed.txt.trashinfo"),
		}
		require.NoError(t, os.WriteFile(r.TrashedPath, []byte("content"), 0o644))
		require.NoError(t, os.WriteFile(r.InfoPath, []byte("[Trash Info]"), 0o644))
		return r
	}

	t.Run("moves the item back and removes its trashinfo", func(t *testing.T) {
		r := setup(t)
		require.NoError(t, Restore(r))
		data, err := os.ReadFile(r.OriginalPath)
		require.NoError(t, err)
		assert.Equal(t, "content", string(data))
		assert.NoFileExists(t, r.TrashedPath)
		assert.NoFileExists(t, r.InfoPath)
	})

	t.Run("never overwrites the original path", func(t *testing.T) {
		r := setup(t)
		require.NoError(t, os.WriteFile(r.OriginalPath, []byte("new"), 0o644))
		require.Error(t, Restore(r))
		data, err := os.ReadFile(r.OriginalPath)
		require.NoError(t, err)
		assert.Equal(t, "new", string(data))
		assert.FileExists(t, r.TrashedPath)
		assert.FileExists(t, r.InfoPath)
	})

	t.Run("needs a trashed path", func(t *testing.T) {
		require.ErrorIs(t, Restore(Result{OriginalPath: "/nonexistent"}), ErrUnsupported)
	})
}
