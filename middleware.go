package celeritas

import (
	"net/http"
	"strconv"

	"github.com/justinas/nosurf"
)

func (c *Celeritas) SessionLoad(next http.Handler) http.Handler {
	return c.Session.LoadAndSave(next)
}

func (c *Celeritas) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	secure, _ := strconv.ParseBool(c.config.cookie.secure)

	csrfHandler.ExemptGlob("/api/*")

	if !c.Debug {
		csrfHandler.SetBaseCookie(http.Cookie{
			HttpOnly: true,
			Path:     "/",
			Secure:   secure,
			SameSite: http.SameSiteStrictMode,
			Domain:   c.config.cookie.domain,
		})
	} else {
		csrfHandler.SetBaseCookie(http.Cookie{
			HttpOnly: false,
			Path:     "/",
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
	}

	return csrfHandler
}
