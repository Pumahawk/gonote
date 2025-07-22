package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

type InfoConf struct {
	XTitle  *regexp.Regexp
	XId     *regexp.Regexp
	Tags    []string
	TagsNot []string
	TagsOr  []string
}

type FilePath struct {
	Absolute string
	Relative string
}

func InfoCommand(conf AppConfig, args []string) {
	infoConf, args := InfoFlags(args)

	files, err := FindAllNotesFiles(conf.RootPath, args)
	if err != nil {
		log.Fatalf("info: Unable to read notes files %v", err)
	}

	noteCount := 0
	tagsCount := make(map[string]int)

	for _, file := range files {
		notes, err := GetNoteData(file.Absolute, file.Relative)
		if err != nil {
			log.Fatalf("info: Unable to read file %s. %v", file, err)
		}
		for _, note := range notes {
			if !NoteTagsFilter(note, infoConf.Tags, infoConf.TagsNot, infoConf.TagsOr) {
				continue
			}
			if !infoConf.XId.MatchString(note.Id()) {
				continue
			}
			if !infoConf.XTitle.MatchString(note.Title()) {
				continue
			}
			noteCount++
			for _, tag := range note.Tags() {
				tagsCount[tag]++
			}
		}
	}

	fmt.Printf("Notes: %d\n", noteCount)
	fmt.Println("Tags: ")
	var tagsCountS []string
	for k, v := range tagsCount {
		tagsCountS = append(tagsCountS, fmt.Sprintf("\t%s: %d", k, v))
	}
	sort.Strings(tagsCountS)
	for _, ts := range tagsCountS {
		fmt.Println(ts)
	}
}

func InfoFlags(args []string) (InfoConf, []string) {
	var conf InfoConf
	lsf := flag.NewFlagSet("info", 0)
	xId := lsf.String("xid", "", "Regex match id")
	xTitle := lsf.String("xtitle", "", "Regex match title")
	tags := lsf.String("t", "", "Tags AND")
	tagsNot := lsf.String("tn", "", "Tags NOT")
	tagsOr := lsf.String("tor", "", "Tags OR")
	err := lsf.Parse(args)
	if err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if *tags != "" {
		conf.Tags = strings.Split(*tags, ",")
	}

	if *tagsNot != "" {
		conf.Tags = strings.Split(*tagsNot, ",")
	}

	if *tagsOr != "" {
		conf.TagsOr = strings.Split(*tagsOr, ",")
	}

	rxId, err := regexp.Compile(*xId)
	if err != nil {
		log.Fatalf("ls command: Invalid regex id. regex=%s. %v", rxId, err)
	}
	conf.XId = rxId

	rxTitle, err := regexp.Compile(*xTitle)
	if err != nil {
		log.Fatalf("ls command: Invalid regex title. regex=%s. %v", rxTitle, err)
	}
	conf.XTitle = rxTitle

	return conf, lsf.Args()
}
