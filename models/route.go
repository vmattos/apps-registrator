package models

type Route struct {
	Path    string
	Backend string
}

type Backend struct {
	Type string
}
