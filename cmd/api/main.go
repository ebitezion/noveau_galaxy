package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/ebitezion/backend-framework/internal/accounts"
	"github.com/ebitezion/backend-framework/internal/appauth"
	"github.com/ebitezion/backend-framework/internal/configuration"
	"github.com/ebitezion/backend-framework/internal/data"
	"github.com/ebitezion/backend-framework/internal/payments"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"

	//_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

// Application version number. Later on this will be generated at build time.
const version = "1.0.0"

// Define a config struct to hold all the configuration settings for the app.
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

// Define an application struct to hold the dependencies for HTTP handlers,
// helpers, and middleware.
type application struct {
	config    config
	logger    *log.Logger
	models    data.Models
	templates map[string]*template.Template
	mu        sync.Mutex
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSIONSTORE")))

var once sync.Once

var Log zerolog.Logger

func main() {

	// Define a file server to serve static files
	fs := http.FileServer(http.Dir("cmd/web/static"))

	// Create a route to serve static files under "/static/"
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	err := godotenv.Load()
	if err != nil {
		log.Println("Environment Loading Error", err)
	}

	//Variables, environment and stuffs
	var cfg config
	PORT, err := strconv.ParseInt(os.Getenv("PORT"), 10, 32)
	if err != nil {
		log.Println("String to Int Conversion Error", err)
	}
	maxOpenConns, err := strconv.ParseInt(os.Getenv("DATABASE_MAX_OPEN_CONNS"), 10, 32)
	if err != nil {
		log.Println("String to Int Conversion Error", err)
	}

	maxIdleConns, err := strconv.ParseInt(os.Getenv("DATABASE_MAX_IDLE_CONNS"), 10, 32)
	if err != nil {
		log.Println("String to Int Conversion Error", err)
	}

	maxIdleTime, err := strconv.ParseInt(os.Getenv("DATABASE_MAX_IDLE_TIME"), 10, 32)
	if err != nil {
		log.Println("String to Int Conversion Error", err)
	}

	// Read the values of the command-line flags into the config struct.
	flag.IntVar(&cfg.port, "port", int(PORT), "API server port")
	flag.StringVar(&cfg.env, "env", os.Getenv("ENV"), "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DATABASE_DSN"), "MySQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", int(maxOpenConns), "MySQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", int(maxIdleConns), "MySQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", fmt.Sprintf("%d%v", maxIdleTime, "m"), "MySQL max connection idle time")
	flag.Parse()

	// Initialize a new logger which writes messages to the standard output stream,
	// prefixed with the current date and time.

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	con, err := configuration.LoadConfig()
	if err != nil {
		fmt.Println(err)

	}
	appauth.SetConfig(&con)
	payments.SetConfig(&con)
	accounts.SetConfig(&con)

	// Call the openDB() helper function (see below) to create the connection pool,
	// passing in the config struct. If this returns an error, we log it and exit the
	// application immediately.
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer db.Close()

	// Also log a message to say that the connection pool has been successfully
	// established.
	logger.Printf("database connection pool established")

	// Declare an instance of the application struct, containing the config struct and
	// the logger.
	app := &application{
		config:    cfg,
		logger:    logger,
		models:    data.NewModels(con.Db),
		templates: make(map[string]*template.Template),
	}

	// Declare a HTTP server with some sensible timeout settings, which listens on the
	// port provided in the config struct and uses the servemux we created as the handler.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// // Start the HTTP server.
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
	// http.ListenAndServe(":5050", nil)
}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.
	db, err := sql.Open("mysql", cfg.db.dsn)
	if err != nil {
		fmt.Println("Connection to DB cannot be established...")
		return nil, err
	}

	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	// Set the maximum number of idle connections in the pool. Again, passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	// Use the time.ParseDuration() function to conver the idle timeout duration string
	// to a time.Duration type.
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		fmt.Println("Connection to DB failed due to parse duration...")
		return nil, err
	}

	// Set the maximum idle timeout.
	db.SetConnMaxIdleTime(duration)

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// Return the sql.DB connection pool.
	return db, nil
}

func Get() zerolog.Logger {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

		logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
		if err != nil {
			logLevel = int(zerolog.InfoLevel) // default to INFO
		}

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		if os.Getenv("ENV") != "development" {
			fileLogger := &lumberjack.Logger{
				Filename:   "wikipedia-demo.log",
				MaxSize:    5, //
				MaxBackups: 10,
				MaxAge:     14,
				Compress:   true,
			}

			output = zerolog.MultiLevelWriter(os.Stderr, fileLogger)
		}

		var gitRevision string

		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, v := range buildInfo.Settings {
				if v.Key == "vcs.revision" {
					gitRevision = v.Value
					break
				}
			}
		}

		Log = zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Str("git_revision", gitRevision).
			Str("go_version", buildInfo.GoVersion).
			Logger()
	})

	return Log
}
