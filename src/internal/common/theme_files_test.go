package common

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
