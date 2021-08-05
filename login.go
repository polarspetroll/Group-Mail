package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type session struct {
	username string
	sid      string
	expires  time.Duration
}

var sessions []session
var wg sync.WaitGroup

func LoginAPI(w http.ResponseWriter, r *http.Request) {

}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmps.ExecuteTemplate(w, "login.gohtml", nil)
		return
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")
		if strings.ReplaceAll(username, " ", "") == "" || strings.ReplaceAll(password, " ", "") == "" {
			tmps.ExecuteTemplate(w, "login.gohtml", HTMLResponse{Username: "", StatusText: "Invalid Username or Password"})
			return
		}

		Encrypt(&password)
		if username != envusr || password != envpwd {
			tmps.ExecuteTemplate(w, "login.gohtml", HTMLResponse{Username: "", StatusText: "Incorrect Username or Password"})
			return
		} else if username == envusr && password == envpwd {
			cookie := GenerateSession(username)
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/index", 302)
			return
		}
	}

}

func CheckCookie(c http.Cookie) (bool, string) {
	value := fmt.Sprintf("%v", c.Value)
	for _, v := range sessions {
		if v.sid == value {
			return true, v.username
		}
	}
	return false, ""
}

func GenerateSession(username string) (cookie http.Cookie) {
	a := make([]byte, 10)
	rand.Read(a)
	value := fmt.Sprintf("%x", a)
	exp, _ := time.ParseDuration("5h")
	sessions = append(sessions, session{username: username, sid: value, expires: exp})
	cookie = http.Cookie{Name: "SID", Expires: time.Now().Add(exp), Value: value, HttpOnly: true}
	return cookie
}

func Encrypt(password *string) {
	hash := sha256.New()
	hash.Write([]byte(*password))
	*password = fmt.Sprintf("%x", hash.Sum(nil))
}

func CookieInterval() {
	arlen := len(sessions)
Loop:
	for {
		if len(sessions) != arlen {
			arlen = len(sessions)
			continue Loop
		}
		wg.Add(len(sessions))
		for i, v := range sessions {
			if len(sessions) != arlen {
				arlen = len(sessions)
				continue Loop
			}
			go func(v time.Duration, i int) {
				time.Sleep(v)
				RemoveSession(i)
				wg.Done()
			}(v.expires, i)
			if len(sessions) != arlen {
				arlen = len(sessions)
				continue Loop
			}
		}
		if len(sessions) != arlen {
			arlen = len(sessions)
			continue Loop
		}
		wg.Wait()

	}
}

func RemoveSession(index int) {
	sessions = append(sessions[:index], sessions[index+1:]...)
}
