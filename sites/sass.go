package sites

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/osteele/gojekyll/helpers"

	libsass "github.com/wellington/go-libsass"
)

// IsSassPath returns a boolean indicating whether the file is a Sass (".sass" or ".scss") file.
func (s *Site) IsSassPath(name string) bool {
	return strings.HasSuffix(name, ".sass") || strings.HasSuffix(name, ".scss")
}

// CopySassFileIncludes copies sass partials into a temporary directory,
// removing initial underscores.
// TODO delete the temp directory when done
func (s *Site) CopySassFileIncludes() error {
	// TODO use libsass.ImportsOption instead?
	if s.sassTempDir == "" {
		dir, err := ioutil.TempDir(os.TempDir(), "_sass")
		if err != nil {
			return err
		}
		s.sassTempDir = dir
	}

	src := filepath.Join(s.Source, "_sass")
	dst := s.sassTempDir
	err := filepath.Walk(src, func(from string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, err := filepath.Rel(src, from)
		if err != nil {
			return err
		}
		to := filepath.Join(dst, strings.TrimPrefix(rel, "_"))
		return helpers.CopyFileContents(to, from, 0644)
	})
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// SassIncludePaths returns an array of sass include directories.
func (s *Site) SassIncludePaths() []string {
	return []string{s.sassTempDir}
}

// WriteSass converts a SASS file and writes it to w.
func (s *Site) WriteSass(w io.Writer, b []byte) error {
	comp, err := libsass.New(w, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	err = comp.Option(libsass.IncludePaths(s.SassIncludePaths()))
	if err != nil {
		log.Fatal(err)
	}
	return comp.Run()
}