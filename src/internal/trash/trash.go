package trash

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type Backend string

const (
	BackendFreeDesktop Backend = "freedesktop"
	BackendMacOS       Backend = "macos"
	BackendWindows     Backend = "windows"
)

var ErrUnsupported = errors.New("trash is not supported on this platform")

type Result struct {
	OriginalPath string
	TrashedPath  string
	// InfoPath is the freedesktop .trashinfo file for TrashedPath, if any
	InfoPath         string
	Backend          Backend
	StrictlyRecycled bool
}

// Restore moves a trashed item back to its original path. It needs the
// TrashedPath that Move returned, so it does not work on Windows, where the
// recycle bin location is not known. It never overwrites: if something already
// exists at the original path, the item stays in the trash.
func Restore(r Result) error {
	if r.TrashedPath == "" {
		return ErrUnsupported
	}
	if _, err := os.Lstat(r.OriginalPath); err == nil {
		return fmt.Errorf("%s already exists", r.OriginalPath)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	if err := restoreMove(r.TrashedPath, r.OriginalPath); err != nil {
		return err
	}
	if r.InfoPath != "" {
		if err := os.Remove(r.InfoPath); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	return nil
}
