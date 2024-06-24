package main

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"time"
)

const (
	port = ":7777"
	user = "adil000"
	pass = "12345000"
)

var templates = template.Must(template.ParseFiles("login.html", "home.html"))
var sessionStore = make(map[string]time.Time)

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == user && password == pass {
			sessionID := generateSessionID()
			expire := time.Now().Add(30 * time.Minute)
			sessionStore[sessionID] = expire

			cookie := http.Cookie{
				Name:    "session_id",
				Value:   sessionID,
				Expires: expire,
				Path:    "/",
			}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}

		templates.ExecuteTemplate(w, "login.html", map[string]string{
			"Error": "INCORRECT USERNAME OR PASSWORD",
		})
		return
	}

	templates.ExecuteTemplate(w, "login.html", nil)

}

func home(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	w.Header().Set("cache-control", "no-store")
	err = templates.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		http.Error(w, "error executing templete", http.StatusInternalServerError)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil && cookie != nil {
		delete(sessionStore, cookie.Value)
		cookie := http.Cookie{
			Name:    "session_id",
			Value:   "",
			Expires: time.Now().Add(-1 * time.Hour),
			Path:    "/",
		}
		http.SetCookie(w, &cookie)
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", login)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/home", home)
	mux.HandleFunc("/logout", logout)

	log.Println("Server started at :7777")
	log.Fatal(http.ListenAndServe(port, mux))

}
