package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	//"fmt"
)

type Page struct {
	Title string
	Body  []byte
}

type User struct {
	Name string
	Password string
	Post []string
}

type LogResult struct {
	Name string
	Password string
	Result string
	Message string
}

//users []User*
var users map[string]*User

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

/*func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}*/

func loginHandler(w http.ResponseWriter, r *http.Request, title string) {
	//http.Redirect(w, r, "/logresult", http.StatusFound)
	var u User
	renderTemplate(w, "login", u)
	u.Name = r.FormValue("name")
	u.Password = r.FormValue("password")
}
func loginresultHandler(w http.ResponseWriter, r *http.Request, title string) {
	logres := LogResult { Name:r.FormValue("name"), Password:r.FormValue("password")}
	if(r.FormValue("choose") == "Log in"){
		pw, ok := users[r.FormValue("name")]
		if(pw == r.FormValue("password")){
			logres.Result = "successfully"
		}else{
			logres.Result = "failed"
			logres.Message = "Wrong user."
			if ok {
				logres.Message = "Wrong password."
			}
		}
		renderTemplate(w, "loginresult", logres)

	}else{
		tmp := User{Name: r.FormValue("name"), Password:r.FormValue("password")}
		users[r.FormValue("name")] = tmp
		renderTemplate(w, "signup", logres)
		//http.Redirect(w, r, "/signup/", http.StatusFound)
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request, title string) {
	new_user := User{ Name:r.FormValue("name"), Password:r.FormValue("password")} 
	renderTemplate(w, "signup", new_user)
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html", "login.html", "loginresult.html", "signup.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view|login|loginresult|signup)/")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[1])
	}
}

func main() {
	users = make(map[string]string)
	//http.HandleFunc("/view/", makeHandler(viewHandler))
	//http.HandleFunc("/edit/", makeHandler(editHandler))
	//http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/login/", makeHandler(loginHandler))
	http.HandleFunc("/loginresult/", makeHandler(loginresultHandler))
	http.HandleFunc("/signup/", makeHandler(signupHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
