package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type LsConf struct {
	XTitle string
	XId string
	Tags []string
	TagsOr []string
}

func LsCommand(conf AppConfig, args []string) {
	lsConf := LsFlags(args)
	rxId, err := regexp.Compile(lsConf.XId)
	if err != nil {
		log.Fatalf("ls command: Invalid regex id. regex=%s. %v", lsConf.XId, err)
	}
	rxTitle, err := regexp.Compile(lsConf.XTitle)
	if err != nil {
		log.Fatalf("ls command: Invalid regex title. regex=%s. %v", lsConf.XTitle, err)
	}

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
			if !NoteTagsFilter(note, lsConf.Tags, lsConf.TagsOr) {
				continue
			}
			if !rxId.MatchString(note.Id) {
				continue
			}
			if !rxTitle.MatchString(note.Title) {
				continue
			}
			fmt.Printf("%s: %s\n", note.Id, note.Path)
		}
	}
}

func LsFlags(args []string) LsConf {
	var conf LsConf
	lsf := flag.NewFlagSet("ls", 0)
	lsf.StringVar(&conf.XId, "xid", ".", "Regex match id")
	lsf.StringVar(&conf.XTitle, "xtitle", ".", "Regex match title")
	tags := lsf.String("t", "", "Tags AND")
	tagsOr := lsf.String("tor", "", "Tags OR")
	err := lsf.Parse(args)
	if err != nil {
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
