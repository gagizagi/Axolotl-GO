package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/eknkc/amber"
	"google.golang.org/appengine"
)

//template map
var templateMap map[string]*template.Template

//initialization of the template map
func init() {
	templateMap, _ = amber.CompileDir("views",
		amber.DirOptions{Ext: ".amber", Recursive: true},
		amber.Options{PrettyPrint: false, LineNumbers: false})
}

//complete web server function
func webServer() {
	//Sets path handler funcions
	http.HandleFunc("/", animeHandler)
	http.HandleFunc("/static/", staticHandler)

	log.Println("Starting http server...")

	if os.Getenv("AX_ENV") == "production" {
		appengine.Main()
	} else {
		addr := ":" + os.Getenv("AX_DEVELOPMENT_PORT")

		http.ListenAndServe(addr, nil)
	}
}

// /anime path handler
func animeHandler(w http.ResponseWriter, r *http.Request) {
	data := getAnimeList()
	sort.Sort(data)
	templateMap["aList"].Execute(w, data)
}

// /static/* file server
func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}
