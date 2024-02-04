package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"crypto/tls"

	"snippetbox.cvclon3.net/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)


type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder *form.Decoder
	sessionManager *scs.SessionManager
}


func main() {
	// LOGGERS
	infoLog := log.New(os.Stdout, "INFO\n", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\n", log.Ldate|log.Ltime|log.Lshortfile)

	// ENVIRONMENT VARS
	err := godotenv.Load(".env")
	if err != nil {
		errorLog.Fatal(err)
	}

	// FLAGS
	addr := flag.String("addr", ":4040", "HTTP network address")
	dsn := flag.String(
		"dsn", 
		fmt.Sprintf("%s:%s@/snippetbox?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASS")), 
		"MySQL data source name",
	)
	flag.Parse()

	// DATABASE
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	
	defer db.Close()

	// TEMPLATE CAHCE
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// FORM DECODER
	formDecoder := form.NewDecoder()

	// SESSION MANAGER
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// APPLICATION
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder: formDecoder,
		sessionManager: sessionManager,
	}

	// TLS CONFIG
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// CUSTOM HTTP SERVER
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
		TLSConfig: tlsConfig,
		// Timeouts
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}


	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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


// CUSTOM FILESYSTEM
// https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
type neuteredFileSystem struct {
    fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
    f, err := nfs.fs.Open(path)
    if err != nil {
        return nil, err
    }

    s, err := f.Stat()
    if s.IsDir() {
        index := filepath.Join(path, "index.html")
        if _, err := nfs.fs.Open(index); err != nil {
            closeErr := f.Close()
            if closeErr != nil {
                return nil, closeErr
            }

            return nil, err
        }
    }

    return f, nil
}    
// CUSTOM FILESYSTEM END
