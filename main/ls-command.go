package main

import (
	"fmt"
	"log"
)

func LsCommand(conf AppConfig, args []string) {
	files, err := FindAllNotesFiles(conf.RootPath)
	if err != nil {
		log.Fatalf("main: Unable to read notes files %v", err)
	}
	for _, file := range files {
		notes, err := GetNoteData(file)
		if err != nil {
			log.Fatalf("main: Unable to read file %s. %v", file, err)
		}
		for _, note := range notes {
			fmt.Printf("%s: %s\n", note.Id, note.Path)
		}
	}
}
