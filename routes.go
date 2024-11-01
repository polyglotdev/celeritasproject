package celeritas

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (c *Celeritas) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	if c.Debug {
		mux.Use(middleware.Logger)
	}
	mux.Use(middleware.Recoverer)

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		n, err := fmt.Fprint(w, "Welcome to Celeritas")
		if err != nil {
			c.ErrorLog.Println(err)
			c.InfoLog.Printf("Number of bytes written: %d", n)
		}
	})

	return mux
}
