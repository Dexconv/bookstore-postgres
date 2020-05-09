package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"github.com/Dexconv/postgres/model"
)

type errNcode struct{
	err error
	code int
}


type books []model.Book

var tpl *template.Template

func init(){
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func Bookindex (w http.ResponseWriter , r *http.Request){
	if r.Method != "GET"{
		errHandler(w, errNcode{errors.New("wrong method used"), 405})
	}
	rows ,err := model.AllBooks()
	errHandler(w, errNcode{err, 500})

	defer rows.Close()
	var bks books
	for rows.Next(){
		var bk model.Book
		err := rows.Scan(&bk.Isbn, &bk.Title, &bk.Author, &bk.Price)
		errHandler(w, errNcode{err, 500})

		bks = append(bks, bk)
	}

	tpl.ExecuteTemplate(w,"books.gohtml", bks)
}
func Onebook(w http.ResponseWriter , r *http.Request){
	if r.Method != "GET"{
		errHandler(w, errNcode{errors.New("wrong method used"), 405})
	}
	isbn := r.FormValue("isbn")
	fmt.Println(isbn, "was requested")
	row := model.OneBooks(isbn)
	var bk model.Book
	err := row.Scan(&bk.Isbn, &bk.Title, &bk.Author, &bk.Price)
	if err == sql.ErrNoRows{
		http.NotFound(w,r)
		return
	}else if err != nil {
		errHandler(w, errNcode{err, 500})
	}
	tpl.ExecuteTemplate(w,"book.gohtml", bk)
}
func BookCreateForm(w http.ResponseWriter , r *http.Request){
	tpl.ExecuteTemplate(w, "create.gohtml", nil )
}
func BookCreateProcess(w http.ResponseWriter , r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, http.StatusText(405), 405)
		return
	}
	bk := model.Book{}
	bk.Isbn = r.FormValue("Isbn")
	bk.Title = r.FormValue("Title")
	bk.Author = r.FormValue("Author")
	p := r.FormValue("Price")

	if bk.Isbn == "" || bk.Isbn == "" || bk.Author == "" || p == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	f64 ,err := strconv.ParseFloat(p, 32)
	if err != nil{
		http.Error(w, http.StatusText(406)+"Please hit back and enter a number for the price", http.StatusNotAcceptable)
		return
	}
	bk.Price = float32(f64)

	err = model.InsertBook(bk)

	if err != nil{
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}
	tpl.ExecuteTemplate(w, "created.gohtml", bk)
}
func BookUpdateForm(w http.ResponseWriter , r *http.Request){
	if r.Method != http.MethodGet{
		http.Error(w,http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	isbn := r.FormValue("isbn")
	if isbn == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	row := model.OneBooks(isbn)

	var bk model.Book

	err := row.Scan(&bk.Isbn,&bk.Title,&bk.Author,&bk.Price)

	if err == sql.ErrNoRows{
		http.NotFound(w,r)
		return
	}else if err != nil{
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	tpl.ExecuteTemplate(w, "update.gohtml", bk)
}
func BookUpdateProcess(w http.ResponseWriter , r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	bk := model.Book{}

	bk.Isbn = r.FormValue("Isbn")
	bk.Title = r.FormValue("Title")
	bk.Author = r.FormValue("Author")
	p := r.FormValue("Price")

	if bk.Isbn == "" || bk.Isbn == "" || bk.Author == "" || p == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	f64 , err := strconv.ParseFloat(p, 32)
	if err != nil{
		http.Error(w, http.StatusText(406)+"Please hit back and enter a number for the price", http.StatusNotAcceptable)
		return
	}
	f32 := float32(f64)
	bk.Price = f32

	err = model.UpdateBook(bk)
	if err != nil{
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteTemplate(w, "book.gohtml", bk)
}
func BookDeleteProcess( w http.ResponseWriter , r *http.Request){
	if r.Method != http.MethodGet{
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	isbn := r.FormValue("isbn")

	if isbn == ""{
		http.Error(w, http.StatusText(400)+isbn, http.StatusBadRequest)
		return
	}
	err := model.DeleteBook(isbn)
	if err != nil{
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	http.Redirect(w,r,"/books", http.StatusSeeOther)
}
func errHandler( w http.ResponseWriter , r errNcode){
	if r.err != nil{
		http.Error(w, http.StatusText(r.code), r.code)
		log.Fatalln(r.err)
		return
	}
}
