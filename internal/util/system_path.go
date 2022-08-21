// Package util contains utility functions.
package util

import (
	"os"
	"path"
	"runtime"
)

// ChangeDir changes current directory to root directory.
// It returns error if failed to change directory.
func ChangeDir(pathStr string) error {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), pathStr)
	err := os.Chdir(dir)
	if err != nil {
		return err
	}

	return nil
}
