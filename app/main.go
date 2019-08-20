package main

import (
	"net/http"

	sessions "github.com/goincremental/negroni-sessions"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
)

var renderer *render.Render

const (
	sessionKey    = "simple_chat_session"
	sessionSecret = "simple_chat_session_secret"
)

func init() {
	renderer = render.New(render.Options{
		IsDevelopment: true,
	})
}

func main() {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		renderer.HTML(w, http.StatusOK, "index", map[string]string{"title": "Simple Chat"})
	})

	router.GET("/login", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		renderer.HTML(w, http.StatusOK, "login", nil)
	})

	router.GET("/logout", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		sessions.GetSession(req).Delete(currentUserKey)
		http.Redirect(w, req, "/login", http.StatusFound)
	})

	router.GET("/auth/:action/:provider", loginHandler)

	n := negroni.Classic()
	// store := cookiestore.New([]byte(sessionSecret))
	// n.Use(sessions.Sessions(sessionKey, store))
	n.UseHandler(router)

	n.Run(":3000")
}
