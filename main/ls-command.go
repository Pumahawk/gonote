package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
)

var validOutputFlags = []string{
	"table",
}

type NotePrintFunc = func(Note, []NoteId)
type LsConf struct {
	Links           bool
	XTitle          *regexp.Regexp
	XId             *regexp.Regexp
	Tags            []string
	TagsOr          []string
	Output          string
	TableIdWidth    int
	TableTitleWidth int
	TableTagsWidth  int
	Since		*time.Time
	Until		*time.Time
}

func LsCommand(conf AppConfig, args []string) {
	lsConf, args := LsFlags(args)

	files, err := FindAllNotesFiles(conf.RootPath, args)
	if err != nil {
		log.Fatalf("ls: Unable to read notes files %v", err)
	}

	repo, err := git.PlainOpen(conf.RootPath)
	if err != nil {
		log.Fatalf("info: Unable to read git repository from root directory %s", conf.RootPath)
	}

	notePrintFunc := NotePrint(lsConf)
	for _, file := range files {
		notes, err := GetNoteData(repo, file.Absolute, file.Relative)
		if err != nil {
			log.Fatalf("main: Unable to read file %s. %v", file, err)
		}
		for _, note := range notes {
			if !NoteTagsFilter(note, lsConf.Tags, lsConf.TagsOr) {
				continue
			}
			if !lsConf.XId.MatchString(note.Id()) {
				continue
			}
			if !lsConf.XTitle.MatchString(note.Title()) {
				continue
			}

			if until := lsConf.Until; until != nil {
				if updateAt := note.LastUpdate(); updateAt == nil || updateAt.Unix() > until.Unix() {
					continue
				}
			}

			if since := lsConf.Since; since != nil {
				if updateAt := note.LastUpdate(); updateAt == nil || updateAt.Unix() < since.Unix() {
					continue
				}
			}

			var links []NoteId
			if lsConf.Links {
				for _, n := range note.Links() {
					links = append(links, n)
				}
			}
			notePrintFunc(note, links)
		}
	}
}

func LsFlags(args []string) (LsConf, []string) {
	var conf LsConf
	lsf := flag.NewFlagSet("ls", 0)
	lsf.StringVar(&conf.Output, "o", "table", "Output format. [table]")
	lsf.BoolVar(&conf.Links, "links", false, "Retrieve links information")
	lsf.IntVar(&conf.TableIdWidth, "tableIdWidth", 24, "Table width Id")
	lsf.IntVar(&conf.TableTitleWidth, "tableTitleWidth", 60, "Table width Title")
	lsf.IntVar(&conf.TableTagsWidth, "tableTagsWidth", 30, "Table width Tags")
	xId := lsf.String("xid", "", "Regex match id")
	xTitle := lsf.String("xtitle", "", "Regex match title")
	tags := lsf.String("t", "", "Tags AND")
	tagsOr := lsf.String("tor", "", "Tags OR")
	since := lsf.String("since", "", "Since date")
	until := lsf.String("until", "", "Until date")
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

	if *tagsOr != "" {
		conf.TagsOr = strings.Split(*tagsOr, ",")
	}

	if !slices.Contains(validOutputFlags, conf.Output) {
		log.Fatalf("Invalid output flag %s", conf.Output)
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

	if until := *until; until != "" {
		untilT, err := time.Parse(time.RFC3339, until)
		if err != nil {
			log.Fatalf("ls command: invalid until %s. %v", until, err)
		}
		conf.Until =&untilT
	}

	if since := *since; since != "" {
		sinceT, err := time.Parse(time.RFC3339, since)
		if err != nil {
			log.Fatalf("ls command: invalid since %s. %v", since, err)
		}
		conf.Since =&sinceT
	}
	return conf, lsf.Args()
}

func NotePrint(conf LsConf) NotePrintFunc {
	switch conf.Output {
	default:
		return TablePrintNote(conf)
	}
}

func TablePrintNote(conf LsConf) NotePrintFunc {
	headerFmt := fmt.Sprintf("%%-%ds  %%-%ds  %%-%ds  %%s\n",
		conf.TableIdWidth, conf.TableTitleWidth, conf.TableTagsWidth)
	rowFmt := fmt.Sprintf("%%-%ds  %%-%ds  %%-%ds  %%s\n",
		conf.TableIdWidth, conf.TableTitleWidth, conf.TableTagsWidth)

	fmt.Printf(headerFmt, "ID", "TITLE", "TAGS", "PATH")
	fmt.Println(strings.Repeat("-", conf.TableIdWidth+conf.TableTitleWidth+conf.TableTagsWidth+10) + "----------------------------------------")

	return func(n Note, noteIds []NoteId) {
		title := truncate(n.Title(), conf.TableTitleWidth)
		tags := truncate(strings.Join(n.Tags(), ", "), conf.TableTagsWidth)
		fmt.Printf(rowFmt, n.Id(), title, tags, n.OpenRef())
		for i, id := range noteIds {
			prefix := "├─"
			if i == len(noteIds) - 1 {
				prefix = "└─"
			}
			fmt.Printf("%s %s\n", prefix, id)
		}
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
