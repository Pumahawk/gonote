package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	yaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

// Read all notes from YAML file
// Use abstolutePath to Open file.
// Use path as referiment for git interaction
func ReadYamlNotes(repo *git.Repository, abstolutePath, relativePath string) ([]Note, error) {
	file, err := os.Open(abstolutePath)
	if err != nil {
		return nil, fmt.Errorf("yaml notes: Unable to open file %s. %w", abstolutePath, err)
	}

	content, lineCounter, err := ReadAndCountLines(file)
	if err != nil {
		return nil, fmt.Errorf("yaml notes: Unable to read all contents from file %s. %w", abstolutePath, err)
	}

	var notes []Note
	f, _ := parser.ParseBytes(content, parser.ParseComments)
	for _, doc := range f.Docs {
		if mapNode, ok := doc.Body.(*ast.MappingNode); ok {
			for _, v := range mapNode.Values {
				if v.Key.String() == "notes" {
					if sqn, ok := v.Value.(*ast.SequenceNode); ok {
						var previousNote *NoteYaml
						for _, n := range sqn.Values {
							if n, ok := n.(*ast.MappingNode); ok {
								note := &NoteYaml{
									BaseNote: &BaseNote{},
								}
								if err := yaml.NodeToValue(n, note.BaseNote); err != nil {
									return nil, fmt.Errorf("Unable to read yaml note, path=%s. %w", abstolutePath, err)
								}
								note.FilePath = relativePath
								note.AbsoluteFilePath = abstolutePath
								note.LineStartNumber = n.GetToken().Position.Line
								if previousNote != nil {
									previousNote.LineEndNumber = note.LineStartNumber - 1
									setLatUpdateTimeNote(previousNote, repo, relativePath, previousNote.LineStartNumber, previousNote.LineEndNumber)
								}
								notes = append(notes, note)
								previousNote = note
							}
						}
						if previousNote != nil {
							previousNote.LineEndNumber = lineCounter
							setLatUpdateTimeNote(previousNote, repo, relativePath, previousNote.LineStartNumber, previousNote.LineEndNumber)
						}
					}
				}
			}
		}
	}
	return notes, nil
}

type LineCounterReader struct {
	Origin io.Reader
	Counter int
}

func NewLineCounterReader(r io.Reader) LineCounterReader {
	return LineCounterReader{
		Origin: r,
	}
}

func (lr *LineCounterReader) Read(p []byte) (n int, err error) {
	n, err = lr.Origin.Read(p)
	if err != nil {
		return
	}

	for i := 0; i < n; i++ {
		if p[i] == '\n' {
			lr.Counter++
		}
	}
	return
}

// Read all bytes from reder and count \n occurences
func ReadAndCountLines(r io.Reader) ([]byte, int, error){
	var bf bytes.Buffer
	lineCounter := NewLineCounterReader(r)
	if _, err := io.Copy(&bf, &lineCounter); err != nil {
		return nil, -1, err
	}
	return bf.Bytes(), lineCounter.Counter, nil
}

type NoteYaml struct {
	*BaseNote
	FilePath string
	AbsoluteFilePath string
	LineStartNumber int
	LineEndNumber int
	lastUpdate *time.Time
}

func (n *NoteYaml) LastUpdate() *time.Time {
	return n.lastUpdate
}

func (n *NoteYaml) OpenRef() string {
	return fmt.Sprintf("%s:%d", n.AbsoluteFilePath, n.LineStartNumber)
}

func setLatUpdateTimeNote(note *NoteYaml, repo *git.Repository, relativePath string, lineStart, lineEnd int) {
	t, err := GetLastUpdateLine(repo, relativePath, lineStart - 1, lineEnd - 1)
	if err != nil {
		log.Printf("yaml note: Unable to retrieve last update from note %s:%d start=%d end=%d. %v", note.FilePath, note.LineStartNumber, note.LineStartNumber, note.LineEndNumber, err)
	} else {
		note.lastUpdate = t
	}
}

func GetLastUpdateLine(repo *git.Repository, path string, lineStart, lineEnd int) (*time.Time, error) {
	opt := git.LogOptions{
		PathFilter: func(s string) bool {
			return path == s
		},
	}
	it, err := repo.Log(&opt)
	if err != nil {
		return nil, fmt.Errorf("git: unable to get Log %w", err)
	}
	defer it.Close()
	c, err := it.Next()
	if err != nil {
		return nil, fmt.Errorf("git: unable to get next log. %w", err)
	}

	br, err := git.Blame(c, path)
	if err != nil {
		return nil, fmt.Errorf("git: unable to get blame. %w", err)
	}
	if lineStart < len(br.Lines) && lineEnd < len(br.Lines) {
			mxd := br.Lines[lineStart].Date
			for _, l := range br.Lines[lineStart:lineEnd] {
				if l.Date.Unix() > mxd.Unix() {
					mxd = l.Date
				}
			}
			return &mxd, nil
	} else {
		return nil, fmt.Errorf("git blame: invalid range of lines. lines: %d start: %d, end: %d.", len(br.Lines), lineStart, lineEnd)
	}
}
