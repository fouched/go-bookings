package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/fouched/go-bookings/internal/config"
	"github.com/fouched/go-bookings/internal/models"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

// TestMain gets called before any test are run, runs the test just before the application closes
func TestMain(m *testing.M) {
	// define complex types that will
	// be stored in the session
	gob.Register(models.Reservation{})

	// change to true for production
	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}
