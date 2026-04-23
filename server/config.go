package main

import "github.com/go-chi/chi/v5"

type serverConfig struct {
	DatabaseLocation string `yaml:"database_location"`
	SkipIndexing     bool   `yaml:"skip_indexing"`
	HTTPUser         string `yaml:"http_user"` // probably switch out authentication middleware
	HTTPPass         string `yaml:"http_pass"`
}

type server struct {
	server *chi.Mux
}

type Service interface{}
