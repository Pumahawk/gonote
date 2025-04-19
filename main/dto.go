package main

import (
	"time"

	"github.com/go-git/go-git/v5"
)

type NoteId = string

// Note main directory
// Must contain git repository
type Repository struct {
	Path string
	Git *git.Repository
}

// Define basic note operations
type Note interface {
    Id() NoteId
    Title() string
    Tags() []string
    Links() []NoteId
    LastUpdate() *time.Time
    OpenRef() string
}

// Define basic data for notes operations
type BaseNote struct {
	IdV NoteId `yaml:"id"`
	TitleV string `yaml:"title"`
	TagsV []string `yaml:"tags"`
	LinksV []NoteId `yaml:"links"`
}

func (n *BaseNote) Id() NoteId {
	return n.IdV
}

func (n *BaseNote) Title() NoteId {
	return n.TitleV
}

func (n *BaseNote) Tags() []string {
	return n.TagsV
}

func (n *BaseNote) Links() []NoteId {
	return n.LinksV
}
