package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"scvis.vaccsedu.org/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	logger         *slog.Logger
	projects       models.ProjectModelInterface
	cfiles         models.CFileModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

var localServer bool

func init() {
	flag.BoolVar(&localServer, "local", false, "Run server in local mode")
}
func main() {
	flag.Parse()

	if localServer {
		err := godotenv.Load("./scvisual.conf")
		if err != nil {
			log.Fatal("Error loading ./scvisual.conf file")
		}
	} else {
		err := godotenv.Load("/etc/scvisual/scvisual.conf")
		if err != nil {
			log.Fatal("Error loading /etc/scvisual/scvisual.conf file")
		}
	}

	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")
	dbName := os.Getenv("DB_NAME")

	dsn := flag.String("dsn", fmt.Sprintf("%s:%s@/%s?parseTime=true", dbUser, dbPasswd, dbName), "MySQL data source name")
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	var logDir string
	if localServer {
		logDir = "./log/scvisual"
	} else {
		logDir = "/var/log/scvisual"
	}
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create log directory: %v", err)
		}
	}

	logFile, err := os.CreateTemp(logDir, "scvisual-*.log")
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()

	logger := slog.New(slog.NewTextHandler(logFile, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		logger:         logger,
		projects:       &models.ProjectModel{DB: db},
		cfiles:         &models.CFileModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12, // Ensure TLS 1.2 or higher
	}

	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logger.Info("getting embedded build folder")

	logger.Info("starting server", "addr", srv.Addr)
	log.Println("Initializing application")

	if localServer {
		err = srv.ListenAndServeTLS("./tls/localhost.pem", "./tls/localhost-key.pem")
		if err != nil {
			log.Println("Error starting server:", err)
		}
	} else {
		err = srv.ListenAndServeTLS("/etc/ssl/private/cs.wmich.edu.chained.crt", "/etc/ssl/private/cs.wmich.edu.key")
	}
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
