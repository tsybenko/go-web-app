package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"web-app/configs"
	"web-app/models/feedback"
)

const PORT = "3000"

type (
	View struct {
		Title     string
		Lang      string
		Content   string
		Feedbacks map[int]feedback.Feedback
	}
	Route struct {
		Name    string
		Handler func(w http.ResponseWriter, r *http.Request)
	}
)

var routes = make(map[string]Route)

func main() {
	godotenv.Load()
	router := mux.NewRouter()
	cnf := configs.DbConfig
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cnf["username"], cnf["password"], cnf["host"], cnf["port"], cnf["dbname"]))

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	routes["index"] = Route{
		Name: "index",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			data := View{
				Title: "Index page",
				Lang:  "uk",
			}

			tpl, err := template.ParseFiles("page/index.html")
			if err != nil {
				panic(err.Error())
			}
			tpl.Execute(w, data)
		},
	}

	routes["feedback"] = Route{
		Name: "Feedback",
		Handler: func(w http.ResponseWriter, r *http.Request) {

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
				http.Redirect(w, r, "/feedback", 301)
			} else {
				//feedbacks := make(map[int]Feedback)
				feedbacks := feedback.GetAll(db)
				data := View{
					Title:     "Feedback page",
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
		},
	}

	routes["feedbackDelete"] = Route{
		Name: "Feedback delete",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			params := mux.Vars(r)
			id, _ := strconv.Atoi(params["id"])
			feedback.Delete(db, id)

			if r.Method == "POST" {
				http.Redirect(w, r, "/feedback", 301)
			} else {
				res := make(map[string]bool)
				res["ok"] = true
				json.NewEncoder(w).Encode(res)
			}
		},
	}

	router.HandleFunc("/", routes["index"].Handler)
	router.HandleFunc("/feedback", routes["feedback"].Handler).Methods("GET", "POST")
	router.HandleFunc("/feedback/{id}", routes["feedbackDelete"].Handler).Methods("DELETE")
	router.HandleFunc("/feedback/delete/{id}", routes["feedbackDelete"].Handler).Methods("POST")

	fmt.Println(fmt.Sprintf("Server starting on :%s port", PORT))
	http.ListenAndServe(fmt.Sprintf(":%s", PORT), router)
}
