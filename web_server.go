package main

import (
	"fmt"
	"os"
	"net/http"
)

func Web_server() {
	http.HandleFunc("/", HomeHandler)
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
