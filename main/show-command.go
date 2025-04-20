package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
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
		notes, err := GetNoteData(repo, f.Absolute, f.Relative)
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
	switch n := n.(type) {
	case *NoteMarkdown:
		absoluteFilePath := n.AbsoluteFilePath
		f, err := os.Open(absoluteFilePath)
		if err != nil {
			return fmt.Errorf("show: Unable to open md file. %v", err)
		}
		defer f.Close()
		_, err = io.Copy(os.Stdout, f)
		if err != nil {
			return fmt.Errorf("show: Unable to stream file to stdout. %v", err)
		}
		return nil
	case *NoteYaml:
		absolutePath := n.absoluteFilePath
		size := n.size
		offset := n.offset
		f, err := os.Open(absolutePath)
		if err != nil {
			return fmt.Errorf("show: unable to open file %s", absolutePath)
		}
		if _, err := io.CopyN(io.Discard, f, int64(offset)); err != nil {
			return fmt.Errorf("show: unable to read offset, content=%d, file=%s", offset, absolutePath)
		}
		if _, err := io.CopyN(os.Stdout, f, int64(size)); err != nil {
			return fmt.Errorf("show: unable to read note content. offset=%d size=%d file=%s", offset, size, absolutePath)
		}
		return nil
	default:
		return fmt.Errorf("show: unsupported note %T", n)
		
	}
}
