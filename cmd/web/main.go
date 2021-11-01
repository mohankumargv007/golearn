package main
import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"alexedwards.net/snippetbox/pkg/models/mysql" // New import 
	_ "github.com/go-sql-driver/mysql"	
)

func main() {
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	
	//Open Connection
	db, err := openDB(*dsn) 
	
	if err != nil {
		errorLog.Fatal(err) 
	}

	//Close Connection
	defer db.Close()
	
	//Snippet Model initialize
	app := &application {
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &mysql.SnippetModel{DB: db},
		employees: &mysql.EmployeeModel{DB: db},
	}
	
	//Routes
	srv := &http.Server{ 
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(), 
	}

	//Listen Serve
	infoLog.Printf("Starting server on %s", *addr) 
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
	

func openDB(dsn string) (*sql.DB, error) { 
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err 
	}
	if err = db.Ping(); err != nil { 
		return nil, err
	}
	return db, nil 
}