package main
import ( "log"
	"net/http"
)
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	//It will check the incoming request and give the appropriate file to user
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//Creating Route For Static Folder
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux) 
	log.Fatal(err)
}