package models

type Route struct {
	Name       string
	Path       string
	ServiceApp string
	Backend    string
}

type Backend struct {
	Type string
}

type Server struct {
	URL string
}

type Frontend struct {
	Type      string
	BackendId string
	Route     string
}
