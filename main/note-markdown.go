package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	yaml "github.com/goccy/go-yaml"
)

func MarkdownNote(absolutePath, relativePath string, r io.Reader) (Note, error) {
	var buf bytes.Buffer
	scanner := bufio.NewScanner(r)

	if !scanner.Scan() {
		return nil, fmt.Errorf("markdown: Unable to read first line. %w", scanner.Err())
	}

	if err := scanner.Err(); err != nil {
		return nil, NotANote
	}
	if scanner.Text() != "---" {
		return nil, NotANote
	}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("markdown: Unable to parse markdown. %w", err)
		} else if scanner.Text() == "---" {
			note := NoteMarkdown{
				BaseNote: &BaseNote{},
			}
			if err := yaml.Unmarshal(buf.Bytes(), note.BaseNote); err != nil {
				return nil, fmt.Errorf("markdown: Unable to parse metadata in note. %w", err)
			} else {
				note.FilePath = relativePath
				note.AbsoluteFilePath = absolutePath
				return &note, nil
			}
		} else {
			buf.WriteString(scanner.Text())
			buf.WriteByte('\n')
		}
	}
	return nil, NotANote
}

type NoteMarkdown struct {
	*BaseNote
	FilePath string
	AbsoluteFilePath string
}

func (n *NoteMarkdown) OpenRef() string {
	return fmt.Sprintf("%s%c%d", n.AbsoluteFilePath, ':', 0)
}
