package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"regexp"
	"slices"

	yaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

var NotANote = fmt.Errorf("Invalid Markdown note")

func main() {
	conf, args := AppParseFlags()
	if len(args) > 0 {
		command := args[0]
		switch command {
		case "ls":
			LsCommand(conf, args[1:])
		case "show":
			ShowCommand(conf, args[1:])
		case "info":
			InfoCommand(conf, args[1:])
		default:
			NotFoundCommand(command)
		}
	} else {
		PrintHelpMessage()
	}
}

func NotFoundCommand(command string) {
	fmt.Printf("Command not found. %s\n", command)
}

func PrintHelpMessage() {
	fmt.Println("Commands:")
	fmt.Println()
	fmt.Println("ls - Print all notes")
	fmt.Println("info - Print result information")
	fmt.Println("show - Print note details")
}

func GetNoteData(filePath string) ([]Note, error) {
	if regexp.MustCompile("\\.yaml$").MatchString(filePath) {
		return YamlNotes(filePath)
	}
	if regexp.MustCompile("\\.md$").MatchString(filePath) {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("main: Unable to open note file. %w", err)
		}
		defer file.Close()
		note, err := MarkdownNote(file)

		if err == NotANote {
			return []Note{}, nil
		}

		if err != nil {
			log.Printf("main: Unable to read Markdown note: %s. %v", filePath, err)
			return []Note{}, nil
		}

		note.PathM = filePath

		return []Note{*note}, nil
	}

	return nil, fmt.Errorf("main: Invalid note extension, supported [md, yaml]. Path %s", filePath)
}

func YamlNotes(path string) ([]Note, error) {
	var notes []Note
	f, _ := parser.ParseFile(path, parser.ParseComments)
	for _, doc := range f.Docs {
		if mapNode, ok := doc.Body.(*ast.MappingNode); ok {
			for _, v := range mapNode.Values {
				if v.Key.String() == "notes" {
					if sqn, ok := v.Value.(*ast.SequenceNode); ok {
						for _, n := range sqn.Values {
							var note NoteYaml
							if err := yaml.NodeToValue(n, &note); err != nil {
								return nil, fmt.Errorf("Unable to read yaml note, path=%s. %w", path, err)
							}
							note.pathY = path
							note.lineY = n.GetToken().Position.Line
							notes = append(notes, note)
						}
					}
				}
			}
		}
	}
	return notes, nil
}

func MarkdownNote(r io.Reader) (*NoteMd, error) {
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
			var note NoteMd
			if err := yaml.Unmarshal(buf.Bytes(), &note); err != nil {
				return nil, fmt.Errorf("markdown: Unable to parse metadata in note. %w", err)
			} else {
				return &note, nil
			}
		} else {
			buf.WriteString(scanner.Text())
			buf.WriteByte('\n')
		}
	}
	return nil, NotANote
}

func FindAllNotesFiles(basePath string, subPath []string) ([]string, error) {
	var files []string
	root := os.DirFS(basePath)
	if len(subPath) == 0 {
		subPath = []string{"."}
	}
	for _, sb := range subPath {
		fs.WalkDir(root, sb, func(path string, d fs.DirEntry, err error) error {
			if d == nil {
				return nil
			}
			if !d.IsDir() {
				if regexp.MustCompile("\\.yaml$").MatchString(d.Name()) || regexp.MustCompile("\\.md$").MatchString(d.Name()) {
					var file string
					if basePath == "." {
						file = path
					} else {
						file = fmt.Sprintf("%s%c%s", basePath, os.PathSeparator, path)
					}
					if !slices.Contains(files, file) {
						files = append(files, file)
					}
				}
			} else if regexp.MustCompile("^\\..").MatchString(d.Name()) {
				return fs.SkipDir
			}
			return nil
		})
	}
	return files, nil
}

func NoteTagsFilter(note Note, tags, tagsOr []string) bool {
	if len(tags) > 0 {
		for _, t := range tags {
			if !slices.Contains(note.Tags(), t) {
				return false
			}
		}
	}
	if len(tagsOr) > 0 {
		for _, t := range tagsOr {
			if slices.Contains(note.Tags(), t) {
				return true
			}
		}
		return false
	}
	return true
}

func GetLinksFromText(r io.Reader) ([]NoteId, error) {
	var b bytes.Buffer
	if _, err := io.Copy(&b, r); err != nil {
		return nil, fmt.Errorf("main get links: Unable to retrieve links from note reader. %w", err)
	}
	s := regexp.MustCompile("\\[\\[(.*)\\]\\]").FindAllStringSubmatch(b.String(), -1)
	var nlinks []NoteId
	for _, m := range s {
		if len(m) > 0 {
			nlinks = append(nlinks, m[1])
		}
	}
	return nlinks, nil
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

