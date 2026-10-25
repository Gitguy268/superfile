package filepanel

import (
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// TODO : Take common.Config.CaseSensitiveSort as a function parameter
// and also consider testing this caseSensitive with both true and false in
// our unit_test TestReturnDirElement
// getDirectoryElements returns the directory elements for the panel's current location
func (m *Model) getDirectoryElements(displayDotFile bool) []Element {
	dirEntries, err := os.ReadDir(m.Location)
	if err != nil {
		slog.Error("Error while returning folder elements", "error", err)
		return nil
	}

	dirEntries = slices.DeleteFunc(dirEntries, func(e os.DirEntry) bool {
		return skipEntry(m.Location, e, displayDotFile, variable.DirOnly)
	})

	// No files/directories to process
	if len(dirEntries) == 0 {
		return nil
	}
	return sortFileElement(m.SortKind, m.SortReversed, dirEntries, m.Location)
}

// getDirectoryElementsBySearch returns filtered directory elements based on search string
func (m *Model) getDirectoryElementsBySearch(displayDotFile bool) []Element {
	searchString := m.SearchBar.Value()
	items, err := os.ReadDir(m.Location)
	if err != nil {
		slog.Error("Error while return folder element function", "error", err)
		return nil
	}

	folderElementMap := map[string]os.DirEntry{}
	fileAndDirectories := []string{}

	for _, item := range items {
		if skipEntry(m.Location, item, displayDotFile, variable.DirOnly) {
			continue
		}

		fileAndDirectories = append(fileAndDirectories, item.Name())
		folderElementMap[item.Name()] = item
	}
	if len(fileAndDirectories) == 0 {
		return nil
	}
	// https://github.com/reinhrst/fzf-lib/blob/main/core.go#L43
	// fzf returns matches ordered by score; we subsequently sort by the chosen sort option.
	fzfResults := utils.FzfSearch(searchString, fileAndDirectories)
	dirElements := make([]os.DirEntry, 0, len(fzfResults))
	for _, item := range fzfResults {
		resultItem := folderElementMap[item.Key]
		dirElements = append(dirElements, resultItem)
	}

	return sortFileElement(m.SortKind, m.SortReversed, dirElements, m.Location)
}

// skipEntry reports whether a directory entry should be left out of the panel:
// unreadable entries, dotfiles when they are hidden, and anything that is not a
// directory (or a symlink to one) when dirOnly is set.
func skipEntry(location string, e os.DirEntry, displayDotFile bool, dirOnly bool) bool {
	if _, err := e.Info(); err != nil {
		return true
	}
	if !displayDotFile && strings.HasPrefix(e.Name(), ".") {
		return true
	}
	if !dirOnly || e.IsDir() {
		return false
	}
	if e.Type()&os.ModeSymlink == 0 {
		return true
	}
	info, err := os.Stat(filepath.Join(location, e.Name()))
	return err != nil || !info.IsDir()
}

// Helper to decide whether to skip updating a panel this tick.
func (m *Model) shouldSkipPanelUpdate(nowTime time.Time) bool {
	if !m.IsFocused {
		return nowTime.Sub(m.LastTimeGetElement) < nonFocussedPanelReRenderTime
	}

	reRenderTime := int(float64(m.ElemCount()) / ReRenderChunkDivisor)
	reRenderTime = min(reRenderTime, ReRenderMaxDelay)
	return !m.NeedsReRender() &&
		nowTime.Sub(m.LastTimeGetElement) < time.Duration(reRenderTime)*time.Second
}

func (m *Model) UpdateElementsIfNeeded(force bool, displayDotFile bool) {
	nowTime := time.Now()
	if force || !m.shouldSkipPanelUpdate(nowTime) {
		// Load elements for this panel (with/without search filter)
		m.element = m.getElements(displayDotFile)
		// Update file panel list
		m.LastTimeGetElement = nowTime
		m.lastFetch = m.currentFetchInputs()

		// For hover to file on first time loading
		if m.TargetFile != "" {
			m.applyTargetFileCursor()
		}

		// If cursor becomes invalid due to element update, reset
		if m.ValidateCursorAndRenderIndex() != nil {
			m.scrollToCursor(0)
		}
	}
}

// Retrieves elements for a panel based on search bar value and sort options.
func (m *Model) getElements(displayDotFile bool) []Element {
	if m.SearchBar.Value() != "" {
		return m.getDirectoryElementsBySearch(displayDotFile)
	}
	return m.getDirectoryElements(displayDotFile)
}
