package main

import (
	"time"
)

type NoteFileYaml struct {
	Notes []NoteYaml `yaml:"notes"`
}

type Note interface {
	Id() string
	Path() string
	Title() string
	Tags() []string
	CreateAt() time.Time
	UpdateAt() time.Time
}

type NoteYaml struct {
	IdY string `yaml:"id"`
	PathY string
	TitleY string `yaml:"title"`
	TagsY []string `yaml:"tags"`
	CreateAtY time.Time `yaml:"createAt"`
	UpdateAtY time.Time `yaml:"updateAt"`
	NoteY string `yaml:"note"`
}

func (n NoteYaml) Id() string {
	return n.IdY
}

func (n NoteYaml) Path() string {
	return n.PathY
}

func (n NoteYaml) Title() string {
	return n.TitleY
}

func (n NoteYaml) Tags() []string {
	return n.TagsY
}

func (n NoteYaml) CreateAt() time.Time {
	return n.CreateAtY
}

func (n NoteYaml) UpdateAt() time.Time {
	return n.UpdateAtY
}

type NoteMd struct {
	IdM string `yaml:"id"`
	PathM string
	TitleM string `yaml:"title"`
	TagsM []string `yaml:"tags"`
	CreateAtM time.Time `yaml:"createAt"`
	UpdateAtM time.Time `yaml:"updateAt"`
}

func (n NoteMd) Id() string {
	return n.IdM
}

func (n NoteMd) Path() string {
	return n.PathM
}

func (n NoteMd) Title() string {
	return n.TitleM
}

func (n NoteMd) Tags() []string {
	return n.TagsM
}

func (n NoteMd) CreateAt() time.Time {
	return n.CreateAtM
}

func (n NoteMd) UpdateAt() time.Time {
	return n.UpdateAtM
}

