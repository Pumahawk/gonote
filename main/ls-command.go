package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type LsConf struct {
	Tags []string
	TagsOr []string
}

func LsCommand(conf AppConfig, args []string) {
	lsConf := LsFlags(args)
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
			if NoteTagsFilter(note, lsConf.Tags, lsConf.TagsOr) {
				fmt.Printf("%s: %s\n", note.Id, note.Path)
			}
		}
	}
}

func LsFlags(args []string) LsConf {
	var conf LsConf
	lsf := flag.NewFlagSet("ls", 0)
	tags := lsf.String("t", "", "Tags AND")
	tagsOr := lsf.String("tor", "", "Tags OR")
	err := lsf.Parse(args)
	if err == flag.ErrHelp {
		os.Exit(0)
	}

	if *tags != "" {
		conf.Tags = strings.Split(*tags, ",")
	}

	if *tagsOr != "" {
		conf.TagsOr = strings.Split(*tagsOr, ",")
	}

	return conf
}
