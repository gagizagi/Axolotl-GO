package main

import (
	"html/template"
	"net/http"
	"fmt"
	"github.com/eknkc/amber"
)

var templateMap map[string]*template.Template

func init() {
	DefaultOptions := amber.Options{true, true}
	DefaultDirOptions := amber.DirOptions{".amber", true}
	templateMap, _ = amber.CompileDir("views", DefaultDirOptions, DefaultOptions)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "/home")
}

func AnimeHandler(w http.ResponseWriter, r *http.Request) {
	templateMap["aList"].Execute(w, Get_anime_list())
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}