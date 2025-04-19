package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	yaml "github.com/goccy/go-yaml"
)

type ShowConf struct {
}

func ShowCommand(conf AppConfig, args []string) {
	ShowFlags(args)

	if len(args) < 1 {
		log.Fatalf("Need note id parameter")
	}

	noteId := args[0]

	files, err := FindAllNotesFiles(conf.RootPath, []string{})
	if err != nil {
		log.Fatalf("show: Unable to read notes files %v", err)
	}

	repo, err := git.PlainOpen(conf.RootPath)
	if err != nil {
		log.Fatalf("info: Unable to read git repository from root directory %s", conf.RootPath)
	}

	for _, f := range files {
		notes, err := GetNoteData(repo, f)
		if err != nil {
			log.Fatalf("show: Unable to read file %s. %v", f, err)
		}
		for _, n := range notes {
			if n.Id() == noteId {
				if err := CatNote(n); err != nil {
					log.Fatalf("show: Unable to stream markdown file. %v", err)
				}
			}
		}
	}
}

func ShowFlags(args []string) ShowConf {
	var conf ShowConf
	lsf := flag.NewFlagSet("show", 0)
	err := lsf.Parse(args)
	if err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
	return conf
}

func CatNote(n Note) error {
	switch n.(type) {
	case NoteMd:
		f, err := os.Open(n.Path())
		if err != nil {
			return fmt.Errorf("show: Unable to open md file. %v", err)
		}
		defer f.Close()
		_, err = io.Copy(os.Stdout, f)
		if err != nil {
			return fmt.Errorf("show: Unable to stream file to stdout. %v", err)
		}
		return nil
	default:
		yaml.NewEncoder(os.Stdout, yaml.UseLiteralStyleIfMultiline(true)).Encode(n)
		return nil
	}
}
