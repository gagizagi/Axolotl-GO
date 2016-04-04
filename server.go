package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/eknkc/amber"
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
	http.HandleFunc("/anime", animeHandler)
	http.HandleFunc("/static/", staticHandler)

	//Sets url and port
	bind := fmt.Sprintf("%s:%s", "127.0.0.1", "422")
	if os.Getenv("OPENSHIFT_GO_IP") != "" &&
		os.Getenv("OPENSHIFT_GO_PORT") != "" {
		bind = fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"),
			os.Getenv("OPENSHIFT_GO_PORT"))
	}

	//Listen and sert to port
	log.Printf("Web server listening on %s", bind)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		log.Fatal("webServer() => ListenAndServer() error:\t", err)
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
