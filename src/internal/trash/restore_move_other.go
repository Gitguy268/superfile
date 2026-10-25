//go:build !linux

package trash

import "os"

func restoreMove(src, dst string) error {
	return os.Rename(src, dst)
}
