package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
	"strings"
	"web-app/models/feedback"
)

const PORT = "3000"

type View struct {
	Title     string
	Lang      string
	Content   string
	Feedbacks map[int]feedback.Feedback
}

var router = make(map[string]func(w http.ResponseWriter, r *http.Request))

func main() {

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/gobase")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	router["index"] = func(w http.ResponseWriter, r *http.Request) {
		data := View{
			Title: "Index page",
			Lang:  "uk",
		}

		tpl, err := template.ParseFiles("page/index.html")
		if err != nil {
			panic(err.Error())
		}
		tpl.Execute(w, data)
	}

	router["contact"] = func(w http.ResponseWriter, r *http.Request) {

		// Catch POST request
		if strings.ToLower("post") == strings.ToLower(r.Method) {
			err := r.ParseForm()
			if err != nil {
				panic(err.Error())
			}
			form := &r.PostForm
			ok := feedback.Create(db, &feedback.Feedback{
				Email: form.Get("email"),
				Text:  form.Get("text"),
			})
			if ok {
				fmt.Println("Record created!")
			} else {
				fmt.Println("Record not created :(")
			}
			http.Redirect(w, r, "/contact", 301)
		} else {
			//feedbacks := make(map[int]Feedback)
			feedbacks := feedback.GetAll(db)
			data := View{
				Title:     "Contact page",
				Lang:      "uk",
				Feedbacks: feedbacks,
			}
			tpl, err := template.ParseFiles("page/contact.html")
			if err != nil {
				panic(err.Error())
			}
			err = tpl.Execute(w, data)
			if err != nil {
				panic(err.Error())
			}
		}
	}

	http.HandleFunc("/", router["index"])
	http.HandleFunc("/contact", router["contact"])

	fmt.Println(fmt.Sprintf("Server starting on :%s port", PORT))
	http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil)
}
