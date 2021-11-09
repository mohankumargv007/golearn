package main
import ( 
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"errors"
	"alexedwards.net/snippetbox/pkg/models"
	"log"
	"alexedwards.net/snippetbox/pkg/models/mysql"
	"encoding/json"
	"database/sql"
)

type EmployeeModel struct {
	DB *sql.DB
}

var db *sql.DB


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

	data := &templateData{Snippets: s}

	ts, err := template.ParseFiles("./ui/html/home.page.tmpl", "./ui/html/base.layout.tmpl", "./ui/html/footer.partial.tmpl")

	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err) 
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}
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

	/* files := []string{
		"./ui/html/show.page.tmpl", 
		"./ui/html/base.layout.tmpl", 
		"./ui/html/footer.partial.tmpl",
	} */

	data := &templateData{Snippet: s}

	ts, err := template.ParseFiles("./ui/html/show.page.tmpl", "./ui/html/base.layout.tmpl", "./ui/html/footer.partial.tmpl") 
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err) 
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
	//Check POST Method
	if r.Method != http.MethodPost { 
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", 405) 
		return
	}

	//Get URL Params
	query := r.URL.Query()

    if len(query) == 0 {
        fmt.Println("Please send params.")
    }

	//Call below one to insert data
	_, err := app.employees.Insert(query["emp_id"][0], query["emp_name"][0], query["role"][0])
	if err != nil {
		w.WriteHeader(405)
		w.Write([]byte(err.Error()))
		app.serverError(w, err)
		return
	}

	//Success Response
	w.WriteHeader(200)
	w.Write([]byte("Created Employee Successfully."))
	return
}

func (app *application) updateEmployee(w http.ResponseWriter, r *http.Request) {
	//Check POST Method
	if r.Method != http.MethodPost { 
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", 405) 
		return
	}

	//Get URL Params
	query := r.URL.Query()

    if len(query) == 0 {
		//Success Response
		w.WriteHeader(204)
		w.Write([]byte("No content found.Please share employee Id"))
    }

	//Call below one to update data
	_, err := app.employees.Update(query["emp_id"][0], query["emp_name"][0], query["role"][0])
	if err != nil {
		w.WriteHeader(500)
		app.serverError(w, err)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("Updated employee successfully."))
	return
}

func (app *application) showAllEmpList(w http.ResponseWriter, r *http.Request) {
	//Call below one to insert data
	w.Header().Set("Content-Type", "application/json")
	empId := 1
	result, err := app.employees.Show(w, r, empId)
	if err != nil {
		w.WriteHeader(405)
		w.Write([]byte(err.Error()))
		app.serverError(w, err)
		return
	}
	json.NewEncoder(w).Encode(result)
}