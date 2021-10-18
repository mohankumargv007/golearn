package main
import (
	"flag" 
	"log"
	"net/http"
	"os"
	"database/sql"
	"alexedwards.net/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *mysql.SnippetModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	//Logging the Info or error
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//MySql Initialize
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	//Connection
	db, err := openDB(*dsn)

	//DB Connection Error
	if err != nil { 
		errorLog.Fatal(err)
	}

	infoLog.Printf("MySql Connected")

	//Db Connection Close
	defer db.Close()

	//App Intialization
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}
	
	//Routes
	//mux := http.NewServeMux()
	//mux.HandleFunc("/", app.home)
	//mux.HandleFunc("/snippet", app.showSnippet)
	//mux.HandleFunc("/snippet/create", app.createSnippet)

	//It will check the incoming request and give the appropriate file to user

	srv := &http.Server{ 
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

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