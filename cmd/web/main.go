package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/fouched/go-bookings/internal/config"
	"github.com/fouched/go-bookings/internal/driver"
	"github.com/fouched/go-bookings/internal/handlers"
	"github.com/fouched/go-bookings/internal/helpers"
	"github.com/fouched/go-bookings/internal/models"
	"github.com/fouched/go-bookings/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":9080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	// close database connection after application has stopped
	defer db.SQL.Close()
	// close mail channel
	defer close(app.MailChan)
	// start listening for mail
	listenForMail()
	fmt.Println("Started mail listener")

	fmt.Println(fmt.Sprintf("Starting application on %s", portNumber))
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// define complex types that will
	// be stored in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// read flags
	inProduction := flag.Bool("production", false, "Application is in production")
	useCache := flag.Bool("cache", false, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "bookings", "Database name")
	dbUser := flag.String("dbuser", "fouche", "Database user")
	dbPass := flag.String("dbpass", "javac", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database SSL settings (disable, prefer, require)")

	flag.Parse()
	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}


	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change to true for production
	app.InProduction = *inProduction
	app.UseCache = *useCache

	// set up loggers
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		*dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)

	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("Connected to database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}
	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
