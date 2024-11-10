package session

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
)

// Package session provides session management functionality using the SCS library.

// Session represents the configuration settings for session management.
// It holds various cookie-related settings as string values that will be
// parsed during session initialization.
type Session struct {
	// CookieLifetime sets how long the session should last in minutes.
	// If empty or invalid, defaults to 60 minutes.
	CookieLifetime string

	// CookiePersist determines if the cookie should persist after browser close.
	// Valid values are "true" or "false" (case insensitive).
	CookiePersist string

	// CookieName sets the name of the session cookie.
	// If empty, defaults to "session" in the SCS library.
	CookieName string

	// CookieDomain sets the domain for the session cookie.
	// If empty, defaults to the domain name that the cookie was issued from.
	CookieDomain string

	// SessionType determines the storage backend for sessions.
	// Valid values are "redis", "mysql", "mariadb", "postgres", "postgresql".
	// Any other value defaults to cookie-based sessions.
	SessionType string

	// CookieSecure determines if the cookie should only be transmitted over HTTPS.
	// Valid values are "true" or "false" (case insensitive).
	CookieSecure string
	DBPool       *sql.DB
}

// InitSession initializes and returns a new session manager using the
// configuration settings in the Session struct. It handles parsing of
// string configuration values and sets up the session manager with
// appropriate timeout, cookie settings, and storage backend.
//
// The returned *scs.SessionManager is configured and ready to use for
// session management in a web application.
func (c *Session) InitSession() *scs.SessionManager {
	var persist, secure bool

	// Parse session lifetime from string to minutes
	minutes, err := strconv.Atoi(c.CookieLifetime)
	if err != nil {
		minutes = 60 // Default to 60 minutes if invalid or empty
	}

	// Parse persistence setting
	if strings.ToLower(c.CookiePersist) == "true" {
		persist = true
	}

	// Parse secure setting
	if strings.ToLower(c.CookieSecure) == "true" {
		secure = true
	}

	// Initialize and configure session manager
	session := scs.New()
	session.Lifetime = time.Duration(minutes) * time.Minute
	session.Cookie.Persist = persist
	session.Cookie.Name = c.CookieName
	session.Cookie.Secure = secure
	session.Cookie.Domain = c.CookieDomain
	session.Cookie.SameSite = http.SameSiteLaxMode

	// Configure session store based on SessionType
	switch strings.ToLower(c.SessionType) {
	case "redis":

	case "mysql", "mariadb":
		session.Store = mysqlstore.New(c.DBPool)

	case "postgres", "postgresql":
		session.Store = postgresstore.New(c.DBPool)

	default:
		// Default to cookie-based session store
	}

	return session
}
