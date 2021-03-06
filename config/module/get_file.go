package module

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

// FileGetter is a Getter implementation that will download a module from
// a file scheme.
type FileGetter struct{}

func (g *FileGetter) Get(dst string, u *url.URL) error {
	// The source path must exist and be a directory to be usable.
	if fi, err := os.Stat(u.Path); err != nil {
		return fmt.Errorf("source path error: %s", err)
	} else if !fi.IsDir() {
		return fmt.Errorf("source path must be a directory")
	}

	fi, err := os.Lstat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// If the destination already exists, it must be a symlink
	if err == nil {
		mode := fi.Mode()
		if mode&os.ModeSymlink == 0 {
			return fmt.Errorf("destination exists and is not a symlink")
		}

		// Remove the destination
		if err := os.Remove(dst); err != nil {
			return err
		}
	}

	// Create all the parent directories
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	return os.Symlink(u.Path, dst)
}
