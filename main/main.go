package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"slices"

	yaml "github.com/goccy/go-yaml"
)

func main() {
	conf, args := AppParseFlags()
	if len(args) > 0 {
		command := args[0]
		switch command {
		case "ls":
			LsCommand(conf, args[1:])
		case "show":
			ShowCommand(conf, args[1:])
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
	fmt.Println("stat - Print result information")
	fmt.Println("show - Print note details")
	fmt.Println("edit - Open note with favourite editor")
}

func GetNoteData(filePath string) ([]Note, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("main: Unable to open note file. %w", err)
	}
	defer file.Close()

	if regexp.MustCompile("\\.yaml$").MatchString(filePath) {
		return YamlNotes(filePath, file)
	}

	if regexp.MustCompile("\\.md$").MatchString(filePath) {
		note, err := MarkdownNote(filePath, file) 
		if err != nil {
			return nil, fmt.Errorf("main: Unable to read Markdown note: %s. %w", filePath, err)
		}
		return []Note{*note}, nil
	}

	return nil, fmt.Errorf("main: Invalid note extension, supported [md, yaml]. Path %s", filePath)
}

func YamlNotes(path string, file *os.File) ([]Note, error) {
	var note NoteFileYaml
	if err := yaml.NewDecoder(file).Decode(&note); err != nil {
		return nil, fmt.Errorf("Unable to read yaml note, path=%s. %w", path, err)
	}
	var notes []Note
	for _, note := range note.Notes {
		note.PathY = path
		notes = append(notes, note)
	}
	return notes, nil
}

func MarkdownNote(path string, file *os.File) (*NoteMd, error) {
	var buf bytes.Buffer
	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return nil, fmt.Errorf("markdown: Unable to read first line, path=%s. %w", path, scanner.Err())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("markdown: Unable to parse first line markdown file %s. %w", path, err)
	}
	if scanner.Text() != "---" {
		return nil, fmt.Errorf("markdown: Invalid first line. Expected ---. path=%s", path)
	}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil{
			return nil, fmt.Errorf("markdown: Unable to parse markdown file %s. %w", path, err)
		} else if scanner.Text() == "---" {
			var note NoteMd
			if err := yaml.Unmarshal(buf.Bytes(), &note); err != nil {
				return nil, fmt.Errorf("markdown: Unable to parse metadata in note. %w", err)
			} else {
				note.PathM = path
				return &note, nil
			}
		} else {
			buf.WriteString(scanner.Text())
			buf.WriteByte('\n')
		}
	}
	return nil, fmt.Errorf("markdown: Not found metadata note")
}

func FindAllNotesFiles(basePath string) ([]string, error) {
	var files []string
	root := os.DirFS(basePath)
	fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			if regexp.MustCompile("\\.yaml$").MatchString(d.Name()) || regexp.MustCompile("\\.md$").MatchString(d.Name()) {
				if basePath == "." {
					files = append(files, path)
				} else {
					files = append(files, fmt.Sprintf("%s%c%s", basePath, os.PathSeparator, path))
				}
			}
		} else if regexp.MustCompile("^\\..").MatchString(d.Name()) {
			return fs.SkipDir
		}
		return nil
	})
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
