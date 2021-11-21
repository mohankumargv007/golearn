package main

import ( 
	"net/http"
	//"github.com/justinas/alice"
)

func (app *application) routes() *http.ServeMux {
	//Route Initialization
	mux := http.NewServeMux()
	//dynamicMiddleware := alice.New(app.session.Enable)
	
	//All Required Routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	//Employee API's

	//Got Issues with migration so created a api to create a table now
	mux.HandleFunc("/employees/createTable", app.createEmpTable)
	
	//mux.HandleFunc("/employees/show", app.showAllEmpList)
	mux.HandleFunc("/employees/create", app.createEmp)
	mux.HandleFunc("/employees/updateEmp", app.updateEmployee)

	// Add the five new routes.
	/* mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm)) 
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser)) 
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm)) 
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser)) 
	mux.Post("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser)) */

	//Static Folder Path
	fileServer := http.FileServer(http.Dir("./ui/static/")) 
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
	//return secureHeaders.Then(mux)
}