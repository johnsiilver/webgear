package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// dotFileHidingFile is the http.File use in dotFileHidingFileSystem.
// It is used to wrap the Readdir method of http.File so that we can
// remove files and directories that start with a period from its output.
type dotFileHidingFile struct {
	http.File
	exts map[string]bool
}

// Readdir is a wrapper around the Readdir method of the embedded File
// that filters out all files that start with a period in their name.
func (f dotFileHidingFile) Readdir(n int) (fis []os.FileInfo, err error) {
	files, err := f.File.Readdir(n)
	for _, file := range files { // Filters out the dot files
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if !f.exts[filepath.Ext(file.Name())] {
			continue
		}
		fis = append(fis, file)
	}
	return
}

// fileSystem is an http.FileSystem that hides hidden "dot files" and only serves files with
// certain extensions.
type fileSystem struct {
	http.FileSystem
	exts map[string]bool
}

// Open is a wrapper around the Open method of the embedded FileSystem
// that serves a 403 permission error when name has a file or directory
// with whose name starts with a period in its path.
func (fs fileSystem) Open(name string) (http.File, error) {
	if containsDotFile(name) { // If dot file, return 403 response
		return nil, os.ErrPermission
	}

	if !fs.exts[filepath.Ext(name)] {
		log.Printf("probe for non-allowed file extension: %s", filepath.Ext(name))
		return nil, os.ErrPermission
	}

	file, err := fs.FileSystem.Open(name)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return dotFileHidingFile{file, fs.exts}, err
}

// containsDotFile reports whether name contains a path element starting with a period.
// The name is assumed to be a delimited by forward slashes, as guaranteed
// by the http.FileSystem interface.
func containsDotFile(name string) bool {
	parts := strings.Split(name, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}
