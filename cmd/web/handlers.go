package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"alexedwards.net/snippetbox/pkg/models"
	"alexedwards.net/snippetbox/pkg/models/mysql"

	//"encoding/json"
	"database/sql"

	"os"

	"github.com/golangcollege/sessions"
)

type EmployeeModel struct {
	DB *sql.DB
}

var db *sql.DB

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	employees     *mysql.EmployeeModel
	session       *sessions.Session
	templateCache map[string]*template.Template
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

	app.session.Put(r, "flash", "Snippet successfully created!")

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

/* func (app *application) showAllEmpList(w http.ResponseWriter, r *http.Request) {
	//Call below one to show data
	w.Header().Set("Content-Type", "application/json")
	result, err := app.employees.Show()
	if err != nil {
		w.WriteHeader(405)
		w.Write([]byte(err.Error()))
		app.serverError(w, err)
		return
	}
	json.NewEncoder(w).Encode(result)
}
*/
//New Functions For User

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display the user signup form...")
}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new user...")
}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display the user login form...")
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}

///File Operations

//1.File Check
func (app *application) fileCheck(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat("ui/static/csv/cities.txt"); err == nil {
		w.Write([]byte("File exists."))
	} else {
		w.Write([]byte("File not exists."))
	}
}

//2.Read File
func (app *application) readFile(w http.ResponseWriter, r *http.Request) {
	file, err := os.Create("ui/static/csv/cities.txt")
	if err != nil {
		w.Write([]byte("File Creation Failed.Please Try After Some Time."))
		return
	}
	defer file.Close()

	var cities = []string{
		"Hyd",
		"Bangalore",
		"Mumbai",
	}
	for i := 0; i < len(cities); i++ {
		file.WriteString(cities[i])
		file.WriteString("\n")
	}
	w.Write([]byte("Cities File Created Successfully."))
}

//3.Renaming the file
func (app *application) renameFile(w http.ResponseWriter, r *http.Request) {
	currentFile := "ui/static/csv/cities.txt"
	newName := "ui/static/csv/new_cities.txt"
	e := os.Rename(currentFile, newName)
	if e != nil {
		fmt.Fprintln(w, e)
		return
	}
	w.Write([]byte("File Renamed Successfully."))
}
