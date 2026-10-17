//go:build darwin

package internal

import (
	"errors"

	"golang.org/x/sys/unix"
)

func renameNoReplace(src, dst string) error {
	err := unix.RenamexNp(src, dst, unix.RENAME_EXCL)
	// Some filesystems (e.g. SMB, FAT) don't support RENAME_EXCL
	if errors.Is(err, unix.ENOTSUP) || errors.Is(err, unix.EINVAL) {
		return renameIfAbsent(src, dst)
	}
	return err
}
