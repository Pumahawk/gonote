package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"time"
)

type NoteId = string

type NoteLink struct {
	Id NoteId
	Title string
	Path string
	Line int
}

type Note interface {
	Line() int
	Id() NoteId
	Path() string
	Links() []NoteId
	Title() string
	Tags() []string
	UpdateAt() *time.Time
}

type NoteYaml struct {
	IdY       NoteId `yaml:"id"`
	pathY     string
	lineY     int
	lineEndY     int
	TitleY    string    `yaml:"title"`
	TagsY     []string  `yaml:"tags"`
	UpdateAtY *time.Time `yaml:"updateAt"`
	NoteY     string    `yaml:"note"`
}

func (n NoteYaml) Line() int {
	return n.lineY
}

func (n NoteYaml) Id() NoteId {
	return n.IdY
}

func (n NoteYaml) Path() string {
	return n.pathY
}

func (n NoteYaml) Links() []NoteId {
	r := bytes.NewReader([]byte(n.NoteY))
	links, err := GetLinksFromText(r)
	if err != nil {
		log.Printf("noteyaml links: Unable to retrieve links from note %s. %v", n.Id(), err)
		return nil
	}
	return links
}

func (n NoteYaml) Title() string {
	return n.TitleY
}

func (n NoteYaml) Tags() []string {
	return n.TagsY
}

func (n NoteYaml) UpdateAt() *time.Time {
	return n.UpdateAtY
}

type NoteMd struct {
	IdM       string `yaml:"id"`
	PathM     string
	TitleM    string    `yaml:"title"`
	TagsM     []string  `yaml:"tags"`
	UpdateAtM *time.Time `yaml:"updateAt"`
}

func (n NoteMd) Line() int {
	return 1
}

func (n NoteMd) Id() string {
	return n.IdM
}

func (n NoteMd) Path() string {
	return n.PathM
}

func (n NoteMd) Links() []NoteId {
	f, err := os.Open(n.Path())
	if err != nil {
		log.Printf("notemd: Unable to retrieve links from md noted, %v", err)
		return nil
	}
	defer f.Close()
	// Consume stream
	r := bufio.NewReader(f)
	MarkdownNote(r)
	links, err := GetLinksFromText(r)
	if err != nil {
		log.Printf("notemd links: Unable to retrieve links from note %s. %v", n.Id(), err)
		return nil
	}
	return links
}

func (n NoteMd) Title() string {
	return n.TitleM
}

func (n NoteMd) Tags() []string {
	return n.TagsM
}

func (n NoteMd) UpdateAt() *time.Time {
	return nil
}
