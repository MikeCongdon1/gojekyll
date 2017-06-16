package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	libsass "github.com/wellington/go-libsass"
)

// IsSassPath returns a boolean indicating whether the file is a Sass (".sass" or ".scss") file.
func (s *Site) IsSassPath(path string) bool {
	return strings.HasSuffix(path, ".sass") || strings.HasSuffix(path, ".scss")
}

func (p *DynamicPage) writeSass(w io.Writer, data []byte) error {
	comp, err := libsass.New(w, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	err = comp.Option(libsass.IncludePaths(site.SassIncludePaths()))
	if err != nil {
		log.Fatal(err)
	}
	return comp.Run()
}

// CopySassFileIncludes copies sass partials into a temporary directory,
// removing initial underscores.
// TODO delete the temp directory when done
func (s *Site) CopySassFileIncludes() {
	// TODO use libsass.ImportsOption instead?
	if site.sassTempDir == "" {
		d, err := ioutil.TempDir(os.TempDir(), "_sass")
		if err != nil {
			panic(err)
		}
		site.sassTempDir = d
	}

	src := filepath.Join(s.Source, "_sass")
	dst := site.sassTempDir
	err := filepath.Walk(src, func(from string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, err := filepath.Rel(src, from)
		if err != nil {
			panic(err)
		}
		to := filepath.Join(dst, strings.TrimPrefix(rel, "_"))
		return CopyFileContents(to, from, 0644)
	})
	if err != nil {
		panic(err)
	}
}

// SassIncludePaths returns an array of sass include directories.
func (s *Site) SassIncludePaths() []string {
	if site.sassTempDir == "" {
		site.CopySassFileIncludes()
	}
	s.CopySassFileIncludes()
	return []string{site.sassTempDir}
}