package session

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
)

// Session represents the session configuration
type Session struct {
	CookieLifetime string
	CookiePersist  string
	CookieName     string
	CookieDomain   string
	SessionType    string
	CookieSecure   string
}

func (c *Session) InitSession() *scs.SessionManager {
	var persist, secure bool

	minutes, err := strconv.Atoi(c.CookieLifetime)
	if err != nil {
		minutes = 60
	}

	persist = strings.ToLower(c.CookiePersist) == "true"
	secure = strings.ToLower(c.CookieSecure) == "true"

	session := scs.New()
	session.Lifetime = time.Duration(minutes) * time.Minute
	session.Cookie.Persist = persist
	session.Cookie.Name = c.CookieName
	session.Cookie.Secure = secure
	session.Cookie.Domain = c.CookieDomain
	session.Cookie.SameSite = http.SameSiteLaxMode

	switch strings.ToLower(c.SessionType) {
	case "redis":
		// session.Redis = c.redisConfig()
	default:
		// cookie
	}

	return session
}
