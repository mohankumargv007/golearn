package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

//Server Errors
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) 
}

//Specific Client Errors
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status) 
}

//404 Errors
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound) 
}
	