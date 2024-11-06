package session

import (
	"net/http"
	"testing"
	"time"
)

func TestSession_InitSession(t *testing.T) {
	tests := []struct {
		name            string
		session         Session
		wantLifetime    time.Duration
		wantPersist     bool
		wantSecure      bool
		wantName        string
		wantDomain      string
		wantSessionType string
	}{
		{
			name: "default values",
			session: Session{
				CookieLifetime: "",
				CookiePersist:  "",
				CookieName:     "session",
				CookieDomain:   "localhost",
				SessionType:    "",
				CookieSecure:   "",
			},
			wantLifetime:    60 * time.Minute,
			wantPersist:     false,
			wantSecure:      false,
			wantName:        "session",
			wantDomain:      "localhost",
			wantSessionType: "cookie",
		},
		{
			name: "custom values",
			session: Session{
				CookieLifetime: "30",
				CookiePersist:  "true",
				CookieName:     "mysession",
				CookieDomain:   "example.com",
				SessionType:    "redis",
				CookieSecure:   "true",
			},
			wantLifetime:    30 * time.Minute,
			wantPersist:     true,
			wantSecure:      true,
			wantName:        "mysession",
			wantDomain:      "example.com",
			wantSessionType: "redis",
		},
		{
			name: "invalid lifetime",
			session: Session{
				CookieLifetime: "invalid",
				CookiePersist:  "false",
				CookieName:     "session",
				CookieDomain:   "localhost",
				SessionType:    "",
				CookieSecure:   "false",
			},
			wantLifetime:    60 * time.Minute, // should default to 60
			wantPersist:     false,
			wantSecure:      false,
			wantName:        "session",
			wantDomain:      "localhost",
			wantSessionType: "cookie",
		},
		{
			name: "case insensitive booleans",
			session: Session{
				CookieLifetime: "60",
				CookiePersist:  "TRUE",
				CookieName:     "session",
				CookieDomain:   "localhost",
				SessionType:    "REDIS",
				CookieSecure:   "True",
			},
			wantLifetime:    60 * time.Minute,
			wantPersist:     true,
			wantSecure:      true,
			wantName:        "session",
			wantDomain:      "localhost",
			wantSessionType: "redis",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			session := tt.session.InitSession()

			if session.Lifetime != tt.wantLifetime {
				t.Errorf("InitSession() Lifetime = %v, want %v", session.Lifetime, tt.wantLifetime)
			}

			if session.Cookie.Persist != tt.wantPersist {
				t.Errorf("InitSession() Cookie.Persist = %v, want %v", session.Cookie.Persist, tt.wantPersist)
			}

			if session.Cookie.Secure != tt.wantSecure {
				t.Errorf("InitSession() Cookie.Secure = %v, want %v", session.Cookie.Secure, tt.wantSecure)
			}

			if session.Cookie.Name != tt.wantName {
				t.Errorf("InitSession() Cookie.Name = %v, want %v", session.Cookie.Name, tt.wantName)
			}

			if session.Cookie.Domain != tt.wantDomain {
				t.Errorf("InitSession() Cookie.Domain = %v, want %v", session.Cookie.Domain, tt.wantDomain)
			}

			if session.Cookie.SameSite != http.SameSiteLaxMode {
				t.Errorf("InitSession() Cookie.SameSite = %v, want %v (http.SameSiteLaxMode)",
					session.Cookie.SameSite, http.SameSiteLaxMode)
			}
		})
	}
}
