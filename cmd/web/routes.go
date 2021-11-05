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

	//Employee API's

	//Got Issues with migration so created a api to create a table now
	mux.HandleFunc("/employees/createTable", app.createEmpTable)
	
	mux.HandleFunc("/employees/show", app.showAllEmpList)
	mux.HandleFunc("/employees/create", app.createEmp)
	mux.HandleFunc("/employees/updateEmp", app.updateEmployee)

	//Static Folder Path
	fileServer := http.FileServer(http.Dir("./ui/static/")) 
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux 
}