package main

import "time"

type NoteFile struct {
	Notes []Note `yaml:"notes"`
}
type Note struct {
	Id string `yaml:"id"`
	Path string
	Title string `yaml:"title"`
	Tag []string `yaml:"tag"`
	CreateAt time.Time `yaml:"createAt"`
	UpdateAt time.Time `yaml:"updateAt"`
	Meta string `yaml:"meta"`
}
