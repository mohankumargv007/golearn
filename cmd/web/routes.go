package main

import ( 
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	//Route Initialization
	mux := http.NewServeMux()

	//All Required Routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	//Static Folder Path
	fileServer := http.FileServer(http.Dir("./ui/static/")) 
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux 
}