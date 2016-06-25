package asset

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// *****************************************************************************
// Thread-Safe Configuration
// *****************************************************************************

var (
	info      Info
	infoMutex sync.RWMutex
)

// Info holds the config
type Info struct {
	// Folder is the parent folder path for the asset folder
	Folder string
}

// SetConfig stores the config
func SetConfig(i Info) {
	infoMutex.Lock()
	info = i
	infoMutex.Unlock()
}

// ResetConfig removes the config
func ResetConfig() {
	infoMutex.Lock()
	info = Info{}
	infoMutex.Unlock()
}

// Config returns the config
func Config() Info {
	infoMutex.RLock()
	i := info
	infoMutex.RUnlock()
	return i
}

// *****************************************************************************
// View Extend
// *****************************************************************************

// Extend returns a template.FuncMap
// JS returns JavaScript tag with timestamp
// CSS returns stylesheet tag with timestamp
func Extend(baseURI string) template.FuncMap {
	f := make(template.FuncMap)

	f["JS"] = func(s string) template.HTML {
		path, err := assetTimePath(baseURI, s)

		if err != nil {
			log.Println("JS Error:", err)
			return template.HTML("<!-- JS Error: " + s + " -->")
		}

		return template.HTML(`<script type="text/javascript" src="` + path + `"></script>`)
	}

	f["CSS"] = func(s string) template.HTML {
		path, err := assetTimePath(baseURI, s)

		if err != nil {
			log.Println("CSS Error:", err)
			return template.HTML("<!-- CSS Error: " + s + " -->")
		}

		return template.HTML(`<link rel="stylesheet" type="text/css" href="` + path + `" />`)
	}

	return f
}

// assetTimePath returns a URL with the proper base uri and timestamp appended
// Works for CSS and JS assets
// Determines if local or on the web
func assetTimePath(baseURI, resource string) (string, error) {
	if strings.HasPrefix(resource, "//") {
		return resource, nil
	}

	resource = strings.TrimLeft(resource, "/")

	abs, err := filepath.Abs(filepath.Join(Config().Folder, resource))
	if err != nil {
		return "", err
	}

	time, err := fileTime(abs)
	if err != nil {
		return "", err
	}

	return baseURI + resource + "?" + time, nil
}

// fileTime returns the modification time of the file
func fileTime(name string) (string, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return "", err
	}
	mtime := fi.ModTime().Unix()
	return fmt.Sprintf("%v", mtime), nil
}
