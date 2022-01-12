package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"alexedwards.net/snippetbox/pkg/forms"
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
	users         *mysql.UserModel
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	//s, err := app.snippets.Latest()

	s, err := app.employees.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	//data := &templateData{Snippets: s}
	data := &templateData{Employees: s}

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

func (app *application) employeeDashboard(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "empdashboard.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
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

	flash := app.session.PopString(r, "flash")

	app.render(w, r, "show.page.tmpl", &templateData{
		Flash:   flash,
		Snippet: s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	w.Write([]byte(title))

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

func (app *application) createEmp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	emp_id := r.PostForm.Get("emp_id")
	emp_name := r.PostForm.Get("emp_name")
	role := r.PostForm.Get("role")

	w.Write([]byte(emp_id))

	///Call below one to insert data
	id, err := app.employees.Insert(emp_id, emp_name, role)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Write([]byte("Create a new emp..."))

	app.session.Put(r, "flash", "Emp successfully created!")

	//Redirecting to created snippet
	http.Redirect(w, r, fmt.Sprintf("/employee?id=%d", id), http.StatusSeeOther)
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
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Validate the form contents using the form helper we made earlier.
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.session.Put(r, "authenticatedUserID", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
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
