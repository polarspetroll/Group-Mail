package main

import (
	"html/template"
	"net/http"
	"os"
)

var tmps = template.Must(template.ParseGlob("templates/*.gohtml"))

type HTMLResponse struct {
	Username   string
	StatusText string
}

var (
	envusr    = os.Getenv("USERNAME")
	envpwd    = os.Getenv("PASSWORD")
	emailaddr = os.Getenv("SMTPUSR")
	smtppwd   = os.Getenv("SMTPPWD")
	smtphost  = os.Getenv("SMTPHOST")
	smtpport  = os.Getenv("SMTPPORT")
)

func initial() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/index", 302) })
	http.Handle("/statics/", http.StripPrefix("/statics/", http.FileServer(http.Dir("statics"))))
	http.HandleFunc("/index", Index)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/api/mail", MailAPI)
}

func main() {
	Encrypt(&envpwd)
	initial()
	go CookieInterval()
	http.ListenAndServe(":80", nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	c, err := r.Cookie("SID")
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	stat, username := CheckCookie(*c)
	if !stat {
		http.Redirect(w, r, "/login", 302)
		return
	}
	tmps.ExecuteTemplate(w, "index.gohtml", HTMLResponse{Username: username, StatusText: ""})
}
