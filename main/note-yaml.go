package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/go-git/go-git/v5"
	yaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

// Read all notes from YAML file
// Use abstolutePath to Open file.
// Use path as referiment for git interaction
func ReadYamlNotes(repo *git.Repository, abstolutePath, path string) ([]Note, error) {
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
								// TODO detine correct mapping from Node to Value
								if err := yaml.NodeToValue(n, note); err != nil {
									return nil, fmt.Errorf("Unable to read yaml note, path=%s. %w", path, err)
								}
								note.FilePath = path
								note.LineStartNumber = n.GetToken().Position.Line
								if previousNote != nil {
									previousNote.LineEndNumber = note.LineStartNumber - 1
									setLatUpdateTimeNote(previousNote, repo, path, previousNote.LineStartNumber, previousNote.LineEndNumber)
								}
								notes = append(notes, note)
								previousNote = note
							}
						}
						if previousNote != nil {
							previousNote.LineEndNumber = lineCounter
							setLatUpdateTimeNote(previousNote, repo, path, previousNote.LineStartNumber, previousNote.LineEndNumber)
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
	LineStartNumber int
	LineEndNumber int
}

func (n *NoteYaml) Note() []byte {
	return nil
}
