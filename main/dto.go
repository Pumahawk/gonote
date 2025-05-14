package main

type NoteId = string

type Repository struct {
	Path string
}

// Define basic note operations
type Note interface {
    Id() NoteId
    Title() string
    Tags() []string
    Links() []NoteId
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
