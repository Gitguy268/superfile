//go:build !darwin && !linux && !windows

package internal

func renameNoReplace(src, dst string) error {
	return renameIfAbsent(src, dst)
}
