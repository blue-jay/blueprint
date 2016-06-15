package view

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PrependBaseURI prepends the base URI to the string
func (v *View) PrependBaseURI(s string) string {
	return v.BaseURI + s
}

// AssetTimePath returns a URL with the proper base uri and timestamp appended.
// Works for CSS and JS assets
// Determines if local or on the web
func (v *View) AssetTimePath(s string) (string, error) {
	if strings.HasPrefix(s, "//") {
		return s, nil
	}

	s = strings.TrimLeft(s, "/")
	abs, err := filepath.Abs(s)

	if err != nil {
		return "", err
	}

	time, err2 := FileTime(abs)
	if err2 != nil {
		return "", err2
	}

	return v.PrependBaseURI(s + "?" + time), nil
}

// FileTime returns the modification time of the file
func FileTime(name string) (string, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return "", err
	}
	mtime := fi.ModTime().Unix()
	return fmt.Sprintf("%v", mtime), nil
}
