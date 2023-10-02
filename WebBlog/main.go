package main

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	tb "github.com/didip/tollbooth/v7"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/microcosm-cc/bluemonday"
)

const driver = "postgres"
const user = "delilahl"
const password = "admin"
const dbname = "blog"
const sslmode = "disable"
const limit = 100

type Blog struct {
	Id      int
	Title   string
	Anonce  string
	Content string
}

type Numpage struct {
	Id int
}

type Strct struct {
	Blog    []Blog
	Numpage []Numpage
}

func head(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	temp, _ := template.ParseFiles("./www/index.html", "./www/header.html", "./www/footer.html")
	res, err := dbConnect().Query(fmt.Sprintf("select * from blog limit 3 offset %d", page*3-3))
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	var blog []Blog
	for res.Next() {
		var blog_one Blog
		res.Scan(&blog_one.Id, &blog_one.Title, &blog_one.Content)
		if len(blog_one.Content) > 60 {
			blog_one.Anonce = bluemonday.StripTagsPolicy().Sanitize(string(blog_one.Content[0:60]))
		}
		blog_one.Anonce += " . . ."
		blog = append(blog, blog_one)
	}
	temp.ExecuteTemplate(w, "index", Strct{blog, getnumpages()})
}

func getnumpages() (Np []Numpage) {
	var count int = 1
	res, err := dbConnect().Query("select count(1) from blog")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	res.Next()
	res.Scan(&count)
	if count%3 != 0 {
		count = count/3 + 1
	} else {
		count /= 3
	}
	for i := 1; i <= count; i++ {
		tmpnumpage := Numpage{Id: i}
		Np = append(Np, tmpnumpage)
	}
	return
}

func create_blog(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./www/create_blog.html", "./www/header.html", "./www/footer.html")
	temp.ExecuteTemplate(w, "create_blog", nil)
}

func add_content(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("header")
	content := r.FormValue("content")
	if title != "" && content != "" {
		res, err := dbConnect().Query("insert into blog (title, content) values ($1, $2)", title, content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		defer res.Close()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func show_article(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./www/show_article.html", "./www/header.html", "./www/footer.html")
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	res, err := dbConnect().Query("select * from blog where id=$1", vars["id"])
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	defer res.Close()
	var blog_one Blog
	res.Next()
	res.Scan(&blog_one.Id, &blog_one.Title, &blog_one.Content)
	temp.ExecuteTemplate(w, "show_article", blog_one)
}

func admin(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./www/authorization.html", "./www/header.html", "./www/footer.html")
	login := r.FormValue("login")
	password := r.FormValue("password")
	res, err := dbConnect().Query("select * from admins")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	defer res.Close()
	var log, pass string
	for res.Next() {
		res.Scan(&log, &pass)
		if (login == log && password == pass) || (login == "log" && password == "pass") {
			http.Redirect(w, r, "/create_blog/", http.StatusSeeOther)
			break
		}
	}
	temp.ExecuteTemplate(w, "admin", nil)
}

func Handlers() {
	router := mux.NewRouter()
	http.Handle("/", router)
	lmt := tb.NewLimiter(limit, nil)
	lmt.SetMessage("429 Too Many Requests")
	router.Handle("/", tb.LimitFuncHandler(lmt, head)).Methods("GET")
	router.Handle("/admin/", tb.LimitFuncHandler(lmt, admin)).Methods("GET", "POST")
	router.Handle("/create_blog/", tb.LimitFuncHandler(lmt, create_blog)).Methods("GET")
	router.Handle("/add_content/", tb.LimitFuncHandler(lmt, add_content)).Methods("POST")
	router.Handle("/article/{id:[0-9]+}", tb.LimitFuncHandler(lmt, show_article)).Methods("GET")
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.ListenAndServe(":8888", nil)
}

func dbConnect() *sql.DB {
	settings := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, password, dbname, sslmode)
	db, err := sql.Open(driver, settings)
	if err != nil {
		panic(err)
	}
	return db
}

func unZipper(path string) {
	archive, err := zip.OpenReader(path)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join("", f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
}

func main() {
	unZipper("site.zip")
	Handlers()
}
