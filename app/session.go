package main

import (
	"encoding/json"
	"net/http"
	"time"

	sessions "github.com/goincremental/negroni-sessions"
)

const (
	currentUserKey  = "oauth2_current_user"
	sessionDuration = time.Hour
)

// User represent a chat user
type User struct {
	UID       string    `json:"uid"`
	Name      string    `json:"name"`
	Email     string    `json:"user"`
	AvatarURL string    `json:"avatar_url"`
	Expired   time.Time `json:"expired"`
}

// Valid check if the session is valid
func (u *User) Valid() bool {
	return u.Expired.Sub(time.Now()) > 0
}

// Refresh expired time
func (u *User) Refresh() {
	u.Expired = time.Now().Add(sessionDuration)
}

// GetCurrentUser get current user info from session
func GetCurrentUser(r *http.Request) *User {
	s := sessions.GetSession(r)
	if s.Get(currentUserKey) == nil {
		return nil
	}

	data := s.Get(currentUserKey).([]byte)
	var u User
	json.Unmarshal(data, &u)
	return &u
}

// SetCurrentUser set current user info to session
func SetCurrentUser(r *http.Request, u *User) {
	if u != nil {
		u.Refresh()
	}

	s := sessions.GetSession(r)
	val, _ := json.Marshal(u)
	s.Set(currentUserKey, val)
}
