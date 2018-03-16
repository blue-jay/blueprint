// Package file provides helpful filesystem functions.
package file

import (
	"fmt"
	"io/ioutil"
	"os"
)

// Exists will return true if the file or folder exists.
func Exists(f string) bool {
	if _, err := os.Stat(f); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Copy copies a file from one location to another. It will return an
// error if the file already exists in the destination.
func Copy(src, dst string) error {
	if Exists(dst) {
		return fmt.Errorf("File, %v, already exists.", dst)
	}

	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dst, data, 0644)
}
