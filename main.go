package main

import (
	"fmt"

	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	// "tawesoft.co.uk/go/dialog"
)

// var s:=sessions.NewCookieStore([]byte("francis"))
var tpl *template.Template
var Store = sessions.NewCookieStore([]byte("francis"))

func init() {
	tpl = template.Must(template.ParseGlob("static/*.html"))
}

type Page struct {
	Status  bool
	Header1 interface{}
	Valid   bool
}

var userDB = map[string]string{
	"email":    "fa_z_il_nbr",
	"password": "123456",
}
var P = Page{
	Status: false,
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,must-revalidate")

	ok := Middleware(w, r)

	if ok {

		http.Redirect(w, r, "/login-submit", http.StatusSeeOther)
		return
	}
	P.Valid = Middleware(w, r)
	filename := "login.html"
	err := tpl.ExecuteTemplate(w, filename, P)
	if err != nil {
		// fmt.Println("error while parsing file", err)
		return
	}

}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "there is an error parsing %v", err)
		return
	}
	emails := r.PostForm.Get("username")

	password := r.PostForm.Get("password")

	if userDB["email"] == emails && userDB["password"] == password && r.Method == "POST" {

		session, _ := Store.Get(r, "started")

		session.Values["id"] = emails
		P.Header1 = session.Values["id"]
		fmt.Println(P.Header1)
		session.Save(r, w)

		// fmt.Println(session)

		w.Header().Set("Cache-Control", "no-cache, must-revalidate")

		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else {
		// dialog.Alert("wrong passwod")
		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return

	}

}
func Logouthandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, must-revalidate")

	if P.Status == true {
		session, _ := Store.Get(r, "started")
		session.Options.MaxAge = -1
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		P.Status = false
	} else if P.Status == false {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func Middleware(w http.ResponseWriter, r *http.Request) bool {
	session, _ := Store.Get(r, "started")

	// fmt.Println(w)
	if session.Values["id"] == nil {
		return false
	}
	P.Header1 = session.Values["id"]
	return true

}

func index(w http.ResponseWriter, r *http.Request) {
	ok := Middleware(w, r)
	if ok {
		P.Status = true

	}
	filenamE := "index.html"
	err := tpl.ExecuteTemplate(w, filenamE, P)
	if err != nil {
		// fmt.Println("error while parsing file", err)
		return
	}

}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login-submit", loginHandler)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", Logouthandler)
	fmt.Println("server starts at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
