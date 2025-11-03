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
	"strconv"
	"time"

	"scvis.vaccsedu.org/internal/github"
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
	inputfiles     models.InputFileModelInterface
	users          models.UserModelInterface
	issues         models.IssueModelInterface
	githubClient   *github.GitHubClient
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
	// First, try to load /opt/scvisual/.env, fall back to local .env if that fails
	var envPath string
	var err error

	// Try production path first
	envPath = "/opt/scvisual/.env"
	err = godotenv.Load(envPath)

	if err != nil {
		// Fall back to local .env file
		envPath = "./.env"
		err = godotenv.Load(envPath)
		if err != nil {
			log.Fatalf("Error loading environment file from both /opt/scvisual/.env and ./.env: %v", err)
		}
	}

	// Now determine server mode from LOCAL_SERVER in the loaded .env file
	localServerEnv := os.Getenv("LOCAL_SERVER")
	if localServerEnv != "" {
		// Explicit configuration via environment variable from .env file
		localServer = localServerEnv == "true"
	} else {
		// Fallback: assume local server if using .env, production if using /opt/scvisual/.env
		localServer = envPath == "./.env"
	}

	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@/%s?parseTime=true", dbUser, dbPasswd, dbName)
	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))

	// Get log directory from environment variable, with fallback logic
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		if localServer {
			logDir = "./log/scvisual"
		} else {
			logDir = "/var/log/scvisual"
		}
	}

	logPath := filepath.Join(logDir, "scvisual.log")

	// Get log rotation settings from environment variables, with defaults
	maxSizeMB := 10 // Default 10MB
	maxBackups := 5 // Default 5 backup files

	// Read from environment variables if set
	if envMaxSize := os.Getenv("LOG_MAX_SIZE_MB"); envMaxSize != "" {
		if size, err := strconv.Atoi(envMaxSize); err == nil {
			maxSizeMB = size
		} else {
			logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
			logger.Warn("Invalid LOG_MAX_SIZE_MB value, using default", "value", envMaxSize, "default", maxSizeMB)
		}
	}

	if envMaxBackups := os.Getenv("LOG_MAX_BACKUPS"); envMaxBackups != "" {
		if backups, err := strconv.Atoi(envMaxBackups); err == nil {
			maxBackups = backups
		} else {
			logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
			logger.Warn("Invalid LOG_MAX_BACKUPS value, using default", "value", envMaxBackups, "default", maxBackups)
		}
	}

	// For development/testing, use configured values or smaller defaults
	if localServer && os.Getenv("LOG_MAX_SIZE_MB") == "" {
		maxSizeMB = 1 // 1MB for local development if not explicitly set
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

	// Initialize GitHub client
	githubOwner := os.Getenv("GITHUB_OWNER")
	githubRepo := os.Getenv("GITHUB_REPO")
	if githubOwner == "" {
		githubOwner = "vaccs" // Default organization
	}
	if githubRepo == "" {
		githubRepo = "scvis-go" // Default repository
	}
	githubClient := github.NewGitHubClient(githubOwner, githubRepo)

	app := &application{
		logger:         logger,
		projects:       &models.ProjectModel{DB: db},
		cfiles:         &models.CFileModel{DB: db},
		inputfiles:     &models.InputFileModel{DB: db},
		users:          &models.UserModel{DB: db},
		issues:         &models.IssueModel{DB: db},
		githubClient:   githubClient,
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

	// Get SSL certificate paths from environment variables, with fallback logic
	var certPath, keyPath string

	if envCertPath := os.Getenv("SSL_CERT_PATH"); envCertPath != "" {
		certPath = envCertPath
	} else if localServer {
		certPath = "./tls/localhost.pem"
	} else {
		certPath = "/etc/letsencrypt/live/vaccs.ceas.wmich.edu/fullchain.pem"
	}

	if envKeyPath := os.Getenv("SSL_KEY_PATH"); envKeyPath != "" {
		keyPath = envKeyPath
	} else if localServer {
		keyPath = "./tls/localhost-key.pem"
	} else {
		keyPath = "/etc/letsencrypt/live/vaccs.ceas.wmich.edu/privkey.pem"
	}

	logger.Info("using SSL certificates", "cert", certPath, "key", keyPath)
	err = srv.ListenAndServeTLS(certPath, keyPath)

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
