// Package replace will search for matched case-sensitive strings in files
// and then replace them with a different string.
//
// Examples:
//	jay replace . red blue
//		Replace the word "red" with the word "blue" in all go files in current folder and in subfolders.
//	jay replace . red blue "*.go" true true
//		Replace the word "red" with the word "blue" in *.go files in current folder including filenames and in subfolders.
//	jay replace . "blue-jay/blueprint" "user/project"
//		Change the name of the project in current folder and in subfolders and all imports to another repository.
package replace

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
	flagReplace   *string
	flagExt       *string
	flagName      *bool
	flagRecursive *bool
	flagCommit    *bool

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

// Run starts the replace filepath walk.
func Run(find, folder, replace, ext *string, recursive, filename, commit *bool) ([]string, error) {
	flagFind = find
	flagFolder = folder
	flagReplace = replace
	flagExt = ext
	flagRecursive = recursive
	flagName = filename
	flagCommit = commit

	contents = []string{}

	if *flagCommit {
		record("Replace Results")
		record("===============")
	} else {
		record("Replace Results (no changes)")
		record("============================")
	}

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

		// If the file name contains the search term, replace the file name
		if *flagName && strings.Contains(fi.Name(), *flagFind) {
			//TODO Fix the bug where if the folder AND file name match, it won't be changed
			// Only change the filename, not the folder, or rename?
			oldpath := path
			path = strings.Replace(path, *flagFind, *flagReplace, -1)
			record(" Rename:", oldpath, "("+path+")")

			if *flagCommit {
				errRename := os.Rename(oldpath, path)
				if errRename != nil {
					record("**ERROR: Could not rename", oldpath, "to", path)
					return nil
				}
			}
		}

		return folderCheck(fi)
	}

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
			// Replace the search term
			newContents := strings.Replace(oldContents, *flagFind, *flagReplace, -1)
			count := strconv.Itoa(strings.Count(oldContents, *flagFind))
			record("Replace:", path, "("+count+")")

			// Write the data back to the file
			if *flagCommit {
				err = ioutil.WriteFile(path, []byte(newContents), 0)
				if err != nil {
					record("**ERROR: Could not write to", path)
					return nil
				}
			}
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
