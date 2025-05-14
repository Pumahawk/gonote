package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	yaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

// Read all notes from YAML file
// Use abstolutePath to Open file.
// Use path as referiment for git interaction
func ReadYamlNotes(abstolutePath, relativePath string) ([]Note, error) {
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
								note.absoluteFilePath = abstolutePath
								position := n.GetToken().Position
								note.lineStartNumber = position.Line
								note.offset = position.Offset
								if previousNote != nil {
									previousNote.lineEndNumber = note.lineStartNumber - 1
									previousNote.size = note.offset - previousNote.offset
								}
								notes = append(notes, note)
								previousNote = note
							}
						}
						if previousNote != nil {
							previousNote.lineEndNumber = lineCounter.Counter
							previousNote.size = lineCounter.Size - previousNote.offset
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
	Size int
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

	lr.Size += n

	for i := 0; i < n; i++ {
		if p[i] == '\n' {
			lr.Counter++
		}
	}
	return
}

// Read all bytes from reder and count \n occurences
func ReadAndCountLines(r io.Reader) ([]byte, *LineCounterReader, error){
	var bf bytes.Buffer
	lineCounter := NewLineCounterReader(r)
	if _, err := io.Copy(&bf, &lineCounter); err != nil {
		return nil, nil, err
	}
	return bf.Bytes(), &lineCounter, nil
}

type NoteYaml struct {
	*BaseNote
	FilePath string
	absoluteFilePath string
	lineStartNumber int
	lineEndNumber int
	offset int
	size int
}

func (n *NoteYaml) OpenRef() string {
	return fmt.Sprintf("%s:%d", n.absoluteFilePath, n.lineStartNumber)
}
