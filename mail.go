package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"strings"
)

type Email struct {
	From    string `json: "from"`
	To      string `json: "to"`
	Body    string `json: "body"`
	Subject string `json: "subject"`
}

func MailAPI(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SID")
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	stat, _ := CheckCookie(*cookie)
	if !stat {
		http.Redirect(w, r, "/login", 302)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", 405)
	}
	/**********************************************************/
	w.Header().Set("Content-Type", "application/json")

	var out Email

	body, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &out)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(`{"ok": false, "status": "invalid request"}`))
		return
	}
	if err = SendMail(out); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`{"ok": false, "status": "%v"}`, err)))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(`{"ok": true, "status": "done"}`))

}

func SendMail(mail Email) (err error) {
	clients := ParseClients(mail.To)
	auth := smtp.PlainAuth("", emailaddr, smtppwd, smtphost)
	text := []byte(fmt.Sprintf("Subject: %v\r\n\r\n%v", mail.Subject, mail.Body))
	err = smtp.SendMail(smtphost+":"+smtpport, auth, emailaddr, clients, text)
	return err
}

func ParseClients(clients string) (out []string) {
	out = strings.Split(clients, "\n")
	return out
}
