package main

import (
	"net/http"

	"github.com/bmizerany/pat"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	//Route Initialization
	mux := pat.New()
	//standardMiddleware := alice.New(app.recoverPanic, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)

	//All Required Routes
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	//Employee API's

	//Got Issues with migration so created a api to create a table now
	//	mux.HandleFunc("/employees/createTable", app.createEmpTable)

	//mux.HandleFunc("/employees/show", app.showAllEmpList)
	//	mux.HandleFunc("/employees/create", app.createEmp)
	//	mux.HandleFunc("/employees/updateEmp", app.updateEmployee)

	// Add the five new routes.
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))

	//File Operation APIs
	mux.Get("/file/check", dynamicMiddleware.ThenFunc(app.fileCheck))
	mux.Get("/file/read", dynamicMiddleware.ThenFunc(app.readFile))
	mux.Get("/file/rename", dynamicMiddleware.ThenFunc(app.renameFile))

	//Static Folder Path
	//fileServer := http.FileServer(http.Dir("./ui/static/"))
	//	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
	//return standardMiddleware(mux)
}
