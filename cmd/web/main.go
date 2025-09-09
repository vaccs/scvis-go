package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
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

// LogRotator manages log file rotation
type LogRotator struct {
	logFile    *os.File
	logPath    string
	maxSize    int64 // Maximum size in bytes before rotation
	maxBackups int   // Maximum number of backup files to keep
}

// NewLogRotator creates a new log rotator
func NewLogRotator(logPath string, maxSize int64, maxBackups int) (*LogRotator, error) {
	lr := &LogRotator{
		logPath:    logPath,
		maxSize:    maxSize,
		maxBackups: maxBackups,
	}

	err := lr.openLogFile()
	if err != nil {
		return nil, err
	}

	return lr, nil
}

// openLogFile opens or creates the log file
func (lr *LogRotator) openLogFile() error {
	// Ensure directory exists
	dir := filepath.Dir(lr.logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(lr.logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	lr.logFile = file
	return nil
}

// Write implements io.Writer interface with rotation logic
func (lr *LogRotator) Write(p []byte) (n int, err error) {
	// Check if rotation is needed
	if lr.needsRotation() {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	return lr.logFile.Write(p)
}

// needsRotation checks if the log file needs to be rotated
func (lr *LogRotator) needsRotation() bool {
	if lr.logFile == nil {
		return false
	}

	info, err := lr.logFile.Stat()
	if err != nil {
		return false
	}

	return info.Size() >= lr.maxSize
}

// rotate rotates the log file
func (lr *LogRotator) rotate() error {
	// Close current log file
	if lr.logFile != nil {
		lr.logFile.Close()
	}

	// Rotate existing backup files
	for i := lr.maxBackups; i > 0; i-- {
		oldPath := fmt.Sprintf("%s.%d", lr.logPath, i)
		newPath := fmt.Sprintf("%s.%d", lr.logPath, i+1)

		if i == lr.maxBackups {
			// Remove the oldest backup
			os.Remove(oldPath)
		} else {
			// Move backup file to next number
			if _, err := os.Stat(oldPath); err == nil {
				os.Rename(oldPath, newPath)
			}
		}
	}

	// Move current log to .1 backup
	if _, err := os.Stat(lr.logPath); err == nil {
		os.Rename(lr.logPath, lr.logPath+".1")
	}

	// Create new log file
	return lr.openLogFile()
}

// Close closes the log file
func (lr *LogRotator) Close() error {
	if lr.logFile != nil {
		return lr.logFile.Close()
	}
	return nil
}

func main() {
	// Check if /etc/scvisual/scvisual.conf is accessible to determine server mode
	if file, err := os.Open("/etc/scvisual/scvisual.conf"); err == nil {
		file.Close()
		localServer = false
	} else {
		localServer = true
	}

	if localServer {
		err := godotenv.Load("./.env")
		if err != nil {
			log.Fatal("Error loading .env file")
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

	dsn := fmt.Sprintf("%s:%s@/%s?parseTime=true", dbUser, dbPasswd, dbName)
	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))

	var logDir string
	if localServer {
		logDir = "./log/scvisual"
	} else {
		logDir = "/var/log/scvisual"
	}

	logPath := filepath.Join(logDir, "scvisual.log")

	// Get log rotation settings from environment, with defaults
	maxSizeMB := 10 // Default 10MB
	maxBackups := 5 // Default 5 backup files

	// For development/testing, allow smaller log files
	if localServer {
		maxSizeMB = 1 // 1MB for local development
	}

	// Create log rotator
	logRotator, err := NewLogRotator(logPath, int64(maxSizeMB)*1024*1024, maxBackups)
	if err != nil {
		log.Fatalf("Failed to create log rotator: %v", err)
	}
	defer logRotator.Close()

	logger := slog.New(slog.NewTextHandler(logRotator, nil))

	db, err := openDB(dsn)
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
		Addr:         addr,
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
	} else {
		err = srv.ListenAndServeTLS("/etc/letsencrypt/live/vaccs.ceas.wmich.edu/fullchain.pem", "/etc/letsencrypt/live/vaccs.ceas.wmich.edu/privkey.pem")
	}

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
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
