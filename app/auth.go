package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/urfave/negroni"

	"github.com/stretchr/objx"

	sessions "github.com/goincremental/negroni-sessions"
	"github.com/julienschmidt/httprouter"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
)

const (
	nextPageKey     = "next_page"
	authSecurityKey = "auth_security_key"
)

func init() {
	gomniauth.SetSecurityKey(authSecurityKey)
	gomniauth.WithProviders(
		google.New("955753969182-09lfebpoa1atq2c85o1n4uln7cio47eu.apps.googleusercontent.com",
			"pldc7t_0bLr8t3JQMZupVcp1",
			"http://127.0.0.1:3000/auth/callback/google"))
}

func loginRequired(ignore ...string) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		for _, s := range ignore {
			if strings.HasPrefix(r.URL.Path, s) {
				next(w, r)
				return
			}
		}

		u := GetCurrentUser(r)

		if u != nil && u.Valid() {
			SetCurrentUser(r, u)
			next(w, r)
			return
		}

		SetCurrentUser(r, nil)
		fmt.Println(r.URL.RequestURI())
		sessions.GetSession(r).Set(nextPageKey, r.URL.RequestURI())
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func loginHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	action := ps.ByName("action")
	provider := ps.ByName("provider")
	s := sessions.GetSession(req)

	switch action {
	case "login":
		p, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln(err)
		}
		loginURL, err := p.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln(err)
		}
		http.Redirect(w, req, loginURL, http.StatusFound)
	case "callback":
		p, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln(err)
		}
		creds, err := p.CompleteAuth(objx.MustFromURLQuery(req.URL.RawQuery))
		if err != nil {
			log.Fatalln(err)
		}

		user, err := p.GetUser(creds)
		if err != nil {
			log.Fatalln(err)
		}

		u := &User{
			UID:       user.Data().Get("id").MustStr(),
			Name:      user.Name(),
			Email:     user.Email(),
			AvatarURL: user.AvatarURL(),
		}

		SetCurrentUser(req, u)
		temp := s.Get(nextPageKey)
		http.Redirect(w, req, temp.(string), http.StatusFound)
	default:
		http.Error(w, "Auth action '"+action+"'is not supported", http.StatusNotFound)
	}
}
