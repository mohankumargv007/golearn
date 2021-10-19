package main
import ( 
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"errors"
	"alexedwards.net/snippetbox/pkg/models"
	"log"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *mysql.SnippetModel
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	ts, err := template.ParseFiles("./ui/html/home.page.tmpl", "./ui/html/base.layout.tmpl", "./ui/html/footer.partial.tmpl")

	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err) 
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	print("oooo");
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w) 
		} else {
			app.serverError(w, err) 
		}
		return
	}
	// Write the snippet data as a plain-text HTTP response body.
	fmt.Fprintf(w, "%v", s)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { 
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", 405) 
		return
	}
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"

	///Call below one to insert data
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Write([]byte("Create a new snippet..."))

	//Redirecting to created snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther) 
}