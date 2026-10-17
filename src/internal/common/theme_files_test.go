package common

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
	"testing/fstest"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	variable "github.com/yorukot/superfile/src/config"
)

// Every bundled theme must set the same keys as the reference theme and load into ThemeType
func TestBundledThemesAreComplete(t *testing.T) {
	themeKeys := func(t *testing.T, path string) []string {
		t.Helper()
		data, err := os.ReadFile(path)
		require.NoError(t, err)
		var raw map[string]any
		require.NoError(t, toml.Unmarshal(data, &raw))
		keys := make([]string, 0, len(raw))
		for k := range raw {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		return keys
	}

	themeDir := filepath.Join("..", "..", "superfile_config", "theme")
	want := themeKeys(t, filepath.Join(themeDir, "catppuccin-mocha.toml"))
	files, err := filepath.Glob(filepath.Join(themeDir, "*.toml"))
	require.NoError(t, err)
	require.NotEmpty(t, files)

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			assert.Equal(t, want, themeKeys(t, file))
			var theme ThemeType
			require.NoError(t, LoadUserTheme(file, &theme))
			assert.Len(t, theme.GradientColor, 2)
		})
	}
}

// Themes added in a new build must reach users whose theme version is already current
func TestWriteThemeFilesOnlyMissing(t *testing.T) {
	oldFolder := variable.ThemeFolder
	variable.ThemeFolder = t.TempDir()                     //nolint:reassign // Needed to tests
	t.Cleanup(func() { variable.ThemeFolder = oldFolder }) //nolint:reassign // Needed to tests

	edited := filepath.Join(variable.ThemeFolder, "old.toml")
	require.NoError(t, os.WriteFile(edited, []byte("user edit"), 0o644))
	content := fstest.MapFS{
		variable.EmbedThemeDir + "/old.toml": {Data: []byte("bundled old")},
		variable.EmbedThemeDir + "/new.toml": {Data: []byte("bundled new")},
	}

	require.NoError(t, WriteThemeFiles(content, true))
	data, err := os.ReadFile(edited)
	require.NoError(t, err)
	assert.Equal(t, "user edit", string(data), "existing themes are kept")
	data, err = os.ReadFile(filepath.Join(variable.ThemeFolder, "new.toml"))
	require.NoError(t, err)
	assert.Equal(t, "bundled new", string(data), "missing themes are added")

	require.NoError(t, WriteThemeFiles(content, false))
	data, err = os.ReadFile(edited)
	require.NoError(t, err)
	assert.Equal(t, "bundled old", string(data), "a version change rewrites all themes")
}
