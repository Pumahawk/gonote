package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"regexp"
	"slices"

	"github.com/go-git/go-git/v5"
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

func GetNoteData(repo *git.Repository, absolutePath, relativePath string) ([]Note, error) {
	if regexp.MustCompile("\\.yaml$").MatchString(relativePath) {
		return ReadYamlNotes(repo, absolutePath, relativePath)
	}
	if regexp.MustCompile("\\.md$").MatchString(relativePath) {
		file, err := os.Open(absolutePath)
		if err != nil {
			return nil, fmt.Errorf("main: Unable to open note file. %w", err)
		}
		defer file.Close()
		note, err := MarkdownNote(repo, absolutePath, relativePath, file)

		if err == NotANote {
			return []Note{}, nil
		}

		if err != nil {
			log.Printf("main: Unable to read Markdown note: %s. %v", relativePath, err)
			return []Note{}, nil
		}

		return []Note{note}, nil
	}

	return nil, fmt.Errorf("main: Invalid note extension, supported [md, yaml]. Path %s", relativePath)
}

func FindAllNotesFiles(basePath string, subPath []string) ([]FilePath, error) {
	var files []FilePath
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
					var absolute string
					if basePath == "." {
						absolute = path
					} else {
						absolute = fmt.Sprintf("%s%c%s", basePath, os.PathSeparator, path)
					}
					files = append(files, FilePath{
						Absolute: absolute,
						Relative: path,
					})
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


