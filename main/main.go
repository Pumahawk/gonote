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
	"time"

	"github.com/go-git/go-git/v5"
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

func GetNoteData(repo *git.Repository, filePath string) ([]Note, error) {
	if regexp.MustCompile("\\.yaml$").MatchString(filePath) {
		return YamlNotes(repo, filePath)
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

func YamlNotes(repo *git.Repository, path string) ([]Note, error) {
	var notes []Note
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("yaml notes: Unable to open file %s. %w", path, err)
	}

	var bf bytes.Buffer
	lineCounter := NewLineCounterReader(file)
	io.Copy(&bf, &lineCounter)
	file.Close()

	f, _ := parser.ParseBytes(bf.Bytes(), parser.ParseComments)
	for _, doc := range f.Docs {
		if mapNode, ok := doc.Body.(*ast.MappingNode); ok {
			for _, v := range mapNode.Values {
				if v.Key.String() == "notes" {
					if sqn, ok := v.Value.(*ast.SequenceNode); ok {
						var pn *NoteYaml
						for _, n := range sqn.Values {
							if n, ok := n.(*ast.MappingNode); ok {
								note := &NoteYaml{}
								if err := yaml.NodeToValue(n, note); err != nil {
									return nil, fmt.Errorf("Unable to read yaml note, path=%s. %w", path, err)
								}
								note.pathY = path
								note.lineY = n.GetToken().Position.Line
								if pn != nil {
									pn.lineEndY = note.lineY - 1
									setLatUpdateTimeNote(pn, repo, path, pn.lineY, pn.lineEndY)
								}
								notes = append(notes, note)
								pn = note
							}
						}
						if pn != nil {
							pn.lineEndY = lineCounter.Counter
							setLatUpdateTimeNote(pn, repo, path, pn.lineY, pn.lineEndY)
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

func setLatUpdateTimeNote(note *NoteYaml, repo *git.Repository, relativePath string, lineStart, lineEnd int) {
	t, err := GetLastUpdateLine(repo, relativePath, lineStart - 1, lineEnd - 1)
	if err != nil {
		log.Printf("yaml note: Unable to retrieve last update from note %s:%d start=%d end=%d. %v", note.Path(), note.Line(), lineStart, lineEnd, err)
	} else {
		note.BaseNote.lastUpdate = t
	}
}

