package main
import ( 
	"fmt"
	//"html/template"
	"net/http"
	"strconv"
	"errors"
	"alexedwards.net/snippetbox/pkg/models"
	"log"
	"alexedwards.net/snippetbox/pkg/models/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *mysql.SnippetModel
	employees *mysql.EmployeeModel
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Latest() 
	if err != nil {
		app.serverError(w, err)
		return
	}
		
	for _, snippet := range s {
		fmt.Fprintf(w, "%v\n", snippet)
	}	

	//ts, err := template.ParseFiles("./ui/html/home.page.tmpl", "./ui/html/base.layout.tmpl", "./ui/html/footer.partial.tmpl")

	/* if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err) 
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	} */
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
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

//Employee Methods
func (app *application) createEmpTable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { 
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", 405) 
		return
	}

	//To create a new employee table
	id, err := app.employees.CreateTable()
	if err != nil {
		app.serverError(w, err)
		return
	}
	print(id)
	return
}

func (app *application) createEmp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { 
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", 405) 
		return
	}
	
	err := r.ParseForm() 
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//Working on it
	print(r.PostForm["emp_id"])
	return

	//Call below one to insert data
	id, err := app.employees.Insert(empID, empName, role)
	if err != nil {
		app.serverError(w, err)
		return
	}
	print(id);
	w.WriteHeader(200)
	w.Write([]byte("Created Employee Successfully."))
	return
}

func (app *application) showAllEmpList(w http.ResponseWriter, r *http.Request) {
	//Call below one to insert data
	result, err := app.employees.show()
	if err != nil {
		w.WriteHeader(500)
		app.serverError(w, err)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(result))
	return
}