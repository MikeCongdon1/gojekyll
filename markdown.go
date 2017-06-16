package main

import (
	"path/filepath"
	"strings"
)

// IsMarkdown returns a boolean indicating whether the file is a Markdown file, according to the current project.
func (s *Site) IsMarkdown(path string) bool {
	ext := filepath.Ext(path)
	return site.MarkdownExtensions()[strings.TrimLeft(ext, ".")]
}

// MarkdownExtensions returns a set of markdown extension, without the final dots.
func (s *Site) MarkdownExtensions() map[string]bool {
	extns := strings.SplitN(s.config.MarkdownExt, `,`, -1)
	return stringArrayToMap(extns)
}