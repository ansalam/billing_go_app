package main

import (
	"billing_api/pkg/models/sqlite"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var infoFile *os.File

type application struct {
	infoLog         *log.Logger
	DBConn          *sqlite.CountersModel
	authenticatorID string
}

func initLoggers() {
	var err error
	infoFile, err = os.OpenFile("./info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Runtime configuration parser
	addr := flag.String("addr", ":4000", "HTTP network address")
	// dsn := flag.String("dsn", "web:Hellodude@13@/snippetbox?parseTime=true", "DB connection")
	dsn := flag.String("dsn", "./dummy.db", "Sqlite3 DB")
	flag.Parse()

	// Logger Initialization
	initLoggers()
	infoLog := log.New(infoFile, "INFO\t", log.Ldate|log.Ltime)
	defer infoFile.Close()
	infoLog.Println("Loggers initialized.")

	// DB Initialization
	db, err := openDB(*dsn)
	if err != nil {
		infoLog.Fatal(err)
	}
	defer db.Close()
	infoLog.Println("DB initialized.")

	app := &application{
		infoLog:         infoLog,
		DBConn:          &sqlite.CountersModel{DB: db},
		authenticatorID: "authenticator-1", // TODO: create a unique-id dynamically.
	}

	_, err = app.DBConn.CreateTable()
	if err != nil {
		infoLog.Fatal(err)
	}

	// Server Initialization
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: infoLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on %s", srv.Addr)
	err = srv.ListenAndServe()
	infoLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
