package main

import (
	"fmt"
	"os"
	"net/http"
	"sort"
	"html/template"
	"github.com/eknkc/amber"
)

var templateMap map[string]*template.Template

func init() {
	DefaultOptions := amber.Options{false, false}
	DefaultDirOptions := amber.DirOptions{".amber", true}
	templateMap, _ = amber.CompileDir("views", DefaultDirOptions, DefaultOptions)
}

func Web_server() {
	http.HandleFunc("/anime", AnimeHandler)
	http.HandleFunc("/static/", StaticHandler)
	
	bind := fmt.Sprintf("%s:%s", "127.0.0.1", "422")
	if os.Getenv("OPENSHIFT_GO_IP") != "" && os.Getenv("OPENSHIFT_GO_PORT") != ""{
		bind = fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	}
	fmt.Printf("Web server listening on %s\n", bind)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		panic(err)
	}
}

func AnimeHandler(w http.ResponseWriter, r *http.Request) {
	data := Get_anime_list()
	sort.Sort(data)
	templateMap["aList"].Execute(w, data)
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}
