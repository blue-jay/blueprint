// Package find will search for matched case-sensitive strings in files.
//
// Examples:
//	jay find . red
//		Find the word "red" in all go files in current folder and in subfolders.
//	jay find . red "*.*"
//		Find the word "red" in all files in current folder and in subfolders.
//	jay find . red "*.go" true false
//		Find word "red" in *.go files in current folder and in subfolders, but will exclude filenames.
package find

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	flagFind      *string
	flagFolder    *string
	flagExt       *string
	flagName      *bool
	flagRecursive *bool

	contents []string // Buffer to hold lines of output

	// MaxSize is the maximum size of a file Go will search through
	MaxSize int64 = DefaultMaxSize()
	// SkipFolders is folders that won't be searched
	SkipFolders = DefaultSkipFolders()
)

// DefaultSkipFolders returns the folders that are skipped by default.
func DefaultSkipFolders() []string {
	return []string{"vendor", "node_modules", ".git"}
}

// DefaultMaxSize returns the default max filesize.
func DefaultMaxSize() int64 {
	return 1048576
}

// record writes the line to the string array.
func record(line ...string) {
	contents = append(contents, strings.Join(line, " "))
}

// Run starts the find filepath walk, will return the results in a string
// array, and will reset the defaults for SkipFolders and MaxSize.
func Run(text, folder, ext *string, recursive, filename *bool) ([]string, error) {
	flagFind = text
	flagFolder = folder
	flagExt = ext
	flagRecursive = recursive
	flagName = filename

	contents = []string{}

	record("Search Results")
	record("==============")

	err := filepath.Walk(*folder, visit)

	// Reset the folders
	SkipFolders = DefaultSkipFolders()

	return contents, err
}

// Visit analyzes a file to see if it matches the parameters.
// Original: https://gist.github.com/tdegrunt/045f6b3377f3f7ffa408
func visit(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	// If path is a folder
	if fi.IsDir() {
		// Ignore specified folders
		if inArray(fi.Name(), SkipFolders) {
			return filepath.SkipDir
		}

		// If the folder name contains the search term, show the folder name
		if *flagName && strings.Contains(fi.Name(), *flagFind) {
			record("Filename:", path)
		}

		return folderCheck(fi)
	}

	// Determine if the extension matches the file
	matched, err := filepath.Match(*flagExt, fi.Name())
	if err != nil {
		return err
	}

	// If the file extension matches
	if matched {
		// Skip file if too big
		if fi.Size() > MaxSize {
			record("**ERROR: Skipping file too big", path)
			return nil
		}

		// Read the entire file into memory
		read, err := ioutil.ReadFile(path)
		if err != nil {
			record("**ERROR: Could not read from", path)
			return nil
		}

		// Convert the bytes array into a string
		oldContents := string(read)

		// If the file contains the search term
		if strings.Contains(oldContents, *flagFind) {
			count := strconv.Itoa(strings.Count(oldContents, *flagFind))
			record("Contents:", path, "("+count+")")

		}
	}

	return nil
}

func folderCheck(fi os.FileInfo) error {
	// Always search current folder
	if fi.Name() == "." {
		return nil
	}

	// If recursive is true
	if *flagRecursive || *flagFolder == fi.Name() {
		return nil
	}

	// Don't walk the folder
	return filepath.SkipDir
}

func inArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
