package models

type Tree struct {
	Title    string
	Children []Tree
}

func (t Tree) HasChild() bool {
	return len(t.Children) > 0
}

func (t Tree) IsRoot() bool {
	return t.Title == ""
}
