package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
)

var validOutputFlags = []string{
	"table",
}

type NotePrintFunc = func(Note)
type LsConf struct {
	XTitle string
	XId string
	Tags []string
	TagsOr []string
	Output string
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
	notePrintFunc := NotePrint(lsConf.Output)
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
			notePrintFunc(note)
		}
	}
}

func LsFlags(args []string) LsConf {
	var conf LsConf
	lsf := flag.NewFlagSet("ls", 0)
	lsf.StringVar(&conf.XId, "xid", "", "Regex match id")
	lsf.StringVar(&conf.XTitle, "xtitle", "", "Regex match title")
	lsf.StringVar(&conf.Output, "o", "table", "Output format. [table]")
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

	if !slices.Contains(validOutputFlags, conf.Output) {
		log.Fatalf("Invalid output flag %s", conf.Output)
	}
	return conf
}

func NotePrint(outputType string) NotePrintFunc {
	switch outputType {
	default:
		return TablePrintNote()
	}
}

func TablePrintNote() NotePrintFunc {
	const (
		idWidth    = 24
		titleWidth = 30
		tagsWidth  = 30
	)

	headerFmt := fmt.Sprintf("%%-%ds  %%-%ds  %%-%ds  %%s\n",
	idWidth, titleWidth, tagsWidth)
	rowFmt := headerFmt

	fmt.Printf(headerFmt, "ID", "TITLE", "TAGS", "PATH")
	fmt.Println(strings.Repeat("-", idWidth+titleWidth+tagsWidth+10) + "----------------------------------------")

	return func(n Note) {
		title := truncate(n.Title, titleWidth)
		tags := truncate(strings.Join(n.Tags, ", "), tagsWidth)
		fmt.Printf(rowFmt, n.Id, title, tags, n.Path)
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
