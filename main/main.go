package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"

	yaml "github.com/goccy/go-yaml"
)

func main() {
	if len(os.Args) > 1 {
		command := os.Args[1]
		args := os.Args[2:]
		switch command {
		case "ls":
			LsCommand(args)
		default:
			NotFoundCommand(command)
		}
	} else {
		PrintHelpMessage()
	}
}

func LsCommand(args []string) {
	files, err := FindAllNotesFiles(".")
	if err != nil {
		log.Fatalf("main: Unable to read notes files %v", err)
	}
	for _, file := range files {
		notes, err := GetNoteData(file)
		if err != nil {
			log.Fatalf("main: Unable to read file %s. %v", file, err)
		}
		for _, note := range notes {
			fmt.Printf("%s\n", note.Id)
		}
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
		return YamlNotes(file)
	}

	if regexp.MustCompile("\\.md$").MatchString(filePath) {
		note, err := MarkdownNote(file) 
		if err != nil {
			return nil, fmt.Errorf("main: Unable to read Markdown note: %s. %w", filePath, err)
		}
		return []Note{*note}, nil
	}

	return nil, fmt.Errorf("main: Invalid note extension, supported [md, yaml]. Path %s", filePath)
}

func YamlNotes(file *os.File) ([]Note, error) {
	var note NoteFile
	if err := yaml.NewDecoder(file).Decode(&note); err != nil {
		return nil, fmt.Errorf("Unable to read yaml note, path=%s. %w", file.Name(), err)
	}
	var notes []Note
	for _, note := range note.Notes {
		note.Path = file.Name()
		notes = append(notes, note)
	}
	return notes, nil
}

func MarkdownNote(file *os.File) (*Note, error) {
	var buf bytes.Buffer
	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return nil, fmt.Errorf("markdown: Unable to read first line, path=%s. %w", file.Name(), scanner.Err())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("markdown: Unable to parse first line markdown file %s. %w", file.Name(), err)
	}
	if scanner.Text() != "---" {
		return nil, fmt.Errorf("markdown: Invalid first line. Expected ---. path=%s", file.Name())
	}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil{
			return nil, fmt.Errorf("markdown: Unable to parse markdown file %s. %w", file.Name(), err)
		} else if scanner.Text() == "---" {
			var note Note
			if err := yaml.Unmarshal(buf.Bytes(), &note); err != nil {
				return nil, fmt.Errorf("markdown: Unable to parse metadata in note. %w", err)
			} else {
				note.Path = file.Name()
				return &note, nil
			}
		} else {
			buf.WriteString(scanner.Text())
			buf.WriteByte('\n')
		}
	}
	return nil, fmt.Errorf("markdown: Not found metadata note")
}

func FindAllNotesFiles(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("Unable extract markdown notes. path=%s", path)
	}
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			if regexp.MustCompile("\\.yaml$").MatchString(entry.Name()) || regexp.MustCompile("\\.md$").MatchString(entry.Name()) {
				files = append(files, entry.Name())
			}
		}
	}
	return files, nil
}
