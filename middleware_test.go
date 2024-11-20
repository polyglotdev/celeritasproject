package celeritas

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
)

func TestCeleritas_SessionLoad(t *testing.T) {
	var cel *Celeritas // declare at test function scope

	tests := []struct {
		name        string
		nextHandler http.Handler
		expectedLog string
		setupFunc   func() *Celeritas
	}{
		{
			name: "session middleware test",
			nextHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Use cel instead of c
				cel.Session.Put(r.Context(), "test", "test value")
				_, err := w.Write([]byte("ok"))
				if err != nil {
					t.Error(err)
				}
			}),
			expectedLog: "SessionLoad called\n",
			setupFunc: func() *Celeritas {
				session := scs.New()
				session.Lifetime = 24 * time.Hour
				session.Cookie.Persist = true
				session.Cookie.SameSite = http.SameSiteLaxMode
				session.Cookie.Secure = false

				c := &Celeritas{
					Session: session,
				}
				var logBuffer bytes.Buffer
				c.InfoLog = log.New(&logBuffer, "", 0)
				return c
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			cel = tt.setupFunc()

			// Create test recorder and request
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)

			// Get the handler from SessionLoad
			handler := cel.SessionLoad(tt.nextHandler)

			// Call the handler
			handler.ServeHTTP(w, r)

			// Verify response
			if w.Code != http.StatusOK {
				ts.Errorf("SessionLoad() status = %v, want %v", w.Code, http.StatusOK)
			}

			// Verify session cookie was set
			if len(w.Result().Cookies()) == 0 {
				ts.Error("SessionLoad() no session cookie set")
			}

			// Verify log message
			logBuffer := cel.InfoLog.Writer().(*bytes.Buffer)
			logOutput := logBuffer.String()
			if logOutput != tt.expectedLog {
				ts.Errorf("SessionLoad() log = %v, want %v", logOutput, tt.expectedLog)
			}
		})
	}
}
