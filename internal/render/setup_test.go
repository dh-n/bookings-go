package render

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dh-n/bookings/internal/config"
)

var (
	session *scs.SessionManager
	testApp config.AppConfig
)

func TestMain(m *testing.M) {
	testApp.InProduction = false

	// setting the session defaults for scs
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session
	app = &testApp

	os.Exit(m.Run())
}

type myWriter struct{}

func (mw *myWriter) Header() http.Header {
	h := make(http.Header)
	return h
}

func (mw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}

func (mw *myWriter) WriteHeader(statusCode int) {}
