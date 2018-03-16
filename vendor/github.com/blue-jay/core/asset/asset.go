// Package asset contains FuncMap for template packages to provide asset paths
// with timestamps.
package asset

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Info holds the config.
type Info struct {
	// Folder is the parent folder path for the asset folder
	Folder string
}

// *****************************************************************************
// FuncMap for Template Packages
// *****************************************************************************

// Map returns a template.FuncMap. JS returns JavaScript tag with a
// timestamp. CSS returns stylesheet tag with a timestamp.
func (c Info) Map(baseURI string) template.FuncMap {
	f := make(template.FuncMap)

	f["JS"] = func(fpath string) template.HTML {
		path, err := c.assetTimePath(baseURI, fpath)

		if err != nil {
			log.Println("JS Error:", err)
			return template.HTML("<!-- JS Error: " + fpath + " -->")
		}

		return template.HTML(`<script type="text/javascript" src="` + path + `"></script>`)
	}

	f["CSS"] = func(fpath, media string) template.HTML {
		path, err := c.assetTimePath(baseURI, fpath)

		if err != nil {
			log.Println("CSS Error:", err)
			return template.HTML("<!-- CSS Error: " + fpath + " -->")
		}

		return template.HTML(fmt.Sprintf(`<link media="%v" rel="stylesheet" type="text/css" href="%v" />`, media, path))
	}

	f["IMG"] = func(fpath string, alt string) template.HTML {
		path, err := c.assetTimePath(baseURI, fpath)

		if err != nil {
			log.Println("IMG Error:", err)
			return template.HTML("<!-- IMG Error: " + fpath + " -->")
		}

		return template.HTML(fmt.Sprintf(`<img src="%v" alt="%v" />`, path, alt))
	}

	return f
}

// *****************************************************************************
// Helpers
// *****************************************************************************

// assetTimePath returns a URL with the proper base URI and timestamp appended.
// Works for CSS and JS assets and determines if local or on the web by the
// number of slashes at the beginning of the string. A prefix of // is web and
// / is local.
func (c Info) assetTimePath(baseURI, resource string) (string, error) {
	if strings.HasPrefix(resource, "//") {
		return resource, nil
	}

	resource = strings.TrimLeft(resource, "/")

	abs, err := filepath.Abs(filepath.Join(c.Folder, resource))
	if err != nil {
		return "", err
	}

	time, err := fileTime(abs)
	if err != nil {
		return "", err
	}

	return baseURI + resource + "?" + time, nil
}

// fileTime returns the modification time of the file.
func fileTime(name string) (string, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return "", err
	}
	mtime := fi.ModTime().Unix()
	return fmt.Sprintf("%v", mtime), nil
}
