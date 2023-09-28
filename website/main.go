package main

import (
	"fmt"
	"net/http"
	"html/template"
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

type Article struct{
	Id uint16
	Title, Anons, FullText string
}
var posts = []Article{}
var showpost = Article{}

func index(w http.ResponseWriter, r *http.Request){
	t, err :=template.ParseFiles("templates/index.html", "templates/header.html","templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			panic(err)
		} 

		defer db.Close()
		// ввыод статей

		res, err :=db.Query("select * from articles")
	if err != nil {
        panic(err)
    }

	posts=[]Article{}

	for res.Next(){
		var post Article
		err =res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post )
		
	}
		

	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request){
	t, err :=template.ParseFiles("templates/create.html", "templates/header.html","templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request){
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")
	

    if title != "" && anons!="" && full_text !=""{
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			panic(err)
		} 
		defer db.Close()
		result, err := db.Exec(fmt.Sprintf("INSERT INTO `articles`(`title`, `anons`, `full_text`) VALUES('%s','%s','%s')", title, anons, full_text))
		if err != nil{
			panic(err)
		}
		fmt.Println(result.LastInsertId()) 
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else{
		
		fmt.Fprintf(w, "Заполнити все поля")
	}
}

func show_post(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	//w.WriteHeader(http.StatusOK) 
	// fmt.Fprintf(w, "ID: %v\n", vars["id"])
	t, err :=template.ParseFiles("templates/show.html", "templates/header.html","templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			panic(err)
		} 

		defer db.Close()
// ввывод данных
		res, err :=db.Query(fmt.Sprintf("SELECT * FROM `articles` WHERE `id` ='%s'", vars["id"]))
	
	if err != nil {
        panic(err)
    }

	showpost = Article{}

	for res.Next(){
		var post Article
		err =res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		showpost= post
		
	}

	t.ExecuteTemplate(w, "show", showpost)
}

func handaleFunc(){

	rtr := mux.NewRouter()

	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")

	http.Handle("/", rtr)
	http.ListenAndServe(":8080", nil)
}

func main() {
	handaleFunc()
}
