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
    Links() []Note
    LastUpdate() *time.Time
    OpenRef() string
    Note() []byte
}

// Define basic data for notes operations
type BaseNote struct {
	id NoteId `yaml:"id"`
	title NoteId `yaml:"title"`
	tags []string `yaml:"tags"`
	links []Note `yaml:"links"`
}

func (n *BaseNote) Id() NoteId {
	return n.id
}

func (n *BaseNote) Title() NoteId {
	return n.title
}

func (n *BaseNote) Tags() []string {
	return n.tags
}

func (n *BaseNote) Links() []Note {
	return n.links
}
