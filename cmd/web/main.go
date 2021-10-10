package main
import (
	"flag" 
	"log"
	"net/http"
)
func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	//It will check the incoming request and give the appropriate file to user
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//Creating Route For Static Folder
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Starting server on :%s", *addr)
	err := http.ListenAndServe(*addr, mux) 
	log.Fatal(err)
}