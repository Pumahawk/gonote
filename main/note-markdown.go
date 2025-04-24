package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/go-git/go-git/v5"
	yaml "github.com/goccy/go-yaml"
)

func MarkdownNote(repo *git.Repository, absolutePath, relativePath string, r io.Reader) (Note, error) {
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
				lastUpdateTime, err := GetLastUpdate(repo, relativePath)
				if err != nil {
					log.Printf("Unable to get last update time from file %s. %v", relativePath, err)
				}
				note.lastUpdate = lastUpdateTime
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
	lastUpdate *time.Time
}

func (n *NoteMarkdown) LastUpdate() *time.Time {
	return n.lastUpdate
}

func (n *NoteMarkdown) OpenRef() string {
	return fmt.Sprintf("%s%c%d", n.AbsoluteFilePath, ':', 0)
}

func GetLastUpdate(repo *git.Repository, path string) (*time.Time, error) {
	ci, err := repo.Log(&git.LogOptions{
		PathFilter: func(s string) bool {
			return s == path
		},
	})
	if err != nil {
		return nil, fmt.Errorf("markdown: Unable to read repository logs, path: %s. %w", path, err)
	}
	c, err := ci.Next()
	if err != nil {
		return nil, fmt.Errorf("markdown: unable to get next commit. path: %s. %w", path, err)
	}
	return &c.Author.When, nil
}
