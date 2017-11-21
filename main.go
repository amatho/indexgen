package main

import (
	"bufio"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Entry struct {
	Name string
	Path string
	Type string
}

type TemplateData struct {
	Folder  string
	Entries []Entry
}

var excludedFiles = [...]string{"index.html", ".DS_Store", ".git"}

func stringInSlice(target string, slice []string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var pathPrefix string
	if len(os.Args) > 1 {
		pathPrefix = os.Args[1]
	}

	wd, wdErr := os.Getwd()
	check(wdErr)
	wdName := filepath.Base(wd)

	var files []os.FileInfo
	var readErr error
	if pathPrefix == "" {
		files, readErr = ioutil.ReadDir(wd)
	} else {
		files, readErr = ioutil.ReadDir(pathPrefix)
	}
	check(readErr)

	var entries []Entry
	for _, f := range files {
		if stringInSlice(f.Name(), excludedFiles[:]) {
			continue
		}
		if filepath.Ext(f.Name()) != ".html" || filepath.Ext(f.Name()) != "" {
			continue
		}

		var path string
		var fileType string

		if f.IsDir() {
			fileType = "folder"
			path = f.Name() + "/"
		} else {
			fileType = "file"
			path = f.Name()
		}

		entries = append(entries, Entry{f.Name(), path, fileType})
	}

	if len(entries) == 0 {
		entry := Entry{"No folders or .html files in directory", "", ""}
		entries = append(entries, entry)
	}

	t, tErr := template.New("template").Parse(templateString)
	check(tErr)

	f, fErr := os.Create(filepath.Join(pathPrefix, "index.html"))
	check(fErr)
	defer f.Close()

	w := bufio.NewWriter(f)
	check(t.Execute(w, TemplateData{wdName, entries}))
	w.Flush()
}
